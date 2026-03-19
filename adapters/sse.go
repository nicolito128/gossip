package adapters

import (
	"fmt"
	"net/http"

	"github.com/nicolito128/gossip"
)

func UpgradeSSE(w http.ResponseWriter, r *http.Request, h http.Header) (up *SSETransport, err error) {
	up = new(SSETransport)
	err = up.Upgrade(w, r, h)
	return
}

type SSETransport struct {
	config        *gossip.TransportConfig
	writer        http.ResponseWriter
	req           *http.Request
	customHeaders http.Header

	closed bool
}

func NewSSETransport(opts ...gossip.TransportOpt) *SSETransport {
	tc := gossip.DefaultTransportConfig(opts...)
	sst := new(SSETransport)
	sst.config = tc
	return sst
}

func (sst *SSETransport) Write(p gossip.TransportMessage) error {
	if sst.closed {
		return fmt.Errorf("error: transport is closed")
	}
	if sst.writer == nil {
		return http.ErrServerClosed
	}

	data := p.RawData
	if data == nil {
		data = []byte{}
	}

	_, err := sst.writer.Write(data)
	return err
}

func (sst *SSETransport) Upgrade(w http.ResponseWriter, r *http.Request, h http.Header) error {
	if !sst.closed {
		return fmt.Errorf("error: transport is already open")
	}

	sst.writer = w
	sst.req = r
	sst.customHeaders = h

	// Set necessary headers for SSE
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
	for key, values := range h {
		for _, value := range values {
			w.Header().Add(key, value)
		}
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		return http.ErrNotSupported
	}

	flusher.Flush()
	return nil
}

func (sst *SSETransport) Close() error {
	if sst.closed {
		return fmt.Errorf("error: transport is already closed")
	}
	// Clean up resources if necessary.
	sst.writer = nil
	sst.req = nil
	sst.customHeaders = nil
	sst.closed = true
	return nil
}
