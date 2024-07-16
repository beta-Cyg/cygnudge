# Cygnudge Online Judge System

## Installation

- **todo**

- add ```export CGO_ENABLED=0``` to your bash configuration file

## Usage

### client

- `cygnudge register [--server {server}]/[--ip {ip} [--port {port}]]`

**default port: 1145**

input email address and password

send code `001` to server

receive code `100` from server

output uid of the new account

- `cygnudge login [[--uid {uid}]/[--email {email}]] [--server {server}]/[--ip {ip} --port {port}]`

if uid/email is not specified in command line, input uid/email

input password

server generates a uuid after getting the request from client, then save it in redis and send it to client.

client save uuid in `~/.cygnudge/login.token`

if the account has logined in other devices / local token was deleted (there is token in redis of server, but no token found in the speciafied path in client), then remove the token in redis and re-login? (optional)

- `cygnudge logout`

require login state (`~/.cygnudge/login.token` exists)

server removes uuid kept in redis

- `cygnudge server --name {name} --ip {ip} [--port {port}]`

server information is stored in `~/.cygnudge/server.json`

format of server.json:

```javascript
{
	"test_server" : {
		"ip" : "127.0.0.1",
		"port" : "900"
	}
}
```

- `cygnudge config`

1. `cygnudge config --password`

**require login and set server**

2. `cygnudge config --server 127.0.0.1:8088`

- `cygnudge submit --problem P1001 --code P1.cpp --language cpp`

## requirements

### runtime requirements
- `python3`

### compile time requirements
- `boost (property_tree)`
- `g++ 13+ (for C++20)`
- `go`
- `xmake`
- `make`

## default config directories

### problem directory

`/var/lib/cygnudge`

### problem judge data directory

`/var/lib/cygnudge/{pid}/data`

### problem achive directory (for tranfering to client)

`/var/lib/cygnudge/achive`

### server config path

`/etc/cygnudge/server/server.json`

### server compile.json path
`/etc/cygnudge/server/compile.json`

## compile.json format

```javascript
{
	"cpp" : "/usr/bin/g++ {0} -o {1}",
	"c" : "/usr/bin/gcc {0} -o {1}",
	"go" : "/usr/bin/go build -o {1} {0}"
}
```

## problem.zip format

### zip content

```
.
|-description.md
|-data/
  |-...
|-judge.json
```

### directory data/ content

```
.
|-0:0.in
|-0:0.out
|-0:1.in
|-0:1.out
|-1:0.in
|-1:0.out
|-...
|-m:n.in
|-m:n.out
```

### judge.json format

- time unit: micro second
- memory unit: MiB

```javascript
{
    "subtask" : 2,
    "s0" : {
        "point" : 2,
        "p0" : {
            "time" : 1000,
            "memory" : 16,
            "score" : 25
        },
        "p1" : {
            "time" : 1000,
            "memory" : 16,
            "score" : 25
        }
    },
    "s1" : {
        "point" : 1,
        "p0" : {
            "time" : 1000,
            "memory" : 16,
            "score" : 50
        }
    }
}
```

no support for optimization options yet

## task.zip format

### task.zip naming

- time format: `%Y-%M-%d_%h:%m:%s`

- task.zip format: `{time}_{uid}_{pid}.zip`

### zip content

```
.
|-code.{language}
|-task.json
```

### task.json format

```javascript
{
	"language" : "cpp",
	"pid" : "P1001",
	"uid" : 1,
	"time" : "2023-12-05_18:28:52"
}
```

## result.json format

non-CE:
```javascript
{
	"score" : 100,
	"status" : "PAC"
	"subtask" : 2,
	"s0" : {
		"point" : 2,
		"p0" : {
			"time" : 0,
			"memory" : 3,
			"return_code" : 0,
			"status" : "AC"
		},
		"p1" : {
			"time" : 0,
			"memory" : 3,
			"return_code" : 0,
			"status" : "WA"
		}
	},
	"s1" : {
		"point" : 1,
		"p0" : {
			"time" : 0,
			"memory" : 3,
			"return_code" : 0,
			"status" : "AC"
		}
	}
}
```

CE:
```javascript
{
	"score" : 0
	"status" : "CE"
}
```

## Redis format

```javascript
user-number : {number of account}
{user-id} : {
	email : {}
	password: {}
	record : {
		{problem id} : {
			content of result.json
			...
		}
	}
}
...
email-uid : {
	{email} : {uid}
	...
}
{user-id}-token : {token (uuid)} //Expiration: 24h
{problem} : {
	{problem-id}
	...
} //Set
```