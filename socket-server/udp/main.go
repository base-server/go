package main

import (
	"fmt"
	"net"

	"github.com/base-server/go/common/config"
	"github.com/base-server/go/common/main_sub"
	"github.com/common-library/go/log/slog"
	"github.com/common-library/go/socket/udp"
)

func main() {
	server := udp.Server{}
	start := func(log *slog.Log) error {
		// UDP packet handler
		packetHandler := func(data []byte, addr net.Addr, conn net.PacketConn) {
			log.Debug("received packet", "from", addr.String(), "data", string(data))

			// Echo response with prefix
			response := fmt.Sprintf("[response] %s", string(data))
			if _, err := conn.WriteTo([]byte(response), addr); err != nil {
				log.Error("failed to send response", "error", err.Error(), "to", addr.String())
			} else {
				log.Debug("sent response", "to", addr.String(), "data", response)
			}
		}

		// Error handler for packet reception errors
		errorHandler := func(err error) {
			log.Error("packet reception error", "error", err.Error())
		}

		// Start UDP server
		// asyncHandler=true: process each packet concurrently
		return server.Start("udp", config.Get("socket.udp.address").(string), 1024, packetHandler, true, errorHandler)
	}

	stop := func(log *slog.Log) error {
		return server.Stop()
	}

	if err := (&main_sub.Main{}).RunWithSlog(main_sub.Socket, start, stop); err != nil {
		panic(err)
	}
}
