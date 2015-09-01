ifndef GOPATH
	GOPATH := ${PWD}/gopath
	export GOPATH
endif

pcd-api: *.go
	go get -v -d
	CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-s' -o pcd
