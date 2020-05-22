package main

import (
	"flag"
	"fmt"
	"os"
)

func downloadInstallerBinary() error {

	return nil

}

func main() {
	// Download the installer binary from mirror.openshift.com. URL to download
	// from:
	// http://mirror.openshift.com/pub/openshift-v4/clients/ocp-dev-preview/latest/
	// Following flags need to present
	// - platform
	// - region
	// - name
	// - pull-secret-path
	// Put bin data code from standard install-config and replace it with latest
	//    after reading from flags and generate the config file
	// Generate the manifest files from config.
	// Copy the cluster-network-operator-03.yaml to the manifests directory
	// Run the installer again. This should complete the setup
	// initFlags()
	flag.Parse()
	err := downloadInstallerBinary()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	// Read assets install-config.yaml

}
