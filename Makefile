all:
	GOOS=linux go build -o ~/bin/revel -ldflags=all="-s -w" ./revel
