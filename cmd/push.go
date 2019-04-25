package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/containerd/containerd/remotes"
	"github.com/containerd/containerd/remotes/docker"
	auth "github.com/deislabs/oras/pkg/auth/docker"
	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
	"github.com/spf13/cobra"
)

type pushCmd struct {
	inputBundle string
	targetRef   string
}

// CNABMediaType represents a *temporary* media type for thin CNAB bundles
// it is not final, and currently acts as a placeholder
const CNABMediaType = "application/vnd.cnab.bundle.thin.v1-wd+json"

func newPushCmd() *cobra.Command {
	const usage = `pushes a CNAB bundle to a registry using ORAS`

	var p pushCmd

	cmd := &cobra.Command{
		Use:   "push",
		Short: usage,
		Long:  usage,
		RunE: func(cmd *cobra.Command, args []string) error {
			p.inputBundle = args[0]
			return p.run()
		},
	}

	cmd.Flags().StringVarP(&p.targetRef, "target", "t", "", "reference where the bundle will be pushed")

	return cmd
}

func (p *pushCmd) run() error {

	b, err := ioutil.ReadFile(p.inputBundle)
	if err != nil {
		return fmt.Errorf("cannot read bundle content: %v", err)
	}

	ms := content.NewMemoryStore()
	desc := ms.Add(p.inputBundle, CNABMediaType, b)
	pushContents := []ocispec.Descriptor{desc}
	resolver := newResolver()
	_, err = oras.Push(context.Background(), resolver, p.targetRef, ms, pushContents)

	return err
}

func newResolver() remotes.Resolver {
	cli, err := auth.NewClient()
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading auth file: %v\n", err)
	}
	resolver, err := cli.Resolver(context.Background())
	if err != nil {
		fmt.Fprintf(os.Stderr, "WARNING: Error loading resolver: %v\n", err)
		resolver = docker.NewResolver(docker.ResolverOptions{})
	}
	return resolver
}
