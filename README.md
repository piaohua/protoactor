# protoactor

    This sample is a game server based on protoactor-go

    Forked from https://github.com/AsynkronIT/protoactor-go

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
