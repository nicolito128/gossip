package adapters

import (
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/nicolito128/gossip"
)

func UpgradeWS(w http.ResponseWriter, r *http.Request, h http.Header) (up *WebSocketTransport, err error) {
	up = NewWebSocketTransport()
	err = up.Upgrade(w, r, h)
	return
}

type WebSocketTransport struct {
	config        *gossip.TransportConfig
	writer        http.ResponseWriter
	req           *http.Request
	customHeaders http.Header
	conn          *websocket.Conn
}

func NewWebSocketTransport(opts ...gossip.TransportOpt) *WebSocketTransport {
	tc := gossip.DefaultTransportConfig(opts...)
	wst := new(WebSocketTransport)
	wst.config = tc
	return wst
}

func (wst *WebSocketTransport) Conn() *websocket.Conn {
	return wst.conn
}

func (wst *WebSocketTransport) Write(p gossip.TransportMessage) error {
	if wst.conn == nil {
		return websocket.ErrBadHandshake
	}

	messageType := websocket.TextMessage
	if p.MessageType != nil {
		messageType = *p.MessageType
	}

	data := p.RawData
	if data == nil {
		data = []byte{}
	}

	return wst.conn.WriteMessage(messageType, data)
}

func (wst *WebSocketTransport) Upgrade(w http.ResponseWriter, r *http.Request, h http.Header) error {
	conn, err := wst.config.WebSocketUpgrader.Upgrade(w, r, h)
	wst.conn = conn
	wst.writer = w
	wst.req = r
	wst.customHeaders = h
	return err
}

func (wst *WebSocketTransport) Close() error {
	if wst.conn != nil {
		return wst.conn.Close()
	}
	return nil
}
