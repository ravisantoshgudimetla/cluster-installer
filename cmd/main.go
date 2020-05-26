package main

import (
	"flag"
	"log"

	"github.com/ravisantoshgudimetla/cluster-installer/pkg/goinstaller"
)

func main() {
	var platform, region, name, publicKeyPath, pullSecretPath, installerPath string
	var skipDownload bool
	flag.StringVar(&platform, "platform", "aws", "specify the platform to use, defaults to aws")
	flag.StringVar(&region, "region", "us-east-2", "specify the region to use, defaults to us-east-2 for aws")
	flag.StringVar(&name, "name", "my-dev-cluster",
		"specify the name of cluster to be created, defaults to my-dev-cluster")
	flag.StringVar(&publicKeyPath, "key", "", "path to public key to be used, defaults to empty string be careful")
	flag.StringVar(&pullSecretPath, "pull-secret", "", "path to pull secret, be careful")
	flag.BoolVar(&skipDownload, "skip-download", false, "path to pull secret, be careful")
	flag.StringVar(&installerPath, "installer-path", "/tmp", "installation directory")
	flag.Parse()

	iOpts := goinstaller.NewInstallerOptions(platform, region, name, pullSecretPath, publicKeyPath, installerPath,
		skipDownload)
	//iOpts := goinstaller.NewInstallerOptions()
	if err := iOpts.Validate(); err != nil {
		log.Fatalf("error validating installer options: %v", err)
	}
	if err := iOpts.RunInstaller(); err != nil {
		log.Fatalf("error running installer: %v", err)
	}
}
