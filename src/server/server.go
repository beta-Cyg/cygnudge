package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"log"
)

func receive_file(conn net.Conn,file_name string){
	f,err:=os.Create(file_name)
	if err!=nil{
		log.Fatal(err)
		return
	}
	defer f.Close()
	buf:=make([]byte,4096)
	for{
		n,err:=conn.Read(buf)
		if err!=nil && err!=io.EOF{
			log.Panicf("Error from %v: %v",conn.RemoteAddr(),err)
			break
		}
		if n==0{
			log.Printf("receive file %s from %v successfully\n",file_name,conn.RemoteAddr())
			break
		}
		f.Write(buf[:n])
	}
}

func process(conn net.Conn){
	defer conn.Close()
	buf:=make([]byte,4096)
	n,err:=conn.Read(buf)
	if err!=nil{
		log.Panic(err)
		return
	}
	file_name:=string(buf[:n])
	conn.Write([]byte("ok"))
	receive_file(conn,file_name)
}

func main(){
	args:=os.Args
	if len(args)!=2{
		fmt.Println("Usage: server <listen ip&port>")
		return
	}
	listen,err:=net.Listen("tcp",args[1])
	if err!=nil{
		log.Fatal(err)
		return
	}
	defer listen.Close()
	for{
		conn,err:=listen.Accept()
		if err!=nil{
			log.Panic(err)
			continue
		}
		go process(conn)
	}
}
