package main

import (
	"context"
	"log"

	"github.com/axiaoxin/netpro"
)

func handler(ctx context.Context) error {
	conn := netpro.CtxUDPConn(ctx)
	buf := make([]byte, 128)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println("rev data:", string(buf[:n]))
	return nil
}

func main() {
	server := netpro.NewUDPServer(":9090")
	server.Run(handler)
}
