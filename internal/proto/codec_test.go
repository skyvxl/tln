package proto_test

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"reflect"
	"strings"
	"testing"

	"github.com/skyvxl/tln/internal/proto"
)

func TestCodec_RoundTrip(t *testing.T) {
	tests := []struct {
		name string
		msg  proto.Message
	}{
		{"echo_req empty", &proto.EchoReq{Text: ""}},
		{"echo_req short", &proto.EchoReq{Text: "hello"}},
		{"echo_req unicode", &proto.EchoReq{Text: "привет, мир 🌍"}},
		{"echo_resp", &proto.EchoResp{Text: "HELLO"}},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			c := proto.NewCodec(&buf)

			if err := c.WriteMessage(tc.msg); err != nil {
				t.Fatalf("WriteMessage: %v", err)
			}

			got, err := c.ReadMessage()
			if err != nil {
				t.Fatalf("ReadMessage: %v", err)
			}
			if got.Type() != tc.msg.Type() {
				t.Errorf("type = %q, want %q", got.Type(), tc.msg.Type())
			}
			if !reflect.DeepEqual(got, tc.msg) {
				t.Errorf("got %#v, want %#v", got, tc.msg)
			}
		})
	}
}

func TestCodec_MultipleFramesInOneStream(t *testing.T) {
	var buf bytes.Buffer
	c := proto.NewCodec(&buf)

	in := []proto.Message{
		&proto.EchoReq{Text: "one"},
		&proto.EchoReq{Text: "two"},
		&proto.EchoResp{Text: "THREE"},
	}
	for _, m := range in {
		if err := c.WriteMessage(m); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	for i, want := range in {
		got, err := c.ReadMessage()
		if err != nil {
			t.Fatalf("read[%d]: %v", i, err)
		}
		if !reflect.DeepEqual(got, want) {
			t.Errorf("frame %d: got %#v, want %#v", i, got, want)
		}
	}

	if _, err := c.ReadMessage(); !errors.Is(err, io.EOF) {
		t.Errorf("expected io.EOF after last frame, got %v", err)
	}
}

func TestCodec_ReadCorruptInput(t *testing.T) {
	tests := []struct {
		name    string
		input   []byte
		wantErr error
	}{
		{
			name:    "empty stream",
			input:   nil,
			wantErr: io.EOF,
		},
		{
			name:    "truncated header",
			input:   []byte{0x00, 0x00},
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "frame too large",
			input:   frameHeader(proto.MaxFrameSize + 1),
			wantErr: proto.ErrFrameTooLarge,
		},
		{
			name:    "zero length frame",
			input:   frameHeader(0),
			wantErr: nil,
		},
		{
			name:    "truncated body",
			input:   append(frameHeader(10), []byte("abc")...),
			wantErr: io.ErrUnexpectedEOF,
		},
		{
			name:    "invalid json",
			input:   frame([]byte("not json at all")),
			wantErr: nil,
		},
		{
			name:    "unknown type",
			input:   frame([]byte(`{"type":"nope","data":{}}`)),
			wantErr: proto.ErrUnknownType,
		},
		{
			name:    "valid envelope with bad data",
			input:   frame([]byte(`{"type":"echo_req","data":{"text":123}}`)),
			wantErr: nil,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var buf bytes.Buffer
			buf.Write(tc.input)
			c := proto.NewCodec(&buf)
			_, err := c.ReadMessage()
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			if tc.wantErr != nil && !errors.Is(err, tc.wantErr) {
				t.Errorf("error = %v, want errors.Is(%v) == true", err, tc.wantErr)
			}
		})
	}
}

func TestCodec_WriteTooLarge(t *testing.T) {
	var buf bytes.Buffer
	c := proto.NewCodec(&buf)

	huge := &proto.EchoReq{Text: strings.Repeat("a", proto.MaxFrameSize)}
	err := c.WriteMessage(huge)
	if !errors.Is(err, proto.ErrFrameTooLarge) {
		t.Errorf("expected ErrFrameTooLarge, got %v", err)
	}
	if buf.Len() != 0 {
		t.Errorf("nothing should be written on size error, got %d bytes", buf.Len())
	}
}

func frame(payload []byte) []byte {
	out := make([]byte, 4+len(payload))
	binary.BigEndian.PutUint32(out[:4], uint32(len(payload)))
	copy(out[4:], payload)
	return out
}

func frameHeader(n uint32) []byte {
	out := make([]byte, 4)
	binary.BigEndian.PutUint32(out, n)
	return out
}
