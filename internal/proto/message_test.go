package proto_test

import (
	"testing"

	"github.com/skyvxl/tln/internal/proto"
)

func TestMessageTypes(t *testing.T) {
	tests := []struct {
		msg  proto.Message
		want proto.MessageType
	}{
		{msg: &proto.EchoReq{}, want: proto.TypeEchoReq},
		{msg: &proto.EchoResp{}, want: proto.TypeEchoResp},
		{msg: &proto.TunnelRequest{}, want: proto.TypeTunnelRequest},
		{msg: &proto.TunnelCreated{}, want: proto.TypeTunnelCreated},
		{msg: &proto.TunnelErr{}, want: proto.TypeTunnelErr},
		{msg: &proto.NewConn{}, want: proto.TypeNewConn},
		{msg: &proto.TunnelClose{}, want: proto.TypeTunnelClose},
		{msg: &proto.TunnelClosed{}, want: proto.TypeTunnelClosed},
	}
	for _, tc := range tests {
		t.Run(string(tc.want), func(t *testing.T) {
			if got := tc.msg.Type(); got != tc.want {
				t.Errorf("Type() = %q, want %q", got, tc.want)
			}
		})
	}
}
