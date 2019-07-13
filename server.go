package main

import (
	"flag"
	"log"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
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
	// Set custom port by running with --port PORT_NUM
	// Default port is 8080
	httpPort := flag.String("port", "8080", "HTTP Listening Address")
	flag.Parse()

	log.Println("Starting Server")

	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(FileServer)
	http.Handle("/", r)
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/dns-check", handleConnections)

	log.Println("Listening on port: ", *httpPort)
	err := http.ListenAndServe(":"+*httpPort, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}

	return nil
}

// Serve web files in public directory
func FileServer(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	extension, _ := regexp.MatchString("\\.+[a-zA-Z]+", r.URL.EscapedPath())
	// If the url contains an extension, use file server
	if extension {
		http.FileServer(http.Dir("./frontend/dist/")).ServeHTTP(w, r)
	} else {
		http.ServeFile(w, r, "./frontend/dist/index.html")
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request) {
	// Upgrade initial GET request to a websocket
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
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
		log.Println(msg)

		if !validateDomainName(msg.DNSname) {
			log.Println("DOMAIN WAS INVALID")
		} else {

			ws.WriteJSON(resolveDNS(msg.DNSname))
			if err != nil {
				log.Printf("error: %v", err)
				ws.Close()
				delete(clients, ws)
			}
		}
	}
}
