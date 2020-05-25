package goinstaller

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"sync"

	"github.com/ravisantoshgudimetla/cluster-installer/pkg/assets"
)

type installerOptions struct {
	platform       string
	region         string
	name           string
	pullSecretPath string
	publicKeyPath  string
	installerPath  string
	skipDownload   bool
}

func NewInstallerOptions(platform, region, name, pullSecretPath, publicKeyPath,
	installerPath string) *installerOptions {
	return &installerOptions{
		platform:       platform,
		region:         region,
		name:           name,
		pullSecretPath: pullSecretPath,
		publicKeyPath:  publicKeyPath,
		installerPath:  installerPath,
	}
}

// Validate validates the code
func (i *installerOptions) Validate() error {
	return nil
}

func (i *installerOptions) RunInstaller() error {
	// Create the Installer directory. If skipDownload is set, we'll skip the download of installer.
	if _, err := os.Stat("/home/ravig/Downloads/openshift-install-linux"); os.IsNotExist(err) {
		// File doesn't exist, so download again.
	}
	installDirectory := i.installerPath + "/" + i.platform
	if err := os.MkdirAll(installDirectory, 0744); err != nil {
		return err
	}
	if err := i.writeInstallConfig(installDirectory); err != nil {
		return err
	}
	// Generate manifests
	if err := i.writeInstallManifests(installDirectory); err != nil {
		return err
	}
	manifestsDirectory := installDirectory + "/" + "manifests/"
	if err := writeClusterNetworkFile(manifestsDirectory); err != nil {
		return err
	}
	// Installer code is working, don't run it for now
	// if err := i.runInstaller(installDirectory); err != nil {
	// 	return err
	// }
	return downloadInstallerBinary()
}

// CapturingPassThroughWriter is a writer that remembers
// data written to it and passes it to w
type CapturingPassThroughWriter struct {
	buf bytes.Buffer
	w   io.Writer
}

// NewCapturingPassThroughWriter creates new CapturingPassThroughWriter
func NewCapturingPassThroughWriter(w io.Writer) *CapturingPassThroughWriter {
	return &CapturingPassThroughWriter{
		w: w,
	}
}

func (w *CapturingPassThroughWriter) Write(d []byte) (int, error) {
	w.buf.Write(d)
	return w.w.Write(d)
}

// Bytes returns bytes written to the writer
func (w *CapturingPassThroughWriter) Bytes() []byte {
	return w.buf.Bytes()
}

// runInstaller runs the openshif installer and captures the output.
// Source code copied from: https://blog.kowalczyk.info/article/wOYk/advanced-command-execution-in-go-with-osexec.html
func (i *installerOptions) runInstaller(installDirectory string) error {
	var errStdout, errStderr error
	cmd := exec.Command(i.installerPath+"/"+"openshift-install", "create", "cluster", "--log-level=debug",
		"--dir="+installDirectory)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	stdout := NewCapturingPassThroughWriter(os.Stdout)
	stderr := NewCapturingPassThroughWriter(os.Stderr)
	err := cmd.Start()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		_, errStdout = io.Copy(stdout, stdoutIn)
		wg.Done()
	}()

	_, errStderr = io.Copy(stderr, stderrIn)
	wg.Wait()

	err = cmd.Wait()
	if err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
	if errStdout != nil || errStderr != nil {
		log.Fatalf("failed to capture stdout or stderr\n")
	}
	outStr, errStr := string(stdout.Bytes()), string(stderr.Bytes())
	fmt.Printf("\nout:\n%s\nerr:\n%s\n", outStr, errStr)
	return nil
}

func (i *installerOptions) writeInstallManifests(installDirectory string) error {
	cmd := exec.Command(i.installerPath+"/"+"openshift-install", "create", "manifests", "--dir="+installDirectory)
	var out bytes.Buffer
	cmd.Stdout = os.Stdout
	err := cmd.Run()
	if err != nil {
		return err
	}
	fmt.Println(out)
	return nil
}

// writeInstallConfig spits out the install config after reading values from the config file
func (i *installerOptions) writeInstallConfig(installDirectory string) error {
	// TODO: Create a map with with options
	//optionsToBeReplaced := map[string]string{ "platform": "aws"}
	publicKey, err := ioutil.ReadFile(i.publicKeyPath)
	if err != nil {
		return err
	}
	pullSecret, err := ioutil.ReadFile(i.pullSecretPath)
	if err != nil {
		return err
	}
	var installConfig []byte
	if i.platform == "aws" {
		installConfig = assets.MustAsset("generated/install-config-aws.yaml")
	} else if i.platform == "azure" {
		installConfig = assets.MustAsset("generated/install-config-azure.yaml")
	}

	installConfigToBeCreated := strings.Replace(string(installConfig), "<sshKey>", string(publicKey), -1)
	installConfigToBeCreated = strings.Replace(installConfigToBeCreated, "<pullSecret>",
		strings.TrimSuffix(string(pullSecret), "\n"), -1)
	installConfigToBeCreated = strings.Replace(installConfigToBeCreated, "<platform>", i.platform, -1)
	installConfigToBeCreated = strings.Replace(installConfigToBeCreated, "<region>", i.region, -1)
	installConfigToBeCreated = strings.Replace(installConfigToBeCreated, "<name>", i.name, -1)
	if err = ioutil.WriteFile(installDirectory+"/"+"install-config.yaml", []byte(installConfigToBeCreated), 0744); err != nil {
		return err
	}
	return nil
}

// writeClusterNetworkFile writes the Cluster-network-config file to stdout
func writeClusterNetworkFile(manifestsDirectory string) error {
	networkConfig := assets.MustAsset("generated/cluster-network-03-config.yaml")
	if err := ioutil.WriteFile(manifestsDirectory+"cluster-network-03-config.yaml", networkConfig, 0640); err != nil {
		return err
	}
	return nil
}

func downloadInstallerBinary() error {
	return nil

}
