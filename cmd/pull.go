package main

import (
	"context"
	"fmt"

	"github.com/deislabs/oras/pkg/content"
	"github.com/deislabs/oras/pkg/oras"
	"github.com/spf13/cobra"
)

type pullCmd struct {
	outputBundle string
	targetRef    string
}

func newPullCmd() *cobra.Command {
	const usage = "pulls a CNAB bundle from a registry using ORAS"

	var p pullCmd

	cmd := &cobra.Command{
		Use:   "pull",
		Short: usage,
		Long:  usage,
		RunE: func(cmd *cobra.Command, args []string) error {
			p.outputBundle = args[0]
			return p.run()
		},
	}

	cmd.Flags().StringVarP(&p.targetRef, "target", "t", "", "reference where the bundle will be pushed")

	return cmd
}

func (p *pullCmd) run() error {

	fs := content.NewFileStore(p.outputBundle)
	defer fs.Close()

	desc, layers, err := oras.Pull(context.Background(), newResolver(), p.targetRef, fs, oras.WithAllowedMediaTypes([]string{CNABMediaType}))
	if err != nil {
		return fmt.Errorf("cannot pull bundle: %v", err)
	}

	fmt.Printf("descriptor: %v\n\n layers: %v", desc, layers)

	return nil
}
