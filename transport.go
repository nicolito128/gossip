package gossip

// TransportMessage represents a message that can be sent through a transport.
type TransportMessage struct {
	// Raw data to be published to the channel
	RawData []byte

	// Websocket message type (e.g., websocket.TextMessage, websocket.BinaryMessage)
	MessageType *int

	// Event name for SSE or other transport types that support event-based messaging
	EventName *string

	// ID for the event, useful for SSE to allow clients to track the last event received
	EventID *string

	// Retry time in milliseconds for SSE clients to attempt reconnection
	EventRetry *int
}

// Transporter interface ...
type Transporter interface {
	Write(p TransportMessage) error
	Close() error
}
