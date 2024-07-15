package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/go-redis/redis"
)

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	}) //todo
	return client
}

var mutex sync.Mutex

func main() {
	//todo complete command arguments processing
	if len(os.Args) != 2 {
		fmt.Println("cygnudge-server: invalid argumants")
		os.Exit(1)
	}
	initVals()
	listener, err := net.Listen("tcp", os.Args[1])
	if err != nil {
		log.Fatalln(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalln(err)
			continue
		}
		go process(conn)
	}
}
