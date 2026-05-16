package client

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"os"

	"github.com/skyvxl/tln/internal/proto"
)

func runEcho(_ context.Context, conn io.ReadWriter, log *slog.Logger) error {
	codec := proto.NewCodec(conn)
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return fmt.Errorf("scan stdin: %w", err)
			}
			return nil
		}
		err := codec.WriteMessage(&proto.EchoReq{Text: scanner.Text()})
		if err != nil {
			return fmt.Errorf("write message: %w", err)
		}
		msg, err := codec.ReadMessage()
		if errors.Is(err, io.EOF) ||
			errors.Is(err, net.ErrClosed) ||
			errors.Is(err, io.ErrUnexpectedEOF) {
			return nil
		}
		if err != nil {
			return fmt.Errorf("read message: %w", err)
		}
		switch m := msg.(type) {
		case *proto.EchoResp:
			fmt.Println(m.Text)
		default:
			log.Warn("unexpected message", "type", m.Type())
		}
	}
}
