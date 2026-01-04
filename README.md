# Base Server for Go

[![CI](https://github.com/base-server/go/workflows/CI/badge.svg)](https://github.com/base-server/go/actions)
[![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/heaven-chp/e51e24bb9338aae48b4465ecd2cbd620/raw/coverage.json)](https://github.com/base-server/go/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/base-server/go)](https://goreportcard.com/report/github.com/base-server/go)
[![Go Version](https://img.shields.io/github/go-mod/go-version/base-server/go?logo=go)](https://github.com/base-server/go)
[![Reference](https://pkg.go.dev/badge/github.com/base-server/go.svg)](https://pkg.go.dev/github.com/base-server/go)
[![License](https://img.shields.io/github/license/base-server/go)](https://github.com/base-server/go/blob/main/LICENSE)
[![GitHub stars](https://img.shields.io/github/stars/base-server/go)](https://github.com/base-server/go/stargazers)

<br/>

## Features
 - cloudevents
 - grpc
 - http (echo, gin, mux)
 - long polling
 - socket (tcp, udp)

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
     - `./bin/cloudevents-server -config-file ./common/config/config.yaml`
 - grpc
   - build
     - `go build -o ./bin/grpc-server ./grpc-server/`
   - run
     - `./bin/grpc-server -config_file ./common/config/config.yaml`
   - log
     - `./grpc-server.log`
 - http
   - echo
     - build
       - `go build -o ./bin/http-server-echo ./http-server/echo/`
     - run
       - `./bin/http-server-echo -config_file ./common/config/config.yaml`
     - log
       - `./http-server.log`
   - gin
     - build
       - `go build -o ./bin/http-server-gin ./http-server/gin/`
     - run
       - `./bin/http-server-gin -config_file ./common/config/config.yaml`
     - log
       - `./http-server.log`
   - mux
     - build
       - `go build -o ./bin/http-server-mux ./http-server/mux/`
     - run
       - `./bin/http-server-mux -config_file ./common/config/config.yaml`
     - log
       - `./http-server.log`
 - long-polling
   - build
     - `go build -o ./bin/long-polling-server ./long-polling-server/`
   - run
     - `./bin/long-polling-server -config_file ./common/config/config.yaml`
   - log
     - `./long-polling-server.log`
 - socket
   - tcp
     - build
       - `go build -o ./bin/socket-server-tcp ./socket-server/tcp/`
     - run
       - `./bin/socket-server-tcp -config_file ./common/config/config.yaml`
     - log
       - `./socket-server.log`
   - udp
     - build
       - `go build -o ./bin/socket-server-udp ./socket-server/udp/`
     - run
       - `./bin/socket-server-udp -config_file ./common/config/config.yaml`
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
