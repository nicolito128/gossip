package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nicolito128/gossip"
	"github.com/nicolito128/gossip/adapters"
)

var (
	addr    = flag.String("addr", ":8080", "http address")
	manager = gossip.NewManager()
)

func main() {
	flag.Parse()

	messages := make(chan string)
	defer close(messages)

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/message", handleMessage(messages))
	http.HandleFunc("/events", serveSSE(messages))

	fmt.Printf("Serving at http://localhost%s/ - Press CTRL+C to exit\n", *addr)
	if err := http.ListenAndServe(*addr, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

func handleMessage(messages chan<- string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}
		msg := r.FormValue("message")
		if msg == "" {
			http.Error(w, "Message is required", http.StatusBadRequest)
			return
		}
		messages <- msg
		w.WriteHeader(http.StatusNoContent)
	}
}

func serveSSE(messages <-chan string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tp, err := adapters.UpgradeSSE(w, r, nil)
		if err != nil {
			log.Printf("Failed to upgrade to SSE: %v", err)
			return
		}
		ch := manager.Subscribe("chat", tp)

		for msg := range messages {
			now := time.Now().Format(time.RFC3339)
			msgWithTimestamp := fmt.Sprintf("[%s] %s", now, msg)

			ch.Publish(gossip.TransportMessage{
				RawData: []byte(msgWithTimestamp),
				SSE: &gossip.SSEMessageOptions{
					Event: "new-message",
				},
			})
			log.Printf("Published message: %s", msg)
		}
	}
}
