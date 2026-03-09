package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/eryalito/smart-dns-proxy/internal/data"
	"github.com/eryalito/smart-dns-proxy/internal/dns"
)

var version = "dev"

func main() {
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(version)
		return
	}

	host := flag.String("host", "127.0.0.1", "The host to listen on")
	port := flag.String("port", "53", "The port to listen on")
	resolver := flag.String("resolver", "8.8.8.8:53", "The upstream DNS resolver to use")
	flag.Parse()

	addr := fmt.Sprintf("%s:%s", *host, *port)

	querier := &data.Querier{
		URL: "https://hayahora.futbol/estado/data.json",
	}
	// Start the background data fetcher
	go fetchDataLoop(querier)

	log.Printf("Starting DNS server on %s", addr)
	server := &dns.Server{Addr: addr, Net: "udp", Querier: querier, Resolver: *resolver}
	log.Fatal(server.Start())
}

func fetchDataLoop(querier *data.Querier) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	// Fetch data immediately at startup, then on every tick
	querier.Tick()
	for range ticker.C {
		querier.Tick()
	}
}
