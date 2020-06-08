all: build update

PACKAGE=github.com/ravisantoshgudimetla/cluster-installer

GO_BUILD_ARGS=CGO_ENABLED=0 GO111MODULE=on

.PHONY: build
build:
	$(GO_BUILD_ARGS) go build -o oci $(PACKAGE)


.PHONY: update-bindata
update-bindata: 
	hack/update-generated-bindata.sh
