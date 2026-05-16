package proto

import "encoding/json"

// MaxFrameSize is the maximum size of payload in bytes the codec will read and write.
const MaxFrameSize = 64 * 1024

// MessageType is the type of message.
type MessageType string

// These are the supported message types. The client and server must agree on these.
const (
	TypeEchoReq       MessageType = "echo_req"
	TypeEchoResp      MessageType = "echo_resp"
	TypeTunnelRequest MessageType = "tunnel_request"
	TypeTunnelCreated MessageType = "tunnel_created"
	TypeTunnelErr     MessageType = "tunnel_err"
	TypeNewConn       MessageType = "new_conn"
	TypeTunnelClose   MessageType = "tunnel_close"
	TypeTunnelClosed  MessageType = "tunnel_closed"
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

// TunnelProto defines the protocol used in the tunnel.
type TunnelProto string

const (
	// ProtoTCP represents the TCP protocol.
	ProtoTCP TunnelProto = "tcp"
	// ProtoHTTP represents the HTTP protocol.
	ProtoHTTP TunnelProto = "http"
)

// TunnelRequest is sent by client to request a new tunnel.
type TunnelRequest struct {
	Proto     TunnelProto `json:"proto"`
	LocalPort int         `json:"local_port"`
	Subdomain string      `json:"subdomain,omitempty"`
}

// Type returns the message type for TunnelRequest.
func (*TunnelRequest) Type() MessageType { return TypeTunnelRequest }

// TunnelCreated is sent by server when the tunnel is successfully created.
type TunnelCreated struct {
	TunnelID   string `json:"tunnel_id"`
	PublicAddr string `json:"public_addr"`
}

// Type returns the message type for TunnelCreated.
func (*TunnelCreated) Type() MessageType { return TypeTunnelCreated }

// TunnelErr is sent by server when it fails.
type TunnelErr struct {
	Reason string `json:"reason"`
}

// Type returns the message type for TunnelErr.
func (*TunnelErr) Type() MessageType { return TypeTunnelErr }

// NewConn is sent by server when a new connection is established to the tunnel.
type NewConn struct {
	TunnelID string `json:"tunnel_id"`
	ConnID   string `json:"conn_id"`
}

// Type returns the message type for NewConn.
func (*NewConn) Type() MessageType { return TypeNewConn }

// TunnelClose is sent by client to close the tunnel.
type TunnelClose struct {
	TunnelID string `json:"tunnel_id"`
}

// Type returns the message type for TunnelClose.
func (*TunnelClose) Type() MessageType { return TypeTunnelClose }

// TunnelClosed is sent by server when the tunnel is closed.
type TunnelClosed struct {
	TunnelID string `json:"tunnel_id"`
	Reason   string `json:"reason"`
}

// Type returns the message type for TunnelClosed.
func (*TunnelClosed) Type() MessageType { return TypeTunnelClosed }
