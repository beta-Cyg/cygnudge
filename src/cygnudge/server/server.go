package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"

	"github.com/go-redis/redis"
)

var redis_server_address string

func newClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     redis_server_address,
		Password: "",
		DB:       0,
	}) //todo
	return client
}

var mutex sync.Mutex

func main() {
	//todo complete command arguments processing
	if len(os.Args) != 3 {
		fmt.Println("cygnudge-server: invalid argumants")
		fmt.Println("command usage: ")
		fmt.Printf("%v {listened address & port} {redis server address & port}\n", os.Args[0])
		os.Exit(1)
	}
	redis_server_address = os.Args[2]
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
