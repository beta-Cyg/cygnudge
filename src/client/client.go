package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"log"
	"flag"
)

func send_file(conn net.Conn,file_path string){
	f,err:=os.Open(file_path)
	if err!=nil{
		log.Fatal(err)
		return
	}
	defer f.Close()
	buf:=make([]byte,4096)
	for{
		n,err:=f.Read(buf)
		if err!=nil{
			if err==io.EOF{
				log.Printf("file %s has been completely read\n",file_path)
			}else{
				log.Fatal(err)
			}
			return
		}
		_,err=conn.Write(buf[:n])
		if err!=nil{
			log.Fatal(err)
			return
		}
	}
}

func main(){
	args:=os.Args
	if len(args)!=3{
		fmt.Println("Usage: client <local_file_path> <remote_ip&port>")
		return
	}
	file_path:=args[1]
	file_info,err:=os.Stat(file_path)
	if err!=nil{
		log.Fatal(err)
		return
	}
	file_name:=file_info.Name()
	conn,err:=net.Dial("tcp",args[2])
	if err!=nil{
		log.Fatal(err)
		return
	}
	defer conn.Close()
	_,err=conn.Write([]byte(file_name))
	if err!=nil{
		log.Fatal(err)
		return
	}
	buf:=make([]byte,16)
	n,err:=conn.Read(buf)
	if err!=nil{
		log.Fatal(err)
		return
	}
	if string(buf[:n])=="ok"{
		send_file(conn,file_path)
	}
}
