package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ociInfo struct {
		installerPath string
		skipDownload  bool
	}
	ociCmd = &cobra.Command{
		Use:   "oci",
		Short: "downloads and runs the OpenShift installer binary to spin up OpenShift cluster",
	}
)

func init() {
	ociCmd.PersistentFlags().StringVar(&ociInfo.installerPath, "installer-path", "/tmp", "installation directory")
	ociCmd.PersistentFlags().BoolVar(&ociInfo.skipDownload, "skip-download", false, "skips the download of the openshift install file")
}

// Execute run the actual oci command
func Execute() {
	if err := ociCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
