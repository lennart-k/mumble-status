package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lennart-k/mumble-status/pkg/status"
)

type GetStatusHandler struct {
	MumbleAddress string
}

func (h *GetStatusHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	status, err := status.GetServerStatus(h.MumbleAddress)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	out, err := json.Marshal(status)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusOK)
	w.Write(out)
}

func main() {
	host, port := "", ""
	mumble_address := ""
	flag.StringVar(&host, "host", "", "hostname, alternatively configurable via LISTEN_HOST, default(0.0.0.0)")
	flag.StringVar(&port, "port", "", "port, alternatively configurable via LISTEN_PORT, default(3000)")
	flag.StringVar(&mumble_address, "mumble-address", "", "mumble address, alternatively configurable via MUMBLE_ADDRESS")
	flag.Parse()

	if host == "" {
		host = os.Getenv("LISTEN_HOST")
	}
	if host == "" {
		host = "0.0.0.0"
	}

	if port == "" {
		port = os.Getenv("LISTEN_PORT")
	}
	if port == "" {
		port = "3000"
	}

	if mumble_address == "" {
		mumble_address = os.Getenv("MUMBLE_ADDRESS")
	}

	if host == "" || port == "" || mumble_address == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	listenAddr := fmt.Sprintf("%s:%s", host, port)

	router := http.ServeMux{}
	router.Handle("/status", &GetStatusHandler{MumbleAddress: mumble_address})

	log.Fatal(http.ListenAndServe(listenAddr, &router))
}
