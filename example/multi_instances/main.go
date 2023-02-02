package main

import (
	"flag"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	socketio "github.com/gonetlib/go-socket.io"
	"github.com/gonetlib/go-socket.io/engineio"
	"github.com/gonetlib/go-socket.io/engineio/transport"
	"github.com/gonetlib/go-socket.io/engineio/transport/polling"
	"github.com/gonetlib/go-socket.io/engineio/transport/websocket"
)

// Easier to get running with CORS. Thanks for help @Vindexus and @erkie
var allowOriginFunc = func(r *http.Request) bool {
	return true
}

var (
	listenAt = "0.0.0.0:6001"
)

func main() {
	flag.StringVar(&listenAt, "listen", listenAt, "listen at")
	flag.Parse()

	server := socketio.NewServer(&engineio.Options{
		Transports: []transport.Transport{
			&polling.Transport{
				CheckOrigin: allowOriginFunc,
			},
			&websocket.Transport{
				CheckOrigin: allowOriginFunc,
			},
		},
	})

	server.OnConnect("/", func(s socketio.Conn) error {
		s.SetContext("")
		log.Println("connected:", s.ID())
		return nil
	})

	server.OnEvent("/", "notice", func(s socketio.Conn, msg string) {
		log.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/chat", "msg", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		return "recv " + msg
	})

	server.OnEvent("/", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		log.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		log.Println("closed", reason)
	})

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()

	http.Handle("/socket.io/", server)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fbytes, err := ioutil.ReadFile("./client.html")
		if err != nil {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(err.Error()))
			return
		}
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		io.WriteString(w, string(fbytes))
	})

	log.Printf("Serving at %s...\n", listenAt)
	log.Fatal(http.ListenAndServe(listenAt, nil))
}
