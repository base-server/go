# Base Server for Go

## Features
 - grpc
 - http
 - long-polling
 - socket

<br/>

## How to add config
 - json type config file add
   - see [config/Sample.config](https://github.com/base-server/go/blob/main/config/Sample.config)
 - struct add
   - see [config/Sample.go](https://github.com/base-server/go/blob/main/config/Sample.go)
 - test add
   - see [Sample_test.go](https://github.com/base-server/go/blob/main/config/Sample_test.go)
 - example of use
   - socketServerConfig of [socket-server/main.go](https://github.com/base-server/go/blob/main/socket-server/main.go)

<br/>

## How to use server
 - cloudevents
   - build
     - `go build -o ./bin/cloudevents-server ./cloudevents-server/`
   - run
     - `./bin/cloudevents-server -config-file ./config/CloudEventsServer.config`
 - grpc
   - build
     - `go build -o ./bin/grpc-server ./grpc-server/`
   - run
     - `./bin/grpc-server -config_file ./config/GrpcServer.config`
   - log
     - `./grpc-server.log`
 - http
   - build
     - `go install github.com/swaggo/swag/cmd/swag@v1.16.3`
     - `$(go env GOPATH)/bin/swag init --dir ./http-server --output ./http-server/swagger_docs`
     - `go build -o ./bin/http-server ./http-server/`
   - run
     - `./bin/http-server -config_file ./config/HttpServer.config`
   - log
     - `./http-server.log`
 - long-polling
   - build
     - `go build -o ./bin/long-polling-server ./long-polling-server/`
   - run
     - `./bin/long-polling-server -config_file ./config/LongPollingServer.config`
   - log
     - `./long-polling-server.log`
 - socket
   - build
     - `go build -o ./bin/socket-server ./socket-server/`
   - run
     - `./bin/socket-server -config_file ./config/SocketServer.config`
   - log
     - `./socket-server.log`

<br/>

## Test and Coverage
 - Test
   - `go clean -testcache && go test -cover ./...`
 - Coverage
   - make coverage file
     - `go clean -testcache && go test -coverprofile=coverage.out -cover $(go list ./... | grep -v "/swagger_docs")`
   - convert coverage file to html file
     - `go tool cover -html=./coverage.out -o ./coverage.html`
