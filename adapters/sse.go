package adapters

import (
	"fmt"
	"net/http"
	"sync"

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
	flusher       http.Flusher
	customHeaders http.Header

	closed bool
	mu     sync.RWMutex
}

func NewSSETransport(opts ...gossip.TransportOpt) *SSETransport {
	tc := gossip.DefaultTransportConfig(opts...)
	sst := new(SSETransport)
	sst.config = tc
	return sst
}

func (sst *SSETransport) Write(p gossip.TransportMessage) error {
	sst.mu.RLock()
	if sst.closed {
		return fmt.Errorf("error: transport is closed")
	}
	if sst.writer == nil {
		return http.ErrServerClosed
	}
	sst.mu.RUnlock()

	sst.mu.Lock()
	defer sst.mu.Unlock()

	raw := p.RawData
	if raw == nil {
		raw = []byte{}
	}

	var eventFormatted string

	if p.EventID != nil {
		eventFormatted += fmt.Sprintf("id: %s\n", *p.EventID)
	}
	if p.EventName != nil {
		eventFormatted += fmt.Sprintf("event: %s\n", *p.EventName)
	}
	if p.EventRetry != nil {
		eventFormatted += fmt.Sprintf("retry: %d\n", *p.EventRetry)
	}
	eventFormatted += fmt.Sprintf("data: %s\n\n", string(raw))

	_, err := sst.writer.Write([]byte(eventFormatted))
	sst.flusher.Flush()

	return err
}

func (sst *SSETransport) Upgrade(w http.ResponseWriter, r *http.Request, h http.Header) error {
	sst.writer = w
	sst.req = r
	sst.customHeaders = h

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
	sst.flusher = flusher

	flusher.Flush()
	return nil
}

func (sst *SSETransport) Close() error {
	if sst.closed {
		return fmt.Errorf("error: transport is already closed")
	}
	// Clean up resources
	sst.writer = nil
	sst.req = nil
	sst.customHeaders = nil
	sst.closed = true
	return nil
}
