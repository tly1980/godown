install:
	go install godown

build:
	go build godown

all: build

ut:
	pushd worker; go test; popd