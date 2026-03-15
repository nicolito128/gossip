package gossip

import "net/http"

func UpgradeSSE(w http.ResponseWriter, r *http.Request, h http.Header) *SSETransport {
	up := new(SSETransport)
	return up
}

type SSETransport struct{}
