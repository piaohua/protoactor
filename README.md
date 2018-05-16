# protoactor

* This sample is a game server based on protoactor-go
* based on [protoactor-go](https://github.com/AsynkronIT/protoactor-go)

## Installation

```
GOPATH=$(cd ../"$(dirname "$0")"; pwd)
go build -o main -ldflags "-w -s" ../src/protoactor/main.go
go build -o robots -ldflags "-w -s" ../src/protoactor/robots.go
```

## Usage:

```
./main -log_dir="log" > /dev/null 2>&1 &
./robots -log_dir="log" > /dev/null 2>&1 &
```

## Note:

    As protoactor is not actively maintained please use one of the following instead
    [goplays](https://github.com/piaohua/goplays)
    [gohappy](https://github.com/piaohua/gohappy)
