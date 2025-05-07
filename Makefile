VERSION := $(shell git describe --tags --match 'v*' --abbrev=0 2>/dev/null || echo "v0.0.0")
BUILD_DATE := $(shell date -u +%Y-%m-%d)

release:
	CGO_ENABLED=0 go build -o ~/bin/revel -ldflags=all="-s -w" -ldflags "-X github.com/wiselike/revel-cmd.Version=$(VERSION) -X github.com/wiselike/revel-cmd.BuildDate=$(BUILD_DATE)" ./revel
build:
	CGO_ENABLED=0 go build -o ~/bin/revel -gcflags=all="-N -l" -ldflags "-X github.com/wiselike/revel-cmd.Version=$(VERSION) -X github.com/wiselike/revel-cmd.BuildDate=$(BUILD_DATE)" ./revel
fmt:
	@if command -v goimports >/dev/null 2>&1; then \
		goimports -local github.com/wiselike/revel-cmd -l -w . \
	else \
		go fmt ./... \
	fi
test: export PATH := ~/bin:$(PATH)
test: export REVEL_BRANCH="develop"
test: build
	go test ./...
	
	@rm -rf my build target
	
	# Ensure the new-app flow works (plus the other commands).
	revel version
	revel new     my/testapp
	revel test    my/testapp
	revel clean   my/testapp
	revel build   my/testapp build/testapp
	revel build   my/testapp build/testapp prod
	revel package my/testapp
	revel package my/testapp prod
	
	# Ensure the new-app flow works (plus the other commands).
	revel new     --gomod-flags "edit -replace=github.com/wiselike/revel=github.com/wiselike/revel@$$REVEL_BRANCH" -a my/testapp2 --package revelframework.com -v
	revel test  -a my/testapp2 -v
	revel clean -a my/testapp2 -v
	revel build -a my/testapp2 -v -t build/testapp2
	revel build -a my/testapp2 -v -t build/testapp2 -m prod
	revel package -a my/testapp2 -v
	revel package -a my/testapp2 -v -m prod
	
	# Check build works with no-vendor flag, only go versions v1.12 or older supported
	# revel new  -a my/testapp3 --no-vendor -v
	# revel test -a my/testapp3 -v
	
	# Check vendored version of revel
	revel new     --gomod-flags "edit -replace=github.com/wiselike/revel=github.com/wiselike/revel@$$REVEL_BRANCH" -a my/testapp4 --package revelframework.com
	revel build -a my/testapp4; # build first to auto update go.mod file
	cd my/testapp4; go mod vendor
	revel test -a my/testapp4

