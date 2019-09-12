netpro
======

Golang net programming


## Example

Running a tcp server [tcpserver.go](https://raw.githubusercontent.com/axiaoxin/netpro/master/example/tcpserver.go)

```

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
	}
}

func main() {
	server := netpro.NewTCPServer(":9090", true)
	server.Run(handler)
}
```

Running an udp server [udpserver.go](https://raw.githubusercontent.com/axiaoxin/netpro/master/example/udpserver.go)

```
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
```

## Config

Copy the `config.toml.example` to your workdir as `config.toml`
