package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
	"strconv"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("cygnudge: invalid arguments")
		os.Exit(1)
	}
	registerCmd := flag.NewFlagSet("register", flag.ExitOnError)
	loginCmd := flag.NewFlagSet("login", flag.ExitOnError)
	logoutCmd := flag.NewFlagSet("logout", flag.ExitOnError)
	judgeCmd := flag.NewFlagSet("judge", flag.ExitOnError)
	switch os.Args[1] {
	case "register":
		{
			var (
				registerServer string
				registerIp     string
				registerPort   int
			)
			registerCmd.StringVar(&registerServer, "server", "", "server name")
			registerCmd.StringVar(&registerIp, "ip", "", "server ip address")
			registerCmd.IntVar(&registerPort, "port", 1145, "server ip port")
			registerCmd.Parse(os.Args[2:])
			//debug
			server_address := parseServer("register", registerServer, registerIp, registerPort)
			register(server_address)
		}
	case "login":
		{
			var (
				loginServer string
				loginIp     string
				loginPort   int
				loginUid    int
				loginEmail  string
			)
			loginCmd.StringVar(&loginServer, "server", "", "server name")
			loginCmd.StringVar(&loginIp, "ip", "", "server ip address")
			loginCmd.IntVar(&loginPort, "port", 1145, "server ip port")
			loginCmd.IntVar(&loginUid, "uid", 0, "user id")
			loginCmd.StringVar(&loginEmail, "email", "", "user email")
			loginCmd.Parse(os.Args[2:])

			server_address := parseServer("login", loginServer, loginIp, loginPort)

			if loginUid != 0 && loginEmail != "" {
				fmt.Println("cygnudge login: cannot use both email and user id to login")
				os.Exit(3)
			}
			if loginUid == 0 && loginEmail == "" {
				fmt.Println("cygnudge login: no user is specified")
				os.Exit(3)
			}
			//debug
			fmt.Printf("user id: %d\nuser email: %s\n", loginUid, loginEmail)

			//check login state
			user, err := user.Current()
			if err != nil {
				log.Fatalln(err)
			}
			login_token_path := path.Join(user.HomeDir, ".cygnudge/login_token")
			_, err = os.Stat(login_token_path)
			var login_token_string string
			if os.IsNotExist(err) { //token doesn't exist
				login_token_string = ""
			} else {
				login_token_tmp, err := os.ReadFile(login_token_path)
				login_token_string = string(login_token_tmp)
				if err != nil {
					log.Fatalln(err)
				}
			}
			login_uid_path := path.Join(user.HomeDir, ".cygnudge/login_uid")
			_, err = os.Stat(login_uid_path)
			var login_uid_string string
			if os.IsNotExist(err) {
				login_uid_string = "0"
			} else {
				login_uid_tmp, err := os.ReadFile(login_uid_path)
				login_uid_string = string(login_uid_tmp)
				if err != nil {
					log.Fatalln(err)
				}
			}
			login_uid, err := strconv.Atoi(login_uid_string)
			if err != nil {
				log.Fatalln(err)
			}
			response_code, response_name := checkLogin(server_address, login_uid, login_token_string)
			switch response_name {
			case "OK":
				fmt.Printf("cygnudge login: another user has login (uid=%d), please logout first\n", login_uid)
				os.Exit(3)
			case "Not Acceptable":
				//not loged in, could login
				break
			case "Remote Login":
				fmt.Printf("cygnudge login: recorded user has login from remote client (uid=%d), please logout first\n", login_uid)
				os.Exit(3)
			default:
				log.Fatalf("unknown response code: %s\n", response_code)
			}

			if loginUid != 0 {
				response_code, response_name := checkLogin(server_address, loginUid, login_token_string)
				switch response_name {
				case "OK":
					fmt.Printf("cygnudge login: user has login (uid=%d)\n", loginUid)
					os.Exit(3)
				case "Not Acceptable":
					break //go to login
				case "Remote Login":
					fmt.Printf("cygnudge login: user has login from remote client (uid=%d)\n", loginUid)
					os.Exit(3)
				default:
					log.Fatalf("unknown response code: %s\n", response_code)
				}
			} else /*if loginEmail!=""*/ {
				response_code, response_name := checkLogin(server_address, getUid(server_address, loginEmail), login_token_string)
				switch response_name {
				case "OK":
					fmt.Printf("cygnudge login: user has login (email=%s)\n", loginEmail)
					os.Exit(3)
				case "Not Acceptable":
					break //go to login
				case "Remote Login":
					fmt.Printf("cygnudge login: user has login from remote client (email=%s)\n", loginEmail)
					os.Exit(3)
				default:
					log.Fatalf("unknown response code %s\n", response_code)
				}
			}

			if loginUid != 0 {
				login(server_address, loginUid)
			} else /*if loginEmail!=""*/ {
				login(server_address, getUid(server_address, loginEmail))
			}
		}
	case "logout":
		{
			var (
				logoutServer string
				logoutIp     string
				logoutPort   int
			)
			logoutCmd.StringVar(&logoutServer, "server", "", "server name")
			logoutCmd.StringVar(&logoutIp, "ip", "", "server ip address")
			logoutCmd.IntVar(&logoutPort, "port", 1145, "server ip port")
			logoutCmd.Parse(os.Args[2:])

			server_address := parseServer("logout", logoutServer, logoutIp, logoutPort)

			//check login state
			user, err := user.Current()
			if err != nil {
				log.Fatalln(err)
			}
			login_token_path := path.Join(user.HomeDir, ".cygnudge/login_token")
			_, err = os.Stat(login_token_path)
			var login_token_string string
			if os.IsNotExist(err) { //token doesn't exist
				fmt.Println("cygnudge logout: no user login")
				os.Exit(4)
			} else {
				login_token_tmp, err := os.ReadFile(login_token_path)
				login_token_string = string(login_token_tmp)
				if err != nil {
					log.Fatalln(err)
				}
			}
			login_uid_path := path.Join(user.HomeDir, ".cygnudge/login_uid")
			_, err = os.Stat(login_uid_path)
			var login_uid_string string
			if os.IsNotExist(err) {
				fmt.Println("cygnudge logout: no user login")
				os.Exit(4)
			} else {
				login_uid_tmp, err := os.ReadFile(login_uid_path)
				login_uid_string = string(login_uid_tmp)
				if err != nil {
					log.Fatalln(err)
				}
			}
			login_uid, err := strconv.Atoi(login_uid_string)
			if err != nil {
				log.Fatalln(err)
			}
			response_code, response_name := checkLogin(server_address, login_uid, login_token_string)
			switch response_name {
			case "OK":
				logout(server_address, login_uid, login_token_string)
				removeLoginUidToken()
			case "Not Acceptable":
				fmt.Println("cygnudge logout: invalid local token")
				removeLoginUidToken()
				os.Exit(4)
			case "Remote Login":
				fmt.Printf("cygnudge logout: user has login from remote client (uid=%d)\n", login_uid)
				removeLoginUidToken()
				os.Exit(4)
			default:
				log.Fatalf("unknown response code: %s\n", response_code)
			}
		}
	case "judge":
		{
			var (
				judgeServer string
				judgeIp     string
				judgePort   int
				judgePid    string
				judgeCode   string // Path of the code file
			)
			judgeCmd.StringVar(&judgeServer, "server", "", "server name")
			judgeCmd.StringVar(&judgeIp, "ip", "", "server ip address")
			judgeCmd.IntVar(&judgePort, "port", 1145, "server ip port")
			judgeCmd.StringVar(&judgePid, "pid", "", "problem id")
			judgeCmd.StringVar(&judgeCode, "code", "", "code file")
			judgeCmd.Parse(os.Args[2:])

			server_address := parseServer("judge", judgeServer, judgeIp, judgePort)

			//todo: check code file and problem id
			_, response_name := checkPid(server_address, judgePid)
			switch response_name {
			case "OK":
				break
			case "Not Acceptable":
				fmt.Println("cygnudge judge: problem does not exist")
				os.Exit(5)
			}

			//check login state: find local uid & token
			user, err := user.Current()
			if err != nil {
				log.Fatalln(err)
			}
			login_token_path := path.Join(user.HomeDir, ".cygnudge/login_token")
			_, err = os.Stat(login_token_path)
			var login_token_string string
			if os.IsNotExist(err) { //token doesn't exist
				fmt.Println("cygnudge judge: no user login")
				os.Exit(5)
			} else {
				login_token_tmp, err := os.ReadFile(login_token_path)
				login_token_string = string(login_token_tmp)
				if err != nil {
					log.Fatalln(err)
				}
			}
			login_uid_path := path.Join(user.HomeDir, ".cygnudge/login_uid")
			_, err = os.Stat(login_uid_path)
			var login_uid_string string
			if os.IsNotExist(err) {
				fmt.Println("cygnudge judge: no user login")
				os.Exit(5)
			} else {
				login_uid_tmp, err := os.ReadFile(login_uid_path)
				login_uid_string = string(login_uid_tmp)
				if err != nil {
					log.Fatalln(err)
				}
			}
			login_uid, err := strconv.Atoi(login_uid_string)
			if err != nil {
				log.Fatalln(err)
			}
			//check login state
			response_code, response_name := checkLogin(server_address, login_uid, login_token_string)
			switch response_name {
			case "OK":
				judge(server_address, login_uid, login_token_string, judgePid, judgeCode)
				removeLoginUidToken()
			case "Not Acceptable":
				fmt.Println("cygnudge judge: invalid local token")
				removeLoginUidToken()
				os.Exit(5)
			case "Remote Login":
				fmt.Printf("cygnudge judge: user has login from remote client (uid=%d)\n", login_uid)
				removeLoginUidToken()
				os.Exit(5)
			default:
				log.Fatalf("unknown response code: %s\n", response_code)
			}
		}
	default:
		{
			fmt.Println("cygnudge: invalid arguments")
			os.Exit(1)
		}
	}
}
