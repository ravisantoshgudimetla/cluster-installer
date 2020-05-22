# cluster-installer
Cluster installer for OpenShift

## Why another wrapper for installer
To just expose the minimal config needed for the OVNKubernetes network variant to work. If some wants to modify the 
config, they should update the bindata.

## Why go
I initially started out in python but I wanted types so I was left with go or Rust. I'll use Rust perhaps for my
next project :P
