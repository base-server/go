# base-server-go

## How to add config
 - json type config file add
   - see [config/Sample.config](https://github.com/heaven-chp/base-server-go/blob/main/config/Sample.config)
 - struct add
   - see [config/Sample.go](https://github.com/heaven-chp/base-server-go/blob/main/config/Sample.go)
 - test add
   - see [Sample_test.go](https://github.com/heaven-chp/base-server-go/blob/main/config/Sample_test.go)
 - example of use
   - socketServerConfig of [socket-server/main.go](https://github.com/heaven-chp/base-server-go/blob/main/socket-server/main.go)

## How to use server
 - grpc
   - build
     - `go build -o ./bin/grpc-server ./grpc-server/`
   - run
     - `./bin/grpc-server -config_file config/GrpcServer.config`
   - log
     - `./log/grpc-server_YYYYMMDD.log`
 - socket
   - build
     - `go build -o ./bin/socket-server ./socket-server/`
   - run
     - `./bin/socket-server -config_file config/SocketServer.config`
   - log
     - `./log/socket-server_YYYYMMDD.log`
