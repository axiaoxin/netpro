package main

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/axiaoxin/netpro"

	"github.com/pkg/errors"
	"github.com/spf13/viper"
)

func dcHandler(ctx context.Context) error {
	// 读取dc转发的数据 (link_dispatch_l5)
	conn := netpro.CtxUDPConn(ctx)
	dcParser := NewDCParser(conn)
	data, err := dcParser.ReadData()
	if err != nil {
		return err
	}

	// 生成pb序列化后可以直接转发的数据，dispatch的pb和logacc中的pb不一样哦
	pData, err := dcParser.ProtoMarshal(data)
	if err != nil {
		return err
	}

	// 转发到logacc
	srv, err := L5Client.GetServerBySid(viper.GetInt32("l5.logacc.mod"), viper.GetInt32("l5.logacc.cmd"))
	if err != nil {
		return err
	}
	logaccAddr := fmt.Sprintf("%s:%d", srv.Ip().String(), srv.Port())
	timeout := time.Duration(viper.GetInt("runtime.dial_timeout_ms")) * time.Millisecond
	logaccConn, err := net.DialTimeout("tcp", logaccAddr, timeout)
	if err != nil {
		return errors.Wrap(err, "connect logacc server error")
	}
	defer logaccConn.Close()
	linkParser := NewLinkParser(logaccConn)
	_, err = linkParser.WriteData(pData)
	if err != nil {
		return err
	}
	return nil
}

func dispatchHandler(ctx context.Context) error {
	// 读取上报的数据 (link_dispatch)
	conn := netpro.CtxUDPConn(ctx)
	dispatchParser := NewDispatchParser(conn)
	data, err := dispatchParser.ReadData()
	if err != nil {
		return err
	}

	// 生成pb序列化后可以直接转发的数据
	pData, err := dispatchParser.ProtoMarshal(data)
	if err != nil {
		return err
	}

	// 转发到logacc
	srv, err := L5Client.GetServerBySid(viper.GetInt32("l5.logacc.mod"), viper.GetInt32("l5.logacc.cmd"))
	if err != nil {
		return errors.Wrap(err, "l5 get server by sid error")
	}
	logaccAddr := fmt.Sprintf("%s:%d", srv.Ip().String(), srv.Port())
	timeout := time.Duration(viper.GetInt("runtime.dial_timeout_ms")) * time.Millisecond
	logaccConn, err := net.DialTimeout("tcp", logaccAddr, timeout)
	if err != nil {
		return errors.Wrap(err, "connect logacc server error")
	}
	defer logaccConn.Close()
	linkParser := NewLinkParser(logaccConn)
	_, err = linkParser.WriteData(pData)
	if err != nil {
		return err
	}
	return nil
}
