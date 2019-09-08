all: linux darwin

darwin: dirs
	GOOS=darwin GOARCH=amd64 go build -o bin/darwin/member-app .

linux: dirs
	GOOS=linux GOARCH=amd64 go build -o bin/linux/member-app .

dirs:
	mkdir -p bin/linux
	mkdir -p bin/darwin
