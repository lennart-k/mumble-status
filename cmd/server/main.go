package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/lennart-k/mumble-status/pkg/status"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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

type StatusExporter struct {
	MumbleAddress     string
	users_online      prometheus.Gauge
	users_max         prometheus.Gauge
	allowed_bandwidth prometheus.Gauge
}

func NewStatusExporter(mumble_address string) StatusExporter {
	users_online := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "users_online",
		Namespace: "mumble_server",
	})
	users_max := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "users_max",
		Namespace: "mumble_server",
	})
	allowed_bandwidth := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:      "allowed_bandwidth",
		Namespace: "mumble_server",
	})
	return StatusExporter{
		MumbleAddress:     mumble_address,
		users_online:      users_online,
		users_max:         users_max,
		allowed_bandwidth: allowed_bandwidth,
	}
}

func (s *StatusExporter) Collect(ch chan<- prometheus.Metric) {
	status, err := status.GetServerStatus(s.MumbleAddress)
	if err != nil {
		// I don't know if this is the preferred way, I'd actually just like to return an error for Collect
		ch <- prometheus.NewInvalidMetric(s.users_online.Desc(), errors.New("Could not fetch server status"))
		ch <- prometheus.NewInvalidMetric(s.users_max.Desc(), errors.New("Could not fetch server status"))
		ch <- prometheus.NewInvalidMetric(s.allowed_bandwidth.Desc(), errors.New("Could not fetch server status"))
		return
	}

	s.users_online.Set(float64(status.UsersOnline))
	s.users_max.Set(float64(status.UsersMax))
	s.allowed_bandwidth.Set(float64(status.AllowedBandwidth))
	ch <- s.users_online
	ch <- s.users_max
	ch <- s.allowed_bandwidth
}

func (s *StatusExporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- s.users_online.Desc()
	ch <- s.users_max.Desc()
	ch <- s.allowed_bandwidth.Desc()
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
	reg := prometheus.NewRegistry()
	collector := NewStatusExporter(mumble_address)
	reg.MustRegister(&collector)
	router.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{EnableOpenMetrics: true}))

	log.Fatal(http.ListenAndServe(listenAddr, &router))
}
