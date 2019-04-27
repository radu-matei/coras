package coras

import (
	"fmt"
	"strings"

	"github.com/deislabs/cnab-go/bundle"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry"
)

// RelocateBundleImages pushes all referenced images to the a new repository.
// In the new repository, images are uniquely identified by the digest.
// Currently, each image also has a unique tag, but the human readable tag should
// never be used when referencing the image.
//
// The bundle is mutated in place, and contains the new image location (and digest, if it wasn't previously present)
func RelocateBundleImages(b *bundle.Bundle, targetRef string) error {
	rc := registry.NewRegistryClient()
	for i := range b.InvocationImages {
		ii := b.InvocationImages[i]
		_, err := relocateImage(&ii.BaseImage, targetRef, rc)
		if err != nil {
			return err
		}
		b.InvocationImages[i] = ii
	}

	for i := range b.Images {
		im := b.Images[i]
		_, err := relocateImage(&im.BaseImage, targetRef, rc)
		if err != nil {
			return err
		}
		b.Images[i] = im
	}

	return nil
}

func relocateImage(i *bundle.BaseImage, targetRef string, client registry.Client) (bool, error) {
	if !isAcceptedImageType(i.ImageType) {
		return false, fmt.Errorf("cannot relocate image of type %v", i.ImageType)
	}

	// original image FQDN
	originalImage, err := image.NewName(i.Image)
	if err != nil {
		return false, fmt.Errorf("cannot get fully qualified image name for %v: %v", i.Image, err)
	}

	// relocated image FQDN
	// the location of the new image is the same target repository as the bundle itself,
	// but a different, unique tag.
	//
	// Note that the tag is just a familiar name, and it is not used
	// All references to this image are through the SHA digest
	//
	// TODO - @radu-matei
	// make sure naming strategy is consistent
	newImage, err := image.NewName(fmt.Sprintf("%s:%s", strings.Split(targetRef, ":")[0], strings.Replace(i.Image, ":", "-", -1)))
	if err != nil {
		return false, fmt.Errorf("cannot get fully qualified image name for the new image in %v: %v", targetRef, err)
	}

	dig, err := client.Copy(originalImage, newImage)
	if err != nil {
		return false, fmt.Errorf("cannot copy original image %v into new image %v: %v", originalImage.Name(), newImage.Name(), err)
	}

	// TODO - @radu-matei
	// make sure the digest is not modified, and return true if it is
	i.Digest = dig.String()
	i.OriginalImage = i.Image
	i.Image = newImage.String()

	return false, nil
}

func isAcceptedImageType(imageType string) bool {
	return imageType == "" || imageType == "oci" || imageType == "docker"
}
