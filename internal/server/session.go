package server

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"strings"

	"github.com/skyvxl/tln/internal/proto"
)

type controlSession struct {
	conn       *tls.Conn
	codec      *proto.Codec
	log        *slog.Logger
	remoteAddr string
}

func newControlSession(conn *tls.Conn, log *slog.Logger) *controlSession {
	remoteAddr := conn.RemoteAddr().String()
	return &controlSession{
		conn:       conn,
		codec:      proto.NewCodec(conn),
		log:        log.With("remote_addr", remoteAddr),
		remoteAddr: remoteAddr,
	}
}

func (s *controlSession) run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// watchdog
	go func() {
		<-ctx.Done()
		_ = s.conn.Close()
	}()

	for {
		msg, err := s.codec.ReadMessage()
		if err == nil {
			if err := s.handle(ctx, msg); err != nil {
				return fmt.Errorf("handle: %w", err)
			}
			continue
		}
		if errors.Is(err, io.EOF) || errors.Is(err, net.ErrClosed) || errors.Is(err, io.ErrUnexpectedEOF) {
			return nil
		}
		if ctx.Err() != nil {
			return nil
		}
		if errors.Is(err, proto.ErrUnknownType) {
			s.log.Warn("unknown message type", "err", err)
			continue
		}
		return fmt.Errorf("read: %w", err)
	}
}

func (s *controlSession) handle(_ context.Context, msg proto.Message) error {
	switch m := msg.(type) {
	case *proto.EchoReq:
		s.log.Debug("received echo request", "text", m.Text)
		err := s.codec.WriteMessage(&proto.EchoResp{Text: strings.ToUpper(m.Text)})
		if err != nil {
			return err
		}
		return nil
	default:
		s.log.Warn("unhandled message type", "type", m.Type())
	}
	return nil
}
