package main

import (
	"context"
	"log"

	"github.com/axiaoxin/netpro"
)

func handler(ctx context.Context) error {
	conn := netpro.CtxTCPConn(ctx)

	for {
		buf := make([]byte, 1024)
		n, err := conn.Read(buf)
		if err != nil {
			log.Println(err)
			return err
		}
		log.Printf("rev %d bytes, data:%s\n", n, string(buf[:n]))
		n, err = conn.Write(buf)
		log.Printf("send %d bytes, data:%s\n", n, string(buf[:n]))
	}
}

func main() {
	server := netpro.NewTCPServer(":9090", true)
	server.Run(handler)
}
