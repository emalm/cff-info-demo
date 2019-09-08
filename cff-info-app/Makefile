all: linux darwin

linux:
	GOOS=linux GOARCH=amd64 go build -o cff-info-app.linux .

darwin:
	GOOS=darwin GOARCH=amd64 go build -o cff-info-app.darwin .

runlocal: darwin
	./cff-info-app.darwin