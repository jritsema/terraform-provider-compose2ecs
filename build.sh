#!/bin/sh
set -e

go vet -x $(go list ./... | grep -v /vendor/)

godep go test -v $(go list ./... | grep -v /vendor/)

gox -osarch="darwin/amd64" -osarch="linux/amd64" -osarch="windows/amd64" -output "dist/ncd_{{.OS}}_{{.Arch}}"
