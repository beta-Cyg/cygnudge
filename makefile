CXXFLAGS=-std=c++20 -static -DCYGNUDGE_DEBUG

all: cygnudge interface cgo_test

cygnudge: src/cygnudge/server/server.go src/cygnudge/client/client.go
	cd src/cygnudge&&make

interface: src/server/judge_task.hpp src/server/judge_interface.h src/server/judge_interface.cpp
	$(CXX) $(CXXFLAGS) -c src/server/judge_interface.cpp -o lib/judge_interface.o
	ar rvs lib/libjudge_interface.a lib/judge_interface.o
	$(CXX) -fPIC $(CXXFLAGS) -c src/server/judge_interface.cpp -o lib/cygnudge_judge.o
	$(CXX) -shared -o lib/libcygnudge_judge.so lib/cygnudge_judge.o

cgo_test: src/test/go/cgot.go
	export CGO_ENABLED=1
	cd src/test&&make

clean:
	cd src/cygnudge&&make clean
	cd src/test&&make clean
	rm -f src/test/go/test bin/cygnudge bin/cygnudge-server lib/*
