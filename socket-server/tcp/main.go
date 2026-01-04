package main

import (
	"fmt"

	"github.com/base-server/go/common/config"
	"github.com/base-server/go/common/main_sub"
	"github.com/common-library/go/log/slog"
	"github.com/common-library/go/socket/tcp"
)

func main() {
	server := tcp.Server{}

	start := func(log *slog.Log) error {
		acceptSuccessFunc := func(client tcp.Client) {
			log.Debug("start", "network", client.GetRemoteAddr().Network(), "address", client.GetRemoteAddr().String())
			log.Debug("end", "network", client.GetRemoteAddr().Network(), "address", client.GetRemoteAddr().String())

			read := func(readJob func(readData string) error) error {
				if readData, err := client.Read(1024); err != nil {
					return err
				} else {
					log.Debug("read", "data", readData)

					return readJob(readData)
				}
			}

			write := func(writeData string) error {
				if writeLen, err := client.Write(writeData); err != nil {
					return err
				} else if writeLen != len(writeData) {
					return fmt.Errorf("invalid write - (%d)(%d)", writeLen, len(writeData))
				} else {
					log.Debug("write", "data", writeData)

					return nil
				}
			}

			if err := write("greeting"); err != nil {
				log.Error(err.Error())
				return
			}

			readJob := func(readData string) error {
				return write("[response] " + readData)
			}
			if err := read(readJob); err != nil {
				log.Error(err.Error())
				return
			}
		}

		acceptFailureFunc := func(err error) {
			log.Error(err.Error())
		}

		return server.Start("tcp", config.Get("socket.tcp.address").(string), config.Get("socket.clientPoolSize").(int), acceptSuccessFunc, acceptFailureFunc)
	}

	stop := func(log *slog.Log) error {
		return server.Stop()
	}

	if err := (&main_sub.Main{}).RunWithSlog(main_sub.Socket, start, stop); err != nil {
		panic(err)
	}
}
