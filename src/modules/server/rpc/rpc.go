package rpc

import (
	"bufio"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"reflect"

	"github.com/ugorji/go/codec"
	"io"
	"net"
	"net/rpc"
	"time"
)

type Server int

func Start(rpcAddr string, logger log.Logger) error {
	// 新建rpc server
	server := rpc.NewServer()
	// 注册rpc对象
	server.Register(new(Server))

	l, err := net.Listen("tcp", rpcAddr)
	if err != nil {
		level.Error(logger).Log("msg", "fail_to_listen_addr", "rpcAddr", rpcAddr, "err", err)
		return err
	}
	level.Info(logger).Log("msg", "rpc_server_available_at", "rpcAddr", rpcAddr)

	var mh codec.MsgpackHandle
	mh.MapType = reflect.TypeOf(map[string]interface{}(nil))

	for {
		// 从accept中拿到一个客户端的连接
		conn, err := l.Accept()
		if err != nil {
			level.Warn(logger).Log("msg", "listener_accept_err", "err", err)
			time.Sleep(100 * time.Millisecond)
			continue
		}
		// 用bufferio做io解析提速

		var bufConn = struct {
			io.Closer
			*bufio.Reader
			*bufio.Writer
		}{conn, bufio.NewReader(conn), bufio.NewWriter(conn)}
		go server.ServeCodec(codec.MsgpackSpecRpc.ServerCodec(bufConn, &mh))
	}

}
