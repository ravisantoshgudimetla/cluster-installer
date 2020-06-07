# cluster-installer
Cluster installer for OpenShift. It downloads the latest installer from:
```
http://mirror.openshift.com/pub/openshift-v4/clients/ocp-dev-preview/latest/openshift-install-linux.tar.gz
```

You can also skip the download if you already happen to have a binary

## Why another wrapper for installer
To just expose the minimal config needed for the OVNKubernetes network variant to work. If some wants to modify the 
config, they should update the bindata.

## Why Go
I initially started out in python but I wanted types so I was left with Go or Rust. I want to use go to generate a single binary for the usage

## How to run
```go
./oci create --installer-path <installer_path> --key <ssh_public_key> --name my-dev-cluster --platform aws --region us-east-2 --pull-secret <pull_secret_path>
```

Options
```go
➜  cluster-installer git:(master) ✗ ./oci --help
downloads and runs the OpenShift installer binary to spin up OpenShift cluster

Usage:
  oci [command]

Available Commands:
  create      Creates OpenShift cluster
  help        Help about any command
  version     Specifies the version of oci

Flags:
  -h, --help                    help for oci
      --installer-path string   installation directory (default "/tmp")
      --skip-download           skips the download of the openshift install file

Use "oci [command] --help" for more information about a command.
```

## How to update the cluster config or network config
If someone wants to update network config or cluster network. They should update files [here](https://github.com/ravisantoshgudimetla/cluster-installer/tree/master/generated) 
and run
``` shell
make update-bindata
```