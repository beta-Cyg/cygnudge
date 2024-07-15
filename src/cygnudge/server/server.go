package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"sync"
	"time"

	"cygnus.beta/cygnudge"
	"github.com/go-redis/redis"
	"github.com/google/uuid"
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
var user_number int

func initVals() {
	rdb := newClient()
	defer rdb.Close()
	tmp, err := rdb.Get("user-number").Result()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	user_num, err := strconv.Atoi(tmp)
	if err != nil {
		log.Fatalln(err)
	}
	user_number = user_num
}

func handleRegister(conn net.Conn) {
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	email := string(buf[:n])
	log.Printf("receive email: %s\n", email)

	rdb := newClient()
	defer rdb.Close()

	//check if email has been registered
	exist, err := rdb.HExists("email-uid", email).Result()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	if exist {
		cygnudge.SendRes("Not Acceptable", conn)
		log.Printf("email has been registered: %s\n", email)
	} else {
		cygnudge.SendRes("OK", conn)
	}
	if exist {
		return
	}

	buf = make([]byte, 32) //length of hashed password (md5) is 32
	n, err = conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	if n != 32 {
		log.Fatalf("invalid password length: %d\n", n)
	}
	hashed_password := string(buf[:n])
	log.Printf("receive hashed password: %s\n", hashed_password)
	cygnudge.SendRes("OK", conn)

	//generate new user information
	mutex.Lock()
	tmp, err := rdb.Get("user-number").Result()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	user_number_tmp, err := strconv.Atoi(tmp)
	if err != nil {
		log.Fatalln(err)
	}
	user_number = user_number_tmp
	user_number += 1
	new_uid := user_number
	err = rdb.Set("user-number", strconv.Itoa(user_number), 0).Err()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	mutex.Unlock()
	log.Printf("generate new user: %d\n", new_uid)
	err = rdb.HSet(strconv.Itoa(new_uid), "email", email).Err()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	err = rdb.HSet(strconv.Itoa(new_uid), "password", hashed_password).Err()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}

	//store map between email and uid
	err = rdb.HSet("email-uid", email, strconv.Itoa(new_uid)).Err()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}

	log.Printf("store new user information: %d\n", new_uid)

	_, err = conn.Write([]byte(strconv.Itoa(new_uid)))
	if err != nil {
		log.Fatalln(err)
	}
	cygnudge.ReceiveRes("OK", conn)
}

func handleLogin(conn net.Conn) {
	buf := make([]byte, 16)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	uid, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("new user login request: %d\n", uid)

	//check if the user exists
	rdb := newClient()
	defer rdb.Close()
	exist, err := rdb.Exists(strconv.Itoa(uid)).Result()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	if exist > 0 {
		cygnudge.SendRes("OK", conn)
	} else {
		cygnudge.SendRes("Not Acceptable", conn)
		log.Printf("user doesn't exist: %d\n", uid)
	}
	if !(exist > 0) {
		return
	}

	//get hashed correct password from redis database
	correct_password, err := rdb.HGet(strconv.Itoa(uid), "password").Result()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}

password_verification:
	//receive hashed unverified password from client
	buf = make([]byte, 32) //md5
	n, err = conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	if n != 32 {
		log.Fatalf("invalid password length: %d\n", n)
	}
	unverified_password := string(buf[:n])
	log.Printf("receive unverified password: %s\n", unverified_password)

	//compare unverified password with correct password
	if unverified_password == correct_password {
		cygnudge.SendRes("OK", conn)
	} else {
		cygnudge.SendRes("Not Acceptable", conn)
		goto password_verification
	}

	//generate token (uuid)
	token := uuid.New().String()

	//save for 24 hours in redis
	token_key := fmt.Sprintf("%d-token", uid)
	err = rdb.Set(token_key, token, time.Hour*24).Err()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}

	//send token to cilent
	_, err = conn.Write([]byte(token))
	if err != nil {
		log.Fatalln(err)
	}
	cygnudge.ReceiveRes("OK", conn)
}

func handleCheckLogin(conn net.Conn) {
	buf := make([]byte, 16)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	uid, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		log.Fatalln(err)
	}
	//receive uid

	rdb := newClient()
	defer rdb.Close()

	//check if the user exists
	exist, err := rdb.Exists(strconv.Itoa(uid)).Result()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	if exist > 0 {
		cygnudge.SendRes("OK", conn)
	} else {
		cygnudge.SendRes("Not Acceptable", conn)
		log.Printf("user doesn't exist: %d\n", uid)
	}
	if !(exist > 0) {
		return
	}

	buf = make([]byte, 64) //36
	n, err = conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	if n != 36 {
		log.Fatalf("invalid token length: %d\n", n)
	}
	token := string(buf[:n])
	//receive token

	token_key := fmt.Sprintf("%d-token", uid)
	exist, err = rdb.Exists(token_key).Result()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	if exist > 0 { //token exist: has login
		correct_token, err := rdb.Get(token_key).Result()
		if err != redis.Nil && err != nil {
			log.Fatalln(err)
		}
		if token == correct_token { //has login from this client
			cygnudge.SendRes("OK", conn)
		} else { //has login from remote client
			cygnudge.SendRes("Remote Login", conn)
		}
	} else { //token doesn't exist
		cygnudge.SendRes("Not Acceptable", conn)
	}
	cygnudge.ReceiveRes("OK", conn)
}

func handleGetUid(conn net.Conn) {
	buf := make([]byte, 128) //maximum length of email is 128
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	email := string(buf[:n])

	rdb := newClient()
	defer rdb.Close()

	//check if user exists
	exist, err := rdb.HExists("email-uid", email).Result()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	if exist {
		cygnudge.SendRes("OK", conn)
	} else { //user doesn't exist
		cygnudge.SendRes("Not Acceptable", conn)
		log.Printf("user doesn't exist: %s\n", email)
	}
	if !exist {
		return
	}

	uid, err := rdb.HGet("email-uid", email).Result()
	if err != redis.Nil && err != nil {
		//todo SendRes Unknown Error
		log.Fatalln(err)
	}
	_, err = conn.Write([]byte(uid))
	if err != nil {
		log.Fatalln(err)
	}
	cygnudge.ReceiveRes("OK", conn)
}

func handleLogout(conn net.Conn) {
	buf := make([]byte, 16)
	n, err := conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	uid, err := strconv.Atoi(string(buf[:n]))
	if err != nil {
		log.Fatalln(err)
	}

	rdb := newClient()
	defer rdb.Close()

	exist, err := rdb.Exists(strconv.Itoa(uid)).Result()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	if exist > 0 {
		cygnudge.SendRes("OK", conn)
	} else {
		cygnudge.SendRes("Not Acceptable", conn)
		log.Printf("user doesn't exist: %d\n", uid)
	}
	if !(exist > 0) {
		return
	}

	buf = make([]byte, 64) // 36
	n, err = conn.Read(buf)
	if err != nil {
		log.Fatalln(err)
	}
	if n != 36 {
		log.Fatalf("invalid token length: %d\n", n)
	}
	token := string(buf[:n])

	token_key := fmt.Sprintf("%d-token", uid)
	exist, err = rdb.Exists(token_key).Result()
	if err != redis.Nil && err != nil {
		log.Fatalln(err)
	}
	if exist > 0 {
		correct_token, err := rdb.Get(token_key).Result()
		if err != redis.Nil && err != nil {
			log.Fatalln(err)
		}
		if correct_token == token {
			err = rdb.Del(token_key).Err()
			if err != redis.Nil && err != nil {
				log.Fatalln(err)
			}
			cygnudge.SendRes("OK", conn)
		} else {
			cygnudge.SendRes("Remote Login", conn) //cannot logout
		}
	} else { //no user login
		cygnudge.SendRes("Not Acceptable", conn)
	}
	cygnudge.ReceiveRes("OK", conn)
}

func process(conn net.Conn) {
	defer conn.Close()
	defer func() {
		recover()
		log.Println("process end")
	}()
	_, request_name := cygnudge.ReceiveGetReq(conn) //if is undefined request code, then Println()
	switch request_name {
	case "Register":
		cygnudge.SendRes("OK", conn)
		handleRegister(conn)
	case "Login":
		cygnudge.SendRes("OK", conn)
		handleLogin(conn)
	case "Check Login":
		cygnudge.SendRes("OK", conn)
		handleCheckLogin(conn)
	case "Get Uid":
		cygnudge.SendRes("OK", conn)
		handleGetUid(conn)
	case "Logout":
		cygnudge.SendRes("OK", conn)
		handleLogout(conn)
	default:
		cygnudge.SendRes("Bad Request", conn)
	}
}

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
