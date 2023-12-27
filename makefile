CXXFLAGS=-std=c++20 -static -DCYGNUDGE_DEBUG

all: client server

client: src/client/client.go
	go build -o bin/client src/client/client.go

server: src/server/server.go
	go build -o bin/server src/server/server.go

cygnudge: src/client/cygnudge.go src/server/cygnudge.go
	go build -o bin/cygnudge src/client/cygnudge.go
	go build -o bin/cygnudge-server src/server/cygnudge.go

interface: src/server/judge_task.hpp src/server/judge_interface.h src/server/judge_interface.cpp
	$(CXX) $(CXXFLAGS) -c src/server/judge_interface.cpp -o lib/judge_interface.o
	ar rvs lib/libjudge_interface.a lib/judge_interface.o

cgo_test: src/test/cgot.go
	export CGO_ENABLED=1
	#go clean -cache
	cd src/test&&go build -o ../../bin/test/cgo_test cgot.go

clean:
	rm -f bin/client bin/server bin/test/cgo_test bin/cygnudge bin/cygnudge-server lib/*
