package main

import (
	"encoding/json"
	"log"
	"net/http"

	_ "github.com/cvasq/dns-lookup-tool/statik"
	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/rakyll/statik/fs"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	"github.com/urfave/cli"
)

var clients = make(map[*websocket.Conn]bool) // connected clients

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func Start(c *cli.Context) error {

	listeningPort := c.GlobalString("listening-port")

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal(err)
	}
	staticHandler := http.FileServer(statikFS)

	router := mux.NewRouter()
	router.PathPrefix("/").Handler(staticHandler)

	http.Handle("/", router)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/dns-check", handleConnections)
	http.HandleFunc("/health", healthCheck)
	http.HandleFunc("/ready", readyCheck)

	log.Println("Server listening on port", listeningPort)
	log.Println("Web Interface: http://localhost:" + listeningPort + "/")
	log.Println("Prometheus Metrics: http://localhost:" + listeningPort + "/metrics")
	log.Println("Liveness Endpoint: http://localhost:" + listeningPort + "/health")
	log.Println("Readiness Endpoint: http://localhost:" + listeningPort + "/ready")
	err = http.ListenAndServe(":"+listeningPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return nil
}

// Healthcheck endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Up")
}

// Readiness endpoint
func readyCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode("Ready")
}

// Metrics Middleware.
var metricsCollector = middleware.New(middleware.Config{
	Recorder: metrics.NewRecorder(metrics.Config{}),
})

// Logging Middleware
func logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Requested URL: %v\n", r.URL.RequestURI())
		next.ServeHTTP(w, r)
	})
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}
	// Make sure we close the connection when the function returns
	defer ws.Close()

	// Register our new client
	clients[ws] = true

	for {
		var msg DNSresponse
		// Read in a new message as JSON and map it to a DNSname object
		err := ws.ReadJSON(&msg)
		if err != nil {
			log.Printf("error: %v", err)
			delete(clients, ws)
		}
		log.Println("Websocket data received:", msg, "from client:", ws.RemoteAddr())
		if !validateDomainName(msg.DNSname) {
			log.Println("DOMAIN INVALID:", msg.DNSname)
		} else {

			response := resolveDNS(msg.DNSname)
			ws.WriteJSON(response)
			if err != nil {
				log.Printf("error: %v", err)
				ws.Close()
				delete(clients, ws)
			}
			log.Println("Websocket data sent. DNS response to client:", ws.RemoteAddr())
		}
	}
}
