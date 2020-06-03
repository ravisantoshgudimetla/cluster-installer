package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var version = "0.0.1"

func init() {
	versionCmd := newVersionCmd()
	ociCmd.AddCommand(versionCmd)
}

func newVersionCmd() *cobra.Command {
	versionCmd := &cobra.Command{
		Use:              "version",
		Short:            "Specifies the version of oci",
		Long:             "Specifies the version of oci",
		TraverseChildren: true,
		Run: func(_ *cobra.Command, args []string) {
			fmt.Printf("version: %s\n", version)
		},
	}
	return versionCmd
}
