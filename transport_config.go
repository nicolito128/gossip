package gossip

import (
	"github.com/gorilla/websocket"
)

var defaultUpgraderWS = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// TransportOpt defines a function type for configuring the TransportConfig.
type TransportOpt func(tc *TransportConfig)

// TransportConfig holds the configuration for transports, such as WebSocket and SSE.
type TransportConfig struct {
	WebSocketUpgrader websocket.Upgrader
}

// DefaultTransportConfig returns a TransportConfig with default values, which can be overridden by providing TransportOpt functions.
func DefaultTransportConfig(opts ...TransportOpt) *TransportConfig {
	tc := new(TransportConfig)
	tc.WebSocketUpgrader = defaultUpgraderWS

	for _, opt := range opts {
		opt(tc)
	}

	return tc
}

// WithWebSocketUpgrader allows users to set a custom WebSocket upgrader in the TransportConfig.
func WithWebSocketUpgrader(upgrader websocket.Upgrader) TransportOpt {
	return func(tc *TransportConfig) {
		tc.WebSocketUpgrader = upgrader
	}
}
