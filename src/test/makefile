GO = go

cgot: ./go/cgot.go ../../lib/libcygnudge_judge.so
	cd ./go && $(GO) env -w CGO_CXXFLAGS='-std=c++20 -O2 -g' && $(GO) build -o cgo_test cgot.go

pack_test: ./go/pack_test.go ../server/go/pack.go
	cd ./go && $(GO) build -o pack_test pack_test.go

clean:
	rm -f ./go/cgo_test
