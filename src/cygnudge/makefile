all: server_ client_

server_: server/server.go
	cd ./server&&go build -o ../build/server .

client_: client/client.go
	cd ./client&&go build -o ../build/client .

tidy:
	go mod tidy
	cd ./server&&go mod tidy
	cd ./client&&go mod tidy

clean:
	rm -f ./build/*
