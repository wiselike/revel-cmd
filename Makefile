all:
	go build -o ~/bin/revel -ldflags=all="-s -w" ./revel
build:
	go build -o ~/bin/revel -gcflags=all="-N -l" ./revel
fmt:
	go fmt ./...
