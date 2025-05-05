all:
	go build -o ~/bin/revel -ldflags=all="-s -w" ./revel
build:
	go build -o ~/bin/revel -gcflags=all="-N -l" ./revel
fmt:
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -local github.com/wiselike/revel-cmd -l -w . \
	else \
		go fmt ./... \
	fi
