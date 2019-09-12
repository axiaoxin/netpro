package netpro

import (
	"net/http"
	_ "net/http/pprof"

	"github.com/spf13/viper"
)

func pprofServer() {
	addr := viper.GetString("pprof.addr")
	Logger.Info("pprof server is running on", addr)
	http.ListenAndServe(addr, nil)
}
