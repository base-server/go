# base-server-go

## How to add config
 - json type config file add
   - see [config/socket_server.config](https://github.com/heaven-chp/base-server-go/blob/main/config/socket_server.config)
 - struct add
   - see [config/config.go](https://github.com/heaven-chp/base-server-go/blob/main/config/config.go)
 - test add
   - see [config_test.go](https://github.com/heaven-chp/base-server-go/blob/main/config/config_test.go)
 - example of use
   - socketServerConfig of [socket_server/main.go](https://github.com/heaven-chp/base-server-go/blob/main/socket_server/main.go)

## How to use grpc server
 - build
   - `go build -o grpc_server ./grpc_server/`
 - run
   - `./grpc_server/grpc_server -config_file config/grpc_server.config`
 - log
   - `./log/grpc_server_YYYYMMDD.log`

## How to use socket server
 - build
   - `go build -o socket_server ./socket_server/`
 - run
   - `./socket_server/socket_server -config_file config/socket_server.config`
 - log
   - `./log/socket_server_YYYYMMDD.log`
