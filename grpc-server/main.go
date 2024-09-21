package main

import (
	"github.com/base-server/go/common/config"
	"github.com/base-server/go/common/main_sub"
	"github.com/common-library/go/grpc"
	"github.com/common-library/go/grpc/sample"
	"github.com/common-library/go/log/slog"
)

func main() {
	server := grpc.Server{}

	start := func(log *slog.Log) error {
		go func() {
			if err := server.Start(config.Get("gRPC.address").(string), &sample.Server{}); err != nil {
				log.Fatal("start error", "error", err)
			}
		}()

		return nil
	}

	stop := func(log *slog.Log) error {
		return server.Stop()
	}

	if err := (&main_sub.Main{}).RunWithSlog(main_sub.GRPC, start, stop); err != nil {
		panic(err)
	}
}
