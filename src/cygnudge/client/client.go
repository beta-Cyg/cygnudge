package main

import (
	"encoding/json"
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
			if registerServer != "" && registerIp != "" { //todo: server and port are both specified
				fmt.Println("cygnudge register: cannot use both server name and ip address to specify the server")
				os.Exit(2)
			}
			if registerServer == "" && registerIp == "" {
				fmt.Println("cygnudge register: no server is specified")
				os.Exit(2)
			}
			fmt.Printf("server name: %s\nserver ip: %s\nserver port: %d\n", registerServer, registerIp, registerPort)
			var server_address ServerAddress
			if registerIp == "" { //use server name
				user, err := user.Current()
				if err != nil {
					log.Fatalln(err)
				}
				server_json_path := path.Join(user.HomeDir, ".cygnudge/server.json")
				server_json_string, err := os.ReadFile(server_json_path)
				if err != nil {
					log.Fatalln(err)
				}
				var servers map[string]ServerAddress
				err = json.Unmarshal(server_json_string, &servers)
				if err != nil {
					log.Fatalln(err)
				}
				tmp, exist := servers[registerServer]
				if !exist {
					fmt.Printf("cygnudge register: server %s not found\n", registerServer)
					os.Exit(2)
				} //not found
				server_address = tmp
			} else {
				server_address.Ip = registerIp
				server_address.Port = registerPort
			}
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
			if loginServer != "" && loginIp != "" {
				fmt.Println("cygnudge login: cannot use both server name and ip address to specify the server")
				os.Exit(3)
			}
			if loginServer == "" && loginIp == "" {
				fmt.Println("cygnudge login: no server is specified")
				os.Exit(3)
			}
			//debug
			fmt.Printf("server name: %s\nserver ip: %s\nserver port: %d\n", loginServer, loginIp, loginPort)
			var server_address ServerAddress
			user, err := user.Current()
			if err != nil {
				log.Fatalln(err)
			}
			if loginIp == "" { //use server name
				server_json_path := path.Join(user.HomeDir, ".cygnudge/server.json")
				server_json_string, err := os.ReadFile(server_json_path)
				if err != nil {
					log.Fatalln(err)
				}
				var servers map[string]ServerAddress
				err = json.Unmarshal(server_json_string, &servers)
				if err != nil {
					log.Fatalln(err)
				}
				tmp, exist := servers[loginServer]
				if !exist {
					fmt.Printf("cygnudge login: server %s not found\n", loginServer)
					os.Exit(3)
				} //not found
				server_address = tmp
			} else {
				server_address.Ip = loginIp
				server_address.Port = loginPort
			}

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

			if logoutServer != "" && logoutIp != "" {
				fmt.Println("cygnudge logout: cannot use both server name and ip address to specify the server")
				os.Exit(4)
			}
			if logoutServer == "" && logoutIp == "" {
				fmt.Println("cygnudge logout: no server is specified")
				os.Exit(4)
			}
			//debug
			fmt.Printf("server name: %s\nserver ip: %s\nserver port: %d\n", logoutServer, logoutIp, logoutPort)
			var server_address ServerAddress
			user, err := user.Current()
			if err != nil {
				log.Fatalln(err)
			}
			if logoutIp == "" { //use server name
				server_json_path := path.Join(user.HomeDir, ".cygnudge/server.json")
				server_json_string, err := os.ReadFile(server_json_path)
				if err != nil {
					log.Fatalln(err)
				}
				var servers map[string]ServerAddress
				err = json.Unmarshal(server_json_string, &servers)
				if err != nil {
					log.Fatalln(err)
				}
				tmp, exist := servers[logoutServer]
				if !exist {
					fmt.Printf("cygnudge logout: server %s not found\n", logoutServer)
					os.Exit(4)
				} //not found
				server_address = tmp
			} else {
				server_address.Ip = logoutIp
				server_address.Port = logoutPort
			}

			//check login state
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
				removeLoginUidToken() //todo: remove login_token and login_uid
				os.Exit(4)
			case "Remote Login":
				fmt.Printf("cygnudge login: user has login from remote client (uid=%d)\n", login_uid)
				removeLoginUidToken()
				os.Exit(4)
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
