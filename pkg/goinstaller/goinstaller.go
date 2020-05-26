package goinstaller

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
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
	installerPath string, skipDownload bool) *installerOptions {
	return &installerOptions{
		platform:       platform,
		region:         region,
		name:           name,
		pullSecretPath: pullSecretPath,
		publicKeyPath:  publicKeyPath,
		installerPath:  installerPath,
		skipDownload:   skipDownload,
	}
}

// Validate validates the code
func (i *installerOptions) Validate() error {
	return nil
}

func (i *installerOptions) RunInstaller() error {
	// If we want to skip download, don't download the installer
	if !i.skipDownload {
		log.Printf("Downloading the installer at path at %s", i.installerPath)
		if err := i.downloadInstallerBinary(); err != nil {
			return err
		}
	}
	installDirectory := i.installerPath + "/openshift-install-linux/" + i.platform
	if err := os.MkdirAll(installDirectory, 0744); err != nil {
		return err
	}
	log.Printf("Successfully created install directory at %s", installDirectory)
	if err := i.writeInstallConfig(installDirectory); err != nil {
		return err
	}
	log.Printf("Successfully created install config at %s", installDirectory)
	// Generate manifests
	if err := i.writeInstallManifests(installDirectory); err != nil {
		return err
	}
	manifestsDirectory := installDirectory + "/" + "manifests/"
	log.Printf("Successfully created install manifest at %s", manifestsDirectory)
	if err := writeClusterNetworkFile(manifestsDirectory); err != nil {
		return err
	}
	log.Printf("Successfully created cluster network file at %s", manifestsDirectory)
	if err := i.runInstaller(installDirectory); err != nil {
		return err
	}
	log.Print("Successfully ran installer")
	return nil
}

func (i *installerOptions) downloadInstallerBinary() error {
	resp, err := http.Get("http://mirror.openshift.com/pub/openshift-v4/clients/ocp-dev-preview/latest/openshift-install-linux.tar.gz")
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	file, err := os.Create(i.installerPath + "/openshift-install-linux.tar.gz")
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)

	downloadFile, err := os.Open(i.installerPath + "/openshift-install-linux.tar.gz")
	if err != nil {
		return err
	}
	destString := i.installerPath + "/" + "openshift-install-linux"
	if err := os.Mkdir(destString, 0744); err != nil {
		return err
	}
	if err := untar(destString, downloadFile); err != nil {
		return err
	}
	return err

}

func untar(dst string, r io.Reader) error {
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()

		switch {

		// if no more files are found return
		case err == io.EOF:
			return nil

		// return any other error
		case err != nil:
			return err

		// if the header is nil, just skip it (not sure how this happens)
		case header == nil:
			continue
		}

		// the target location where the dir/file should be created
		target := filepath.Join(dst, header.Name)

		// the following switch could also be done using fi.Mode(), not sure if there
		// a benefit of using one vs. the other.
		// fi := header.FileInfo()

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
	}
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
	cmd := exec.Command(i.installerPath+"/openshift-install-linux/"+"openshift-install", "create", "cluster", "--log-level=debug",
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
	cmd := exec.Command(i.installerPath+"/openshift-install-linux/"+"openshift-install", "create", "manifests",
		"--dir="+installDirectory)
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
