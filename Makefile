install:
	go install godown

build:
	go build godown

all: build

ut:
	pushd worker; go test -stderrthreshold=INFO; popd

deps:
	go get github.com/stretchr/testify/assert
	go get github.com/golang/glog