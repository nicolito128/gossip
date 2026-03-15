package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/nicolito128/gossip"
)

var (
	addr    = flag.String("addr", ":8080", "http address")
	manager = gossip.NewManager()
)

func main() {
	flag.Parse()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", serveWS)

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

func serveWS(w http.ResponseWriter, r *http.Request) {
	tp, err := gossip.UpgradeWS(w, r, nil)
	if err != nil {
		log.Println("serveWS:", err)
		return
	}

	ch := manager.Subscribe("ticker", tp)
	ticker := time.NewTicker(time.Second * 1)
	for tick := range ticker.C {
		tp.SetMessageType(websocket.TextMessage)
		msg := []byte(fmt.Sprintf("%s", tick.Format(time.RFC3339)))
		ch.Publish(msg)
	}
}
