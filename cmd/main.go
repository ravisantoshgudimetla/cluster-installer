package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/ravisantoshgudimetla/cluster-installer/pkg/goinstaller"
)

func main() {
	// Download the installer binary from mirror.openshift.com. URL to download
	// from:
	// http://mirror.openshift.com/pub/openshift-v4/clients/ocp-dev-preview/latest/
	// Following flags need to present
	// - platform
	// - region
	// - name
	// - public-key-path
	// - pull-secret-path
	// Put bin data code from standard install-config and replace it with latest
	//    after reading from flags and generate the config file
	// Generate the manifest files from config.
	// Copy the cluster-network-operator-03.yaml to the manifests directory
	// Run the installer again. This should complete the setup
	// initFlags()
	// TODO: Switch to having a root command for platform and subsequent sub commands for region
	var platform, region, name, publicKeyPath, pullSecretPath, installerPath string
	flag.StringVar(&platform, "p", "aws", "specify the platform to use, defaults to aws")
	flag.StringVar(&region, "r", "us-east-2", "specify the region to use, defaults to us-east-2 for aws")
	flag.StringVar(&name, "n", "my-dev-cluster",
		"specify the name of cluster to be created, defaults to my-dev-cluster")
	flag.StringVar(&publicKeyPath, "k", "", "path to public key to be used, defaults to empty string be careful")
	flag.StringVar(&pullSecretPath, "s", "", "path to pull secret, be careful")
	flag.StringVar(&installerPath, "i", "/tmp", "path to pull secret, be careful")
	flag.Parse()
	iOpts := goinstaller.NewInstallerOptions(platform, region, name, pullSecretPath, publicKeyPath, installerPath)
	if err := iOpts.Validate(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if err := iOpts.RunInstaller(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
