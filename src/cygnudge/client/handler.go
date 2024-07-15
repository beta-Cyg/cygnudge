package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"os/user"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
	"time"

	"cygnus.beta/cygnudge"
	"golang.org/x/term"
)

type ServerAddress struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

func register(add ServerAddress) {
	conn, err := net.Dial("tcp", add.Ip+":"+strconv.Itoa(add.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	cygnudge.SendReq("Register", conn)
	cygnudge.ReceiveRes("OK", conn)

	//input email
	var email string
	fmt.Printf("input email: ")
	fmt.Scanln(&email)
	_, err = conn.Write([]byte(email))
	if err != nil {
		log.Fatalln(err)
	}
	_, response_name := cygnudge.ReceiveGetRes(conn)
	switch response_name {
	case "OK":
		break
	case "Not Acceptable":
		fmt.Println("cygnudge register: email has been registered")
		os.Exit(2)
	}

password_input:
	//input password and encrypt it
	fmt.Printf("input password: ")
	password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalln(err)
	}
	//debug
	fmt.Printf("\npassword: %s\n", string(password))
	fmt.Printf("input password again: ")
	confirm_password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalln(err)
	}
	//debug
	fmt.Printf("\nconfirm password: %s\n", string(confirm_password))
	if !bytes.Equal(password, confirm_password) {
		goto password_input
	}
	hashed_password := fmt.Sprintf("%x", md5.Sum(password))
	//debug
	fmt.Printf("hashed password: %s\n", hashed_password)
	_, err = conn.Write([]byte(hashed_password))
	if err != nil {
		log.Fatalln(err)
	}
	cygnudge.ReceiveRes("OK", conn)

	buf := make([]byte, 16) //the length of uid is impossible to be more than 16
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	new_uid, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		log.Fatalln(err)
	}
	cygnudge.SendRes("OK", conn)
	fmt.Printf("register new user successfully: %d\n", new_uid)
}

func getUid(add ServerAddress, email string) int {
	conn, err := net.Dial("tcp", add.Ip+":"+strconv.Itoa(add.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	cygnudge.SendReq("Get Uid", conn)
	cygnudge.ReceiveRes("OK", conn)

	_, err = conn.Write([]byte(email))
	if err != nil {
		log.Fatalln(err)
	}
	_, response_name := cygnudge.ReceiveGetRes(conn)
	switch response_name {
	case "OK":
		break
	case "Not Acceptable":
		fmt.Println("cygnudge get uid: user doesn't exist")
		os.Exit(3)
	}

	//receive uid from server
	buf := make([]byte, 16)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	cygnudge.SendRes("OK", conn)

	uid, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		log.Fatalln(err)
	}
	return uid
}

func login(add ServerAddress, uid int) {
	conn, err := net.Dial("tcp", add.Ip+":"+strconv.Itoa(add.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	cygnudge.SendReq("Login", conn)
	cygnudge.ReceiveRes("OK", conn)

	_, err = conn.Write([]byte(strconv.Itoa(uid)))
	if err != nil {
		log.Fatalln(err)
	}
	_, response_name := cygnudge.ReceiveGetRes(conn)
	switch response_name {
	case "OK":
		break
	case "Not Acceptable":
		fmt.Println("cygnudge login: user doesn't exist")
		os.Exit(3)
	}

password_input:
	//input password and encrypt it, then send to server for verification
	fmt.Printf("input password: ")
	unverified_password, err := term.ReadPassword(int(syscall.Stdin))
	if err != nil {
		log.Fatalln(err)
	}
	//debug
	fmt.Printf("\nunverified password: %s\n", string(unverified_password))

	//encrypt the password
	hashed_unverified_password := fmt.Sprintf("%x", md5.Sum(unverified_password))
	//debug
	fmt.Printf("hashed unverified password: %s\n", hashed_unverified_password)

	//send to server to verify
	_, err = conn.Write([]byte(hashed_unverified_password))
	if err != nil {
		log.Fatalln(err)
	}
	response_code, response_name := cygnudge.ReceiveGetRes(conn)
	switch response_name {
	case "OK":
		fmt.Println("login successfully")
	case "Not Acceptable":
		fmt.Println("incorrect password, please input again")
		goto password_input
	default:
		log.Fatalf("%s %s\n", response_code, response_name)
	}

	//receive uuid and save it in token
	buf := make([]byte, 64) //36
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	if n != 36 {
		log.Fatalf("invalid token length: %d\n", n)
	}
	token := string(buf[:n])
	log.Printf("receive new login token: %s\n", token)
	cygnudge.SendRes("OK", conn)

	user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	token_path := path.Join(user.HomeDir, ".cygnudge/login_token")
	token_file, err := os.OpenFile(token_path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	token_file.WriteString(token)
	token_file.Close()

	uid_path := path.Join(user.HomeDir, ".cygnudge/login_uid")
	uid_file, err := os.OpenFile(uid_path, os.O_TRUNC|os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		log.Fatalln(err)
	}
	uid_file.WriteString(strconv.Itoa(uid))
	uid_file.Close()
}

func checkLogin(add ServerAddress, uid int, local_token string) (string, string) {
	if local_token == "" { //token file doesn't exist
		return "120", "Not Acceptable"
	}
	conn, err := net.Dial("tcp", add.Ip+":"+strconv.Itoa(add.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	cygnudge.SendReq("Check Login", conn)
	cygnudge.ReceiveRes("OK", conn)

	//send uid
	_, err = conn.Write([]byte(strconv.Itoa(uid)))
	if err != nil {
		log.Fatalln(err)
	}
	_, response_name := cygnudge.ReceiveGetRes(conn)
	switch response_name {
	case "OK":
		break
	case "Not Acceptable":
		fmt.Printf("cygnudge check login: user doesn't exist (uid=%d)\n", uid)
		os.Exit(3)
	}

	//send token (uuid)
	_, err = conn.Write([]byte(local_token))
	if err != nil {
		log.Fatalln(err)
	}
	response_code, response_name := cygnudge.ReceiveGetRes(conn)
	cygnudge.SendRes("OK", conn)
	return response_code, response_name
	//"OK":             user has login from this client
	//"Not Acceptable": user hasn't login
	//"Remote Login":   user has login from another client
}

func logout(add ServerAddress, uid int, token string) {
	/*
		1. send uid
		2. send token
		3. server: if token does not exist: return response "Not Acceptable"
		if rdb.Get(uid+"-token")==token: remove token in redis, return response "OK"
		else: return response "Remote Login"
	*/
	conn, err := net.Dial("tcp", add.Ip+":"+strconv.Itoa(add.Port))
	if err != nil {
		log.Fatalln(err)
	}
	defer conn.Close()
	cygnudge.SendReq("Logout", conn)
	cygnudge.ReceiveRes("OK", conn)

	//send uid
	_, err = conn.Write([]byte(strconv.Itoa(uid)))
	if err != nil {
		log.Fatalln(err)
	}
	_, response_name := cygnudge.ReceiveGetRes(conn)
	switch response_name {
	case "OK":
		break
	case "Not Acceptable":
		fmt.Printf("cygnudge logout: user doesn't exist (uid=%d)\n", uid)
		os.Exit(4)
	}

	//send token (uuid)
	_, err = conn.Write([]byte(token))
	if err != nil {
		log.Fatalln(err)
	}
	_, response_name = cygnudge.ReceiveGetRes(conn)
	cygnudge.SendRes("OK", conn)
	switch response_name {
	case "OK":
		break
	case "Not Acceptable":
		fmt.Println("cygnudge logout: no user login")
		os.Exit(4)
	case "Remote Login":
		fmt.Printf("cygnudge logout: user has login from remote client")
		os.Exit(4)
	}
}

func generate_task_json(language string, pid string, uid int) (string, string) {
	loc, _ := time.LoadLocation("UTC")
	time_point := time.Now().In(loc).Format("2006-01-02_15:04:05")
	JsonMap := map[string]any{
		"language": language,
		"pid":      pid,
		"uid":      uid,
		"time":     time_point, //todo
	}
	data, err := json.Marshal(JsonMap)
	if err != nil {
		log.Fatalln(err)
	}
	return string(data), time_point
}

// language is just the suffix of code file
func pack_judgement(pid string, uid int, code_file string) string {
	dot_point := strings.LastIndex(code_file, ".")
	lang := code_file[dot_point+1:]
	task_json_string, time_point := generate_task_json(lang, pid, uid)
	fmt.Println(task_json_string)

	code_file, err := filepath.Abs(code_file)
	if err != nil {
		log.Fatalln(err)
	}

	zip_name := fmt.Sprintf("%v_%v_%v.zip", time_point, uid, pid)
	zip_path := path.Join("/tmp", zip_name)
	file, err := os.Create(zip_path) //todo: complete the path and file name
	if err != nil {
		log.Fatalln(err)
	}
	defer file.Close()

	w := zip.NewWriter(file)
	defer w.Close()

	f, err := w.Create("task.json")
	if err != nil {
		log.Fatalln(err)
	}
	f.Write([]byte(task_json_string))

	f, err = w.Create("code." + lang)
	if err != nil {
		log.Fatalln(err)
	}
	code, err := os.Open(code_file)
	if err != nil {
		log.Fatalln(err)
	}
	info, err := code.Stat()
	if err != nil {
		log.Fatalln(err)
	}
	if info.IsDir() {
		fmt.Println("cygnudge judge: code file is a directory")
		os.Exit(5)
	}
	_, err = io.Copy(f, code)
	if err != nil {
		log.Fatalln(err)
	}
	return zip_path
} //warning: should handle error and send response to server, but not directly call log.Fatalln

func judge(add ServerAddress, uid int, token string, pid string, code_file string) {
	zip_path := pack_judgement(pid, uid, code_file)
	fmt.Println(zip_path) //debug
}

func removeLoginUidToken() {
	user, err := user.Current()
	if err != nil {
		log.Fatalln(err)
	}
	token_path := path.Join(user.HomeDir, ".cygnudge/login_token")
	uid_path := path.Join(user.HomeDir, ".cygnudge/login_uid")
	err = os.Remove(token_path)
	if err != nil {
		log.Fatalln(err)
	}
	err = os.Remove(uid_path)
	if err != nil {
		log.Fatalln(err)
	}
}
