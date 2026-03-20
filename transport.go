package gossip

type TransportKind int

const (
	TransportWS TransportKind = iota
	TransportSSE
)

type TransportMessage struct {
	RawData []byte

	Kind TransportKind
	WS   *WSMessageOptions
	SSE  *SSEMessageOptions
}

type WSMessageOptions struct {
	MessageType int
}

type SSEMessageOptions struct {
	Event string
	ID    string
	Retry int
}

// Transporter interface ...
type Transporter interface {
	Write(p TransportMessage) error
	Close() error
}
