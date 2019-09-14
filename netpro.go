// Package netpro means the abbreviation of golang net programming.
// netpro can be used to build a network server conveniently,
// just create the corresponding server, then specify the running port and implement the corresponding HandlerFunc processing function.
// default load configuration file in current directory which named "config", you can also start the program through the command line parameters "config_path" and "config_name" to specify the path and name of the configuration file
// the name of the configuration file does not contain suffix name, support JSON, TOML, YAML, HCL, envelope file or Java properties format.
// refer to specific configuration items through "config.toml.example"
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
