package proxy

import (
	"context"
	"fmt"
	"github.com/cbuschka/go-runproxy/internal/config"
	"github.com/cbuschka/go-runproxy/internal/console"
	"github.com/cbuschka/go-runproxy/internal/util"
	"io"
	"net"
	"sync/atomic"
)

type TcpProxyStrategy struct {
	listenAddress         string
	targetEndpointAddress string
}

func NewProxy(ctx context.Context, cfg *config.Config) (*TcpProxyStrategy, error) {
	return &TcpProxyStrategy{listenAddress: cfg.ListenAddress, targetEndpointAddress: cfg.TargetAddress}, nil
}

func (h *TcpProxyStrategy) String() string {
	return fmt.Sprintf("TcpProxyStrategy{listenAddress=%s,targetEndpointAddress=%s}", h.listenAddress, h.targetEndpointAddress)
}

func (s *TcpProxyStrategy) Start(ctx context.Context, eventChan chan interface{}) {

	console.Infof("Listening on %s...", s.listenAddress)

	l, err := net.Listen("tcp", s.listenAddress)
	if err != nil {
		eventChan <- err
		return
	}

	eventChan <- "proxy started"

	go s.acceptAndProcess(ctx, l, eventChan)
}

func (s *TcpProxyStrategy) acceptAndProcess(ctx context.Context, l net.Listener, eventChan chan interface{}) {

	defer l.Close()

	connIdSeq := uint64(0)
	for {
		conn, err := l.Accept()
		if err != nil {
			eventChan <- err
			return
		}

		connId := atomic.AddUint64(&connIdSeq, 1)
		go s.handleRequest(connId, conn)
	}
}

func (s *TcpProxyStrategy) handleRequest(connId uint64, upstreamConn net.Conn) {
	console.Debugf("Conn #%d opened.", connId)

	downstreamConn, err := net.Dial("tcp", s.targetEndpointAddress)
	if err != nil {
		defer upstreamConn.Close()
		console.Infof("error connecting to downstream, %v", err.Error())
		return
	}

	defer upstreamConn.Close()
	defer downstreamConn.Close()

	nursery := util.NewNursery()
	nursery.Start(func() { s.forward(upstreamConn, downstreamConn) })
	nursery.Start(func() { s.forward(downstreamConn, upstreamConn) })
	nursery.Wait()
	console.Debugf("Conn #%d closed.", connId)
}

func (s *TcpProxyStrategy) forward(in net.Conn, out net.Conn) {
	bbuf := make([]byte, 1024)
	for {
		cnt, err := in.Read(bbuf)
		if err != nil {
			if err != io.EOF {
				console.Infof("error reading, %v", err.Error())
			}
			return
		}

		_, err = out.Write(bbuf[0:cnt])
		if err != nil {
			if err != io.EOF {
				console.Infof("error writing, %v", err.Error())
			}
			return
		}
	}
}
