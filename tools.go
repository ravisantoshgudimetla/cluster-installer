// +build tools

// tools is a dummy package that will be ignored for builds, but included for dependencies.
package tools

import (
	_ "github.com/go-bindata/go-bindata"
	_ "github.com/go-bindata/go-bindata/go-bindata"
)
