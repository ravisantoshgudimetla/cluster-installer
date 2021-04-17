package cmd

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"sync"

	"github.com/ravisantoshgudimetla/cluster-installer/pkg/goinstaller"
	"github.com/spf13/cobra"
)

var deleteInfo struct {
	deleteInstaller bool
	platform        string
}

func init() {
	deleteCmd := newDeleteCmd()
	ociCmd.AddCommand(deleteCmd)
}

func newDeleteCmd() *cobra.Command {
	deleteCmd := &cobra.Command{
		Use:              "delete",
		Short:            "Deletes OpenShift cluster",
		Long:             "deletes an OpenShift cluster",
		TraverseChildren: true,
		RunE: func(_ *cobra.Command, args []string) error {
			return deleteInstaller(ociInfo.installerPath, deleteInfo.platform, deleteInfo.deleteInstaller)
		},
	}
	deleteCmd.PersistentFlags().BoolVar(&deleteInfo.deleteInstaller, "delete-installer", false,
		"deletes the installer binary")
	deleteCmd.PersistentFlags().StringVar(&deleteInfo.platform, "platform", "",
		"deletes the specific directory after installation is complete")
	return deleteCmd
}

func deleteInstaller(installerPath, platform string, deleteInstaller bool) error {
	var errStdout, errStderr error
	var osType string
	if runtime.GOOS == "darwin" {
		osType = "mac"
	} else if runtime.GOOS == "linux" {
		osType = "linux"
	}
	installDirectory := installerPath + "/openshift-install-"+ osType + "/"+ platform
	cmd := exec.Command(installerPath+"/openshift-install-"+osType+"/openshift-install", "destroy", "cluster", "--log-level=debug",
		"--dir="+installDirectory)
	stdoutIn, _ := cmd.StdoutPipe()
	stderrIn, _ := cmd.StderrPipe()
	stdout := goinstaller.NewCapturingPassThroughWriter(os.Stdout)
	stderr := goinstaller.NewCapturingPassThroughWriter(os.Stderr)
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
	if deleteInstaller {
		// delete the directory containing
		if err := os.RemoveAll(installerPath + "/openshift-install-"+osType); err != nil {
			return fmt.Errorf("error deleting file %v", err)
		}
		// delete the installer tar file we downloaded earlier
		if err := os.RemoveAll(installerPath + "/openshift-install-"+osType+".tar.gz"); err != nil {
			return fmt.Errorf("error deleting installer tar file %v", err)
		}
	}
	return nil
}
