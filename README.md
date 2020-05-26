# cluster-installer
Cluster installer for OpenShift

## Why another wrapper for installer
To just expose the minimal config needed for the OVNKubernetes network variant to work. If some wants to modify the 
config, they should update the bindata.

## Why go
I initially started out in python but I wanted types so I was left with go or Rust. I'll use Rust perhaps for my
next project.

## How to run
```go
 ./oci -installer-path /tmp/new_installer -key <use-openshift-dev> -name my-cool-cluster -platform azure -region centralus -pull-secret pull-secret 
```

Options
```go
➜  cluster-installer git:(master) ✗ ./oci --help  
Usage of ./oci:
  -installer-path string
        installation directory (default "/tmp")
  -key string
        path to public key to be used, defaults to empty string be careful
  -name string
        specify the name of cluster to be created, defaults to my-dev-cluster (default "my-dev-cluster")
  -platform string
        specify the platform to use, defaults to aws (default "aws")
  -pull-secret string
        path to pull secret, be careful
  -region string
        specify the region to use, defaults to us-east-2 for aws (default "us-east-2")
  -skip-download
        path to pull secret, be careful
  -version
        version (default true)
```

## How to update the cluster config or network config
If someone wants to update network config or cluster network. They should update files [here](https://github.com/ravisantoshgudimetla/cluster-installer/tree/master/generated) 
and run
``` shell
make update-bindata
```