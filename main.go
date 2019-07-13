package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var clients = make(map[*websocket.Conn]bool) // connected clients

// Configure the upgrader
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func main() {
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

type DNSresponse struct {
	DNSname string   `json:"dnsname"`
	A       []string `json:"A"`
	AAAA    []string `json:"AAAA"`
	MX      []string `json:"MX"`
	NS      []string `json:"NS"`
}

func resolveDNS(dnsname string) DNSresponse {

	log.Println("Resolving DNS")
	ips, err := net.LookupIP(dnsname)
	if err != nil {
		log.Println("Could not get IPs: ", err)
	}

	mx, err := net.LookupMX(dnsname)
	if err != nil {
		log.Println(err)
	}

	nss, err := net.LookupNS(dnsname)
	if err != nil {
		log.Println("NS Lookup error:", err)
	}
	if len(nss) == 0 {
		log.Println("no record")
	}
	for _, ns := range nss {
		log.Printf("%s\n", ns.Host)
	}

	ss := DNSresponse{DNSname: dnsname}

	for _, ip := range ips {
		if ip.To4() != nil {
			ss.A = append(ss.A, ip.String())
		} else {
			ss.AAAA = append(ss.AAAA, ip.String())
		}
	}

	for _, ip := range mx {
		ss.MX = append(ss.MX, fmt.Sprintf("%v\t%v", ip.Pref, ip.Host))
	}

	for _, ip := range nss {
		ss.NS = append(ss.NS, ip.Host)
	}

	return ss
}
