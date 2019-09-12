package netpro

import (
	"context"
	"net"

	"github.com/spf13/viper"
)

// HandlerFunc developer should implement this function to handle connection
// Use CtxUDPConn or CtxTCPConn to get the connection
type HandlerFunc func(ctx context.Context) error

// TCPServer is tcp server struct
type TCPServer struct {
	// Addr is server run address
	Addr string
	// AutoCloseConn will close the connection when HandlerFunc called if it is true
	// if it is false, connection will not close when HandlerFunc called
	AutoCloseConn bool
}

// UDPServer is udp server struct
type UDPServer struct {
	// Addr is server run address
	Addr string
}

// NewUDPServer will create an UDP server
func NewUDPServer(addr string) *UDPServer {
	srv := &UDPServer{
		Addr: addr,
	}
	return srv
}

// NewTCPServer will create a TCP server
func NewTCPServer(addr string, autoCloseConn bool) *TCPServer {
	srv := &TCPServer{
		Addr:          addr,
		AutoCloseConn: autoCloseConn,
	}
	return srv
}

// Run will start the TCP server, developer should implement the handler
func (srv *TCPServer) Run(handler HandlerFunc) {
	lis, err := net.Listen("tcp", srv.Addr)
	if err != nil {
		Logger.Fatal("TCP Server listen error:", err)
		return
	}

	goroutineNum := viper.GetInt("runtime.goroutine_num")
	doneChan := make(chan struct{}, goroutineNum)
	Logger.Infof("TCP Server is listening on %s, runtime.goroutine_num=%d, handler=%s", lis.Addr(), goroutineNum, GetFunctionName(handler))
	for {
		conn, err := lis.Accept()
		if err != nil {
			Logger.Fatal("accept error:", err)
			break
		}
		doneChan <- struct{}{}
		go srv.handleConn(conn, handler, doneChan)
	}
}

// Run will start the UDP server, developer should implement the handler
func (srv *UDPServer) Run(handler HandlerFunc) {
	udpAddr, err := net.ResolveUDPAddr("udp", srv.Addr)
	if err != nil {
		Logger.Fatal("UDP Server resolve UDP addr error:", err)
	}
	conn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		Logger.Fatal("UDP Server listen error:", err)
	}
	goroutineNum := viper.GetInt("runtime.goroutine_num")
	doneChan := make(chan struct{}, goroutineNum)
	Logger.Infof("UDP Server is listening on %s, runtime.goroutine_num=%d, handler=%s", conn.LocalAddr(), goroutineNum, GetFunctionName(handler))
	for {
		doneChan <- struct{}{}
		go srv.handleConn(conn, handler, doneChan)
	}
}

func (srv *TCPServer) handleConn(conn net.Conn, handler HandlerFunc, doneChan chan struct{}) {
	if srv.AutoCloseConn {
		defer conn.Close()
	}
	ctx := allocTCPConnCtx(conn)
	if err := handler(ctx); err != nil {
		Logger.Error(err)
	}
	<-doneChan
}

func (srv *UDPServer) handleConn(conn *net.UDPConn, handler HandlerFunc, doneChan chan struct{}) {
	ctx := allocUDPConnCtx(conn)
	if err := handler(ctx); err != nil {
		Logger.Error(err)
	}
	<-doneChan
}
