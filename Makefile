all:
	GOOS=linux GOARCH=arm64 go build -o ~/bin/revel -ldflags=all="-s -w" ./revel
