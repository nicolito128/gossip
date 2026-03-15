package gossip

import (
	"io"

	"github.com/gorilla/websocket"
)

var defaultUpgraderWS = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Transporter interface ...
type Transporter interface {
	io.WriteCloser
}

type TransportOpt func(tc *TransportConfig)

type TransportConfig struct {
	WebSocketUpgrader websocket.Upgrader
}

func DefaultTransportConfig(opts ...TransportOpt) *TransportConfig {
	tc := new(TransportConfig)
	tc.WebSocketUpgrader = defaultUpgraderWS
	for _, opt := range opts {
		opt(tc)
	}
	return tc
}
