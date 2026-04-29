// Package proto provides a Codec for encoding and decoding framed JSON messages.
package proto

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

var (
	// ErrFrameTooLarge means the frame's length exceeds MaxFrameSize.
	ErrFrameTooLarge = errors.New("frame exceeds max size")

	// ErrUnknownType means we received a valid envelope with a Type we don't recognize.
	ErrUnknownType = errors.New("unknown message type")
)

// Codec encodes and decodes messages to and from an underlying io.ReadWriter.
type Codec struct {
	r *bufio.Reader
	w io.Writer
}

// NewCodec creates a new Codec that reads from and writes to rw.
func NewCodec(rw io.ReadWriter) *Codec {
	return &Codec{
		r: bufio.NewReader(rw),
		w: rw,
	}
}

// WriteMessage encodes msg as JSON, prefixes it with a 4-byte big-endian length,
// and writes it to the underlying writer.
func (c *Codec) WriteMessage(msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message data: %w", err)
	}

	env := envelope{Type: msg.Type(), Data: data}
	payload, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("marshal envelope")
	}

	if len(payload) > MaxFrameSize {
		return fmt.Errorf("write %s: %w (got %d, max %d)",
			msg.Type(), ErrFrameTooLarge, len(payload), MaxFrameSize)
	}

	frame := make([]byte, 4+len(payload))
	binary.BigEndian.PutUint32(frame[:4], uint32(len(payload))) //nolint:gosec // bounded above by MaxFrameSize
	copy(frame[4:], payload)

	if _, err := c.w.Write(frame); err != nil {
		return fmt.Errorf("write frame: %w", err)
	}
	return nil
}

// ReadMessage reads one frame from the underlying reader, checks the length,
// decodes the JSON envelope, and returns the decoded message.
func (c *Codec) ReadMessage() (Message, error) {
	var header [4]byte
	if _, err := io.ReadFull(c.r, header[:]); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(header[:])
	if length > MaxFrameSize {
		return nil, fmt.Errorf("read frame: %w (got %d, max %d)",
			ErrFrameTooLarge, length, MaxFrameSize)
	}

	if length == 0 {
		return nil, fmt.Errorf("read frame: empty frame")
	}

	payload := make([]byte, length)
	if _, err := io.ReadFull(c.r, payload); err != nil {
		return nil, fmt.Errorf("read frame body: %w", err)
	}

	var env envelope
	if err := json.Unmarshal(payload, &env); err != nil {
		return nil, fmt.Errorf("unmarshal envelope: %w", err)
	}

	return decodeByType(env.Type, env.Data)
}

func decodeByType(t MessageType, data json.RawMessage) (Message, error) {
	switch t {
	case TypeEchoReq:
		var m EchoReq
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("unmarshal %s: %w", t, err)
		}
		return &m, nil
	case TypeEchoResp:
		var m EchoResp
		if err := json.Unmarshal(data, &m); err != nil {
			return nil, fmt.Errorf("unmarshal %s: %w", t, err)
		}
		return &m, nil
	default:
		return nil, fmt.Errorf("%w: %q", ErrUnknownType, t)
	}
}
