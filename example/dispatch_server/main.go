package main

import (
	"fmt"
	"strings"

	"github.com/axiaoxin/netpro"
	"github.com/spf13/viper"
)

func init() {
	InitL5()
}

func runDispatch(addr string) {
	server := netpro.NewUDPServer(addr)
	handlerName := strings.ToLower(viper.GetString("runtime.handler"))
	if handlerName == "dispatch" {
		server.Run(dispatchHandler)
	} else if handlerName == "dc" {
		server.Run(dcHandler)
	}
	netpro.Logger.Error("no handler")
}

func main() {
	for i := 0; i < viper.GetInt("server.num"); i++ {
		port := viper.GetInt("server.port") + i
		addr := fmt.Sprintf(":%d", port)
		go runDispatch(addr)
	}
	select {}
}
