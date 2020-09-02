# base-server-go

## Installation
```
go get -u github.com/heaven-chp/base-server-go
```

## How to add config
 - json type config file add
   - see config/socket_server.config 
 - test add
   - see config_test.go
 - struct add
   - see config/config.go
 - example of use
   - socketServerConfig of socket_server/main.go

## How to use socket server
 - install
   - go install github.com/heaven-chp/base-server-go/socket_server
 - run
   - ./bin/socket_server -config_file src/github.com/heaven-chp/base-server-go/config/socket_server.config
 - log
   - ./log/socket_server_YYYYMMDD.log 
