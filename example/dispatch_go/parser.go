package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"net/url"
	"strconv"
	"time"

	"github.com/axiaoxin/netpro"
	"github.com/pkg/errors"

	"dispatch_go/dispatch_pb"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
)

// dc 转发出来的数据协议格式为 0x02|cVer|wlen|data|0x03 (wlen: c++中word类型即unsigned short对应go中uint16占用2byte)
// link 接收的数据协议格式为 0x02|len|data|0x03 (len: c++中double word类型即unsigned long对应go中的uint32占用4byte)

const (
	LinkPkgType int32 = 1
	DCPkgType   int32 = 2

	DCHeadFlag          byte = 0x02
	DCTailFlag          byte = 0x03
	DCHeadLength             = 4
	DCHeadFlagIndex          = 0
	DCHeadVerIndex           = 1
	DCHeadLenStartIndex      = 2
	DCHeadLenEndIndex        = 3

	LinkHeadFlag          byte = 0x02
	LinkTailFlag          byte = 0x03
	LinkHeadLength             = 5
	LinkHeadFlagIndex          = 0
	LinkHeadLenStartIndex      = 1
	LinkHeadLenEndIndex        = 4
)

type DataParser interface {
	ReadData() ([]byte, error)
	WriteData([]byte) (int, error)
}

// DCParser 解析DC转发的数据
type DCParser struct {
	Conn *net.UDPConn
}

type LinkParser struct {
	Conn net.Conn
}

// DispatchParser 解析通过link_dispatch链路上报的数据，上报的数据协议和DC链路的协议相同
type DispatchParser struct {
	DCParser
}

// NewDCParser 实例化DC协议UDP的Parser用于读写conn的数据
func NewDCParser(conn *net.UDPConn) *DCParser {
	p := &DCParser{}
	p.Conn = conn
	return p
}

// ReadData 读取DC发送过来的数据，返回原始的不包含头尾的原始数据
func (p *DCParser) ReadData() ([]byte, error) {
	buf := make([]byte, viper.GetInt("runtime.udp_read_size"))
	n, remote, err := p.Conn.ReadFrom(buf)
	if err != nil {
		return nil, err
	}
	// 检查头部标志
	if buf[DCHeadFlagIndex] != DCHeadFlag {
		return nil, fmt.Errorf("dc parser read invalid data head %v\n", buf[DCHeadFlagIndex])
	}
	// 检查数据长度
	if len(buf) <= DCHeadLength {
		return nil, errors.New("dc parser read incomplete data")
	}
	// 去掉头尾
	data := buf[DCHeadLength : n-1]
	netpro.Logger.Debugf("read %d bytes from %s, data:%s", n, remote, string(data))
	return data, nil
}

func (p *DCParser) ProtoMarshal(data []byte) ([]byte, error) {
	// data string: __logname=dc00911&__timestamp=1568100181&__clientip=100.97.6.100&__clientip=100.97.6.100&masterip=10.219.134.238&bid=101021758&op_type=get&val=1&avg_ms=99.986&over_50ms=1&recatmax=10.229.146.214&access_type=ccns_access
	m, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, errors.Wrap(err, "dc parser parse query error")
	}
	ts, err := strconv.Atoi(m["__timestamp"][0])
	if err != nil {
		return nil, errors.Wrap(err, "dc parser convert timestamp error")
	}
	timestamp := uint32(ts)
	logDataPkg := &dispatch_pb.LOG_DATA_PKG{
		Dataid:    &m["__logname"][0],
		Timestamp: &timestamp,
		Clientip:  &m["__clientip"][0],
		Logmsg:    data,
	}
	mRealPkg, err := proto.Marshal(logDataPkg)
	if err != nil {
		return nil, errors.Wrap(err, "dc parser marshal logdatapkg error")
	}
	pkgType := DCPkgType // 常量不能取地址
	commPkg := &dispatch_pb.COMM_PKG{
		PkgType: &pkgType,
		RealPkg: mRealPkg,
	}
	mCommPkg, err := proto.Marshal(commPkg)
	if err != nil {
		return nil, errors.Wrap(err, "dc parser marshal commpkg error")
	}
	return mCommPkg, nil
}

// NewLinkParser 实例化link服务的TCP Parser用于读写conn的数据
func NewLinkParser(conn net.Conn) *LinkParser {
	p := &LinkParser{}
	p.Conn = conn
	return p
}

// WriteData 向用net_svc框架实现的link服务发送数据，参数为不包含头尾的原始数据
func (p *LinkParser) WriteData(data []byte) (int, error) {
	dataLen := len(data)
	if dataLen == 0 {
		return 0, nil
	}
	// 发送数据的总长度 头+数据+尾
	totalLen := LinkHeadLength + dataLen + 1

	buff := bytes.NewBuffer([]byte{})
	binary.Write(buff, binary.BigEndian, LinkHeadFlag)     // 添加协议头
	binary.Write(buff, binary.BigEndian, uint32(totalLen)) // 添加完整数据长度
	binary.Write(buff, binary.BigEndian, data)             // 数据部分
	binary.Write(buff, binary.BigEndian, LinkTailFlag)     // 添加协议尾

	n, err := p.Conn.Write(buff.Bytes())
	if err != nil {
		return -1, errors.Wrap(err, "link parser write error")
	}
	netpro.Logger.Debugf("write %d bytes from %s to %s", n, p.Conn.LocalAddr(), p.Conn.RemoteAddr())
	return n, nil
}

// NewDispatchParser 实例化用于解析link_dispatch链路上报的数据的parser
func NewDispatchParser(conn *net.UDPConn) *DispatchParser {
	// 套嵌匿名结构体必须要这样写才能对字段赋值，不能直接在{}里面写
	p := &DispatchParser{}
	p.Conn = conn
	return p
}

// ProtoMarshal 转换为可用于直接进行转发到一下链路的pb协议的数据
func (p *DispatchParser) ProtoMarshal(data []byte) ([]byte, error) {
	// data string: __logname=dc00911&data=a%7Cb  %7C -> |
	m, err := url.ParseQuery(string(data))
	if err != nil {
		return nil, errors.Wrap(err, "dispatch parser parse query error")
	}
	timestamp := uint32(time.Now().Unix())
	localIP, err := netpro.GetOutboundIP()
	if err != nil {
		return nil, errors.Wrap(err, "get outbound ip error")
	}
	localIPStr := localIP.String()
	logDataPkg := &dispatch_pb.LOG_DATA_PKG{
		Dataid:    &m["__logname"][0],
		Timestamp: &timestamp,
		Clientip:  &localIPStr,
		Logmsg:    []byte(m["data"][0]),
	}
	mRealPkg, err := proto.Marshal(logDataPkg)
	if err != nil {
		return nil, errors.Wrap(err, "dispatch parser marshal logdatapkg error")
	}
	pkgType := LinkPkgType // 常量不能取地址
	commPkg := &dispatch_pb.COMM_PKG{
		PkgType: &pkgType,
		RealPkg: mRealPkg,
	}
	mCommPkg, err := proto.Marshal(commPkg)
	if err != nil {
		return nil, errors.Wrap(err, "dispatch parser marshal commpkg error")
	}
	return mCommPkg, nil
}
