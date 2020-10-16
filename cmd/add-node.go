package cmd

import (
	"fmt"

	"github.com/ravisantoshgudimetla/cluster-installer/pkg/goinstaller"
	"github.com/spf13/cobra"
)

var addNode struct {
	nodeCount int32
	machineSetName string
	imageName string
}

func init() {
	addNodeCommand := newAddCmd()
	ociCmd.AddCommand(addNodeCommand)
}

func newAddCmd() *cobra.Command {
	createCmd := &cobra.Command{
		Use:              "add",
		Short:            "Adds a new windows node to the cluster",
		Long:             "Adds a new windows node to the cluster",
		TraverseChildren: true,
		RunE: func(_ *cobra.Command, args []string) error {
			iOpts := goinstaller.NewInstallerOptions(createInfo.platform, createInfo.region, createInfo.name,
				createInfo.pullSecretPath, createInfo.publicKeyPath, ociInfo.installerPath, ociInfo.skipDownload)
			if err := iOpts.Validate(); err != nil {
				return fmt.Errorf("error validating installer options: %v", err)
			}
			if err := iOpts.RunInstaller(); err != nil {
				return fmt.Errorf("error running installer: %v", err)
			}
			return nil
		},
	}
	createCmd.PersistentFlags().StringVar(&createInfo.platform, "platform", "aws",
		"specify the platform to use, defaults to aws")
	createCmd.PersistentFlags().StringVar(&createInfo.region, "region", "us-east-2",
		"specify the region to use, defaults to us-east-2 for aws")
	createCmd.PersistentFlags().StringVar(&createInfo.name, "name", "my-dev-cluster",
		"specify the name of cluster to be created, defaults to my-dev-cluster")
	createCmd.PersistentFlags().StringVar(&createInfo.publicKeyPath, "key", "",
		"path to public key to be used, defaults to empty string be careful")
	createCmd.PersistentFlags().StringVar(&createInfo.pullSecretPath, "pull-secret", "",
		"path to pull secret, be careful")
	return createCmd
}
