package gossip

import (
	"net/http"

	"github.com/gorilla/websocket"
)

func UpgradeWS(w http.ResponseWriter, r *http.Request, h http.Header) (up *WebSocketTransport, err error) {
	up = NewWebSocketTransport()
	_, err = up.Upgrade(w, r, h)
	return up, err
}

type WebSocketTransport struct {
	config      *TransportConfig
	writer      http.ResponseWriter
	req         *http.Request
	conn        *websocket.Conn
	messageType int
}

func NewWebSocketTransport(opts ...TransportOpt) *WebSocketTransport {
	tc := DefaultTransportConfig(opts...)
	wst := new(WebSocketTransport)
	wst.config = tc
	return wst
}

func (wst *WebSocketTransport) Upgrade(w http.ResponseWriter, r *http.Request, h http.Header) (*websocket.Conn, error) {
	conn, err := wst.config.WebSocketUpgrader.Upgrade(w, r, nil)
	wst.conn = conn
	return wst.conn, err
}

func (wst *WebSocketTransport) SetMessageType(typ int) {
	wst.messageType = typ
}

func (wst *WebSocketTransport) Write(p []byte) (n int, err error) {
	err = wst.conn.WriteMessage(wst.messageType, p)
	n = len(p)
	return
}

func (wst *WebSocketTransport) Close() error {
	if wst.conn != nil {
		return wst.conn.Close()
	}
	return nil
}
