package main

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/deislabs/cnab-go/bundle"
	"github.com/pivotal/image-relocation/pkg/image"
	"github.com/pivotal/image-relocation/pkg/registry"
	"github.com/spf13/cobra"
)

type relocateCmd struct {
	inputBundle string
	targetRef   string

	client registry.Client
}

func newRelocateCmd() *cobra.Command {
	const usage = `relocates the images referenced in a bundle to a new repository`

	var r relocateCmd
	cmd := &cobra.Command{
		Use:   "relocate",
		Short: usage,
		Long:  usage,
		RunE: func(cmd *cobra.Command, args []string) error {
			r.inputBundle = args[0]
			r.targetRef = args[1]
			r.client = registry.NewRegistryClient()
			return r.run()
		},
	}
	return cmd
}

func (r *relocateCmd) run() error {
	data, err := ioutil.ReadFile(r.inputBundle)
	if err != nil {
		return fmt.Errorf("cannot read bundle file: %v", err)
	}

	b, err := bundle.Unmarshal(data)
	if err != nil {
		return fmt.Errorf("cannot unmarshal bundle file: %v", err)
	}

	for i := range b.InvocationImages {
		ii := b.InvocationImages[i]
		_, err := r.relocateImage(&ii.BaseImage)
		if err != nil {
			return fmt.Errorf("cannot relocate invocation image: %v", err)
		}
		//if modified {
		b.InvocationImages[i] = ii
		//}
	}

	for k := range b.Images {
		im := b.Images[k]
		_, err := r.relocateImage(&im.BaseImage)
		if err != nil {
			return err
		}
		//if modified {
		b.Images[k] = im
		//}
	}

	// TODO - @radu-matei
	// make sure digest is persisted
	// get image size from the registry client

	err = b.WriteFile("output.json", 0644)
	if err != nil {
		return fmt.Errorf("cannot write output file")
	}

	p := pushCmd{
		inputBundle: "output.json",
		targetRef:   r.targetRef,
	}

	return p.run()
}

func (r *relocateCmd) relocateImage(i *bundle.BaseImage) (bool, error) {
	if !isOCI(i.ImageType) && !isDocker(i.ImageType) {
		return false, nil
	}
	// map the image name
	n, err := image.NewName(i.Image)
	if err != nil {
		return false, err
	}

	nn, err := image.NewName(fmt.Sprintf("%s:%s", strings.Split(r.targetRef, ":")[0], strings.Replace(i.Image, ":", "-", -1)))
	if err != nil {
		return false, fmt.Errorf("canot get new image name: %v", err)
	}
	// tag/push the image to its new repository
	dig, err := r.client.Copy(n, nn)
	if err != nil {
		return false, err
	}
	if dig.String() != i.Digest {
		i.Digest = dig.String()
		i.OriginalImage = i.Image
		i.Image = nn.String()
		return false, nil
	}
	// if i.Digest != "" && dig.String() != i.Digest {
	// 	// should not happen
	// 	return false, fmt.Errorf("digest of image %s not preserved: old digest %s; new digest %s", i.Image, i.Digest, dig.String())
	// }

	// update the imagemap
	i.Digest = dig.String()
	i.OriginalImage = i.Image
	i.Image = nn.String()
	return true, nil
}

func isOCI(imageType string) bool {
	return imageType == "" || imageType == "oci"
}

func isDocker(imageType string) bool {
	return imageType == "docker"
}

// func (r *relocateCmd) relocateImage(i *bundle.BaseImage) error {
// 	n, err := image.NewName(r.inputBundle)
// 	if err != nil {
// 		return fmt.Errorf("cannot original image fully qualified name: %v", err)
// 	}

// 	fmt.Printf("original name: %v", n.Name())

// 	nn, err := image.NewName(r.targetRef)
// 	if err != nil {
// 		return fmt.Errorf("canot get new image name: %v", err)
// 	}

// 	fmt.Printf("new name: %v", nn.Name())

// 	digest, err := r.client.Copy(n, nn)
// 	if err != nil {
// 		return fmt.Errorf("cannot copy source %v into destination %v: %v", n, r.targetRef, err)
// 	}

// 	fmt.Printf("digest: %v", digest.String())

// 	return nil
//}
