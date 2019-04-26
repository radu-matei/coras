package main

import (
	"fmt"
	"io/ioutil"

	"github.com/radu-matei/coras/pkg/coras"

	"github.com/deislabs/cnab-go/bundle"
	"github.com/spf13/cobra"
)

type pushCmd struct {
	inputBundle string
	targetRef   string
	exported    bool
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
			p.targetRef = args[1]
			return p.run()
		},
	}

	cmd.Flags().BoolVarP(&p.exported, "exported", "", false, "When passed, this command will push an exported (thick) bundle")
	return cmd
}

func (p *pushCmd) run() error {
	if p.exported {
		return coras.PushThick(p.inputBundle, p.targetRef)
	}

	data, err := ioutil.ReadFile(p.inputBundle)
	if err != nil {
		return fmt.Errorf("cannot read input bundle: %v", err)
	}
	b, err := bundle.Unmarshal(data)
	if err != nil {
		return fmt.Errorf("cannot unmarshal input bundle: %v", err)
	}

	return coras.PushThin(b, p.targetRef)
}
