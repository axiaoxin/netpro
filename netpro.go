package netpro

import (
	"github.com/spf13/viper"
)

func init() {
	initViper()
	initLogrus()

	if viper.GetBool("server.pprof") {
		go pprofServer()
	}
}
