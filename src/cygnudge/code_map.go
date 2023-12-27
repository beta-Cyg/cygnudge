package cygnudge

import (
	"fmt"
	"net"
	"log"
)

var CodeRequestMap=map[string]string{
	"001" : "Register",
	"002" : "Login",
	"003" : "Check Login",
	"004" : "Get Uid",
	"005" : "Logout",
	"010" : "Submit",
}

var RequestCodeMap=map[string]string{
	"Register" : "001",
	"Login" : "002",
	"Check Login" : "003",
	"Get Uid" : "004",
	"Logout" : "005",
	"Submit" : "010",
}

var CodeResponseMap=map[string]string{
	"100" : "OK",
	"110" : "Bad Request",
	"120" : "Not Acceptable",
	"121" : "Remote Login",
}

var ResponseCodeMap=map[string]string{
	"OK" : "100",
	"Bad Request" : "110",
	"Not Acceptable" : "120",
	"Remote Login" : "121",
}

func GetReqCode(name string)(string,error){
	code,exist:=RequestCodeMap[name]
	if !exist{
		return "",fmt.Errorf("unknown request message: %s",name)
	}
	return code,nil
}

func GetReqName(code string)(string,error){
	name,exist:=CodeRequestMap[code]
	if !exist{
		return "",fmt.Errorf("unknown request code: %s",code)
	}
	return name,nil
}

func GetResCode(name string)(string,error){
	code,exist:=ResponseCodeMap[name]
	if !exist{
		return "",fmt.Errorf("unknown response name: %s",name)
	}
	return code,nil
}

func GetResName(code string)(string,error){
	name,exist:=CodeResponseMap[code]
	if !exist{
		return "",fmt.Errorf("unknown response code: %s",code)
	}
	return name,nil
}

func SendRes(name string,conn net.Conn){
	response_code,err:=GetResCode(name)
	if err!=nil{
		log.Fatalln(err)
	}
	log.Printf("SendRes: %s %s\n",response_code,name)
	_,err=conn.Write([]byte(response_code))
	if err!=nil{
		log.Fatalln(err)
	}
}

func ReceiveRes(name string,conn net.Conn){
	buf:=make([]byte,4)
	n,err:=conn.Read(buf)
	if err!=nil{
		log.Fatalln(err)
	}
	if n!=3{
		log.Fatalf("invalid response code (ReceiveRes): %s\n",string(buf[:n]))
	}
	response_code:=string(buf[:n])
	response_name,err:=GetResName(response_code)
	if err!=nil{
		log.Fatalln(err)
	}
	if response_name==name{
		log.Printf("ReceiveRes: %s %s\n",response_code,response_name)
	}else{
		log.Fatalf("ReceiveRes: %s %s\n",response_code,response_name)
	}
}

func ReceiveGetRes(conn net.Conn)(string,string){
	buf:=make([]byte,4)
	n,err:=conn.Read(buf)
	if err!=nil{
		log.Fatalln(err)
	}
	if n!=3{
		log.Fatalf("invalid response code (ReceiveGetRes): %s\n",string(buf[:n]))
	}
	response_code:=string(buf[:n])
	response_name,err:=GetResName(response_code)
	if err!=nil{
		log.Fatalln(err)
	}
	log.Printf("ReceiveGetRes: %s %s\n",response_code,response_name)
	return response_code,response_name
}

func SendReq(name string,conn net.Conn){
	request_code,err:=GetReqCode(name)
	if err!=nil{
		log.Fatalln(err)
	}
	log.Printf("SendReq: %s %s\n",request_code,name)
	_,err=conn.Write([]byte(request_code))
	if err!=nil{
		log.Fatalln(err)
	}
}

func ReceiveReq(name string,conn net.Conn){
	buf:=make([]byte,4)
	n,err:=conn.Read(buf)
	if err!=nil{
		log.Fatalln(err)
	}
	if n!=3{
		log.Fatalf("invalid request code (ReceiveReq): %s\n",string(buf[:n]))
	}
	request_code:=string(buf[:n])
	request_name,err:=GetReqName(request_code)
	if err!=nil{
		log.Fatalln(err)
	}
	if request_name==name{
		log.Printf("ReceiveReq: %s %s\n",request_code,request_name)
	}else{
		log.Fatalf("ReceiveReq: %s %s\n",request_code,request_name)
	}
}

func ReceiveGetReq(conn net.Conn)(string,string){
	buf:=make([]byte,4)
	n,err:=conn.Read(buf)
	if err!=nil{
		log.Fatalln(err)
	}
	if n!=3{
		log.Fatalf("invalid request code (ReceiveGetReq): %s\n",string(buf[:n]))
	}
	request_code:=string(buf[:n])
	request_name,err:=GetReqName(request_code)
	if(err!=nil){
		log.Fatalln(err)
	}
	log.Printf("ReceiveGetReq: %s %s\n",request_code,request_name)
	return request_code,request_name
}
