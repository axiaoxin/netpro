package netpro

import (
	"context"
	"net"

	uuid "github.com/satori/go.uuid"
	"github.com/sirupsen/logrus"
)

const (
	// ConnIDKey use it to get conn id in ctx
	ConnIDKey = "conn_id"
	// UDPConnKey use it to get udp conn in ctx
	UDPConnKey = "udp_conn"
	// TCPConnKey use it to get tcp conn in ctx
	TCPConnKey = "tcp_conn"
)

func allocUDPConnCtx(conn *net.UDPConn) context.Context {
	ctx := context.WithValue(context.Background(), UDPConnKey, conn)

	connID := uuid.NewV4().String()
	ctx = context.WithValue(ctx, ConnIDKey, connID)

	Logger = Logger.WithFields(logrus.Fields{ConnIDKey: connID})

	return ctx
}

// CtxUDPConn get UDP conn from context
func CtxUDPConn(ctx context.Context) *net.UDPConn {
	return ctx.Value(UDPConnKey).(*net.UDPConn)
}

func allocTCPConnCtx(conn net.Conn) context.Context {
	ctx := context.WithValue(context.Background(), TCPConnKey, conn)

	connID := uuid.NewV4().String()
	ctx = context.WithValue(ctx, ConnIDKey, connID)

	Logger = Logger.WithFields(logrus.Fields{ConnIDKey: connID})

	return ctx
}

// CtxTCPConn get TCP conn from context
func CtxTCPConn(ctx context.Context) net.Conn {
	return ctx.Value(TCPConnKey).(net.Conn)
}

// CtxConnID get conn id from context
// conn id like request id
func CtxConnID(ctx context.Context) string {
	return ctx.Value(ConnIDKey).(string)
}
