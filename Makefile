ifndef GOPATH
	GOPATH := ${PWD}/gopath
	export GOPATH
endif

COMMITHASH := $(shell git describe --always --tags --dirty)

pcd: *.go
	go get -v -d
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags "-s -X main.COMMITHASH=${COMMITHASH}" -o pcd

.PHONY: clean
clean:
	rm -f pcd
	rm -f go-pcd
