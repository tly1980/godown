install:
	go install godown

build:
	go build godown

all: build

ut:
	pushd worker; go test -v; popd

deps:
	go get github.com/stretchr/testify/assert
	go get github.com/golang/glog