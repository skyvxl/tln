package proto

import "encoding/json"

// MaxFrameSize is the maximum size of payload in bytes the codec will read and write.
const MaxFrameSize = 64 * 1024

// MessageType is the type of message.
type MessageType string

// These are the supported message types. The client and server must agree on these.
const (
	TypeEchoReq  MessageType = "echo_req"
	TypeEchoResp MessageType = "echo_resp"
)

// Message shows the type of message in envelope.
type Message interface {
	Type() MessageType
}

// EchoReq is sent by client. The server replies with EchoResp.
type EchoReq struct {
	Text string `json:"text"`
}

// Type returns the message type for EchoReq.
func (*EchoReq) Type() MessageType { return TypeEchoReq }

// EchoResp is the server's reply to EchoReq.
type EchoResp struct {
	Text string `json:"text"`
}

// Type returns the message type for EchoResp.
func (*EchoResp) Type() MessageType { return TypeEchoResp }

type envelope struct {
	Type MessageType     `json:"type"`
	Data json.RawMessage `json:"data"`
}
