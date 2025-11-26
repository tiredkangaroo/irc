package main

import (
	"log/slog"
	"net"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

var debug = os.Getenv("DEBUG") == "true"

func main() {
	listenerAddr := ":0"
	if debug {
		listenerAddr = ":9000" // this is so that vite dev server can reach it easily
		// for the prod build, we use a random port (bc it runs on user's machine and we don't know what ports are free)
		// and the frontend will just use /proxy instead of localhost:port/proxy
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}

	listener, err := net.Listen("tcp", listenerAddr)
	if err != nil {
		panic(err)
	}
	slog.Info("listening", "addr", listener.Addr().String())

	// /proxy expects a request to upgrade to ws, and will relay data between
	// the ws and the host. the hostport should be specified in the irc_host query parameter.
	http.HandleFunc("/proxy", func(w http.ResponseWriter, r *http.Request) {
		irc_host := r.URL.Query().Get("irc_host")
		if irc_host == "" {
			http.Error(w, "no irc host specified", http.StatusBadRequest)
			return
		}

		upgrader := websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true // warn: this might be bad twin
			},
		}
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "failed upgrade to websocket", http.StatusInternalServerError)
			slog.Error("upgrade to websocket", "error", err)
			return
		}
		defer conn.Close()

		if err := relay(conn, irc_host); err != nil {
			slog.Error("relay", "remote_addr", r.RemoteAddr, "host", irc_host, "error", err) // may be beyond cooked to try to even return an error back to the conn, we're giving up
			return
		}
	})

	if !debug {
		fs := os.DirFS("dist")
		srv := http.FileServer(http.FS(fs))
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			srv.ServeHTTP(w, r)
		})
	} else {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			// during development, redirect to vite dev server, just convienient
			http.Redirect(w, r, "http://localhost:5173", http.StatusFound)
		})
	}

	http.Serve(listener, nil)
}
