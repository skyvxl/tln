package server

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"sync"
)

func (s *Server) runControl(ctx context.Context) error {
	rawListener, err := net.Listen("tcp", s.cfg.ControlAddr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", s.cfg.ControlAddr, err)
	}
	listener := tls.NewListener(rawListener, s.tlsConfig)

	s.log.Info("control listener started", "addr", listener.Addr().String())

	go func() {
		<-ctx.Done()
		s.log.Info("control listener stopping")
		err := listener.Close()
		if err != nil {
			s.log.Debug("failed to close listener", "err", err)
		}
	}()

	var wg sync.WaitGroup
	defer wg.Wait()

	for {
		conn, err := listener.Accept()
		if err != nil {
			if ctx.Err() != nil {
				return nil
			}
			return fmt.Errorf("accept: %w", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			s.handleConnection(ctx, conn)
		}()
	}
}

func (s *Server) handleConnection(ctx context.Context, conn net.Conn) {
	remoteAddr := conn.RemoteAddr().String()
	defer func() {
		if err := conn.Close(); err != nil {
			s.log.Debug("failed to close connection", "remote_addr", remoteAddr, "err", err)
		}
	}()
	s.log.Info("connection accepted", "remote_addr", remoteAddr)

	tlsConn, ok := conn.(*tls.Conn)
	if ok {
		if err := tlsConn.HandshakeContext(ctx); err != nil {
			s.log.Warn("tls handshake failed", "remote_addr", remoteAddr, "err", err)
			return
		}
	}

	session := newControlSession(tlsConn, s.log)
	if err := session.run(ctx); err != nil {
		s.log.Warn("control session failed", "remote_addr", remoteAddr, "err", err)
	}

	s.log.Info("connection closed", "remote_addr", remoteAddr)
}
