run:
	go build -o bin/api ./cmd/api
	bin/api

linux:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o bin/api_linux ./cmd/api
