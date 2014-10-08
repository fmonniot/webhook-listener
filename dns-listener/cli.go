package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	server "github.com/fmonniot/dns-webhook-listener"
)

var listenAddr = flag.String("listen", "localhost:8080", "<address>:<port> to listen on")
var configFile = flag.String("config-file", "", "Location of handler config file")

func main() {
	flag.Parse()

	config, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on %s", config.ListenAddr)
	if err := server.Serve(config); err != nil {
		log.Fatal(err)
	}
}

func parseConfig() (*server.Config, error) {
	config := &server.Config{}

	if *configFile != "" {
		file, _ := os.Open(*configFile)
		decoder := json.NewDecoder(file)
		err := decoder.Decode(&config)
		if err != nil {
			return nil, err
		}
	}
	log.Printf("%+v", config)

	if *listenAddr != "localhost:8080" || config.ListenAddr == "" {
		config.ListenAddr = *listenAddr
	}

	return config, nil
}
