package main

import (
	"encoding/json"
	"flag"
	"log"
	"os"

	registry "github.com/fmonniot/webhook-listener"
	server "github.com/fmonniot/webhook-listener"
)

var listenAddr = flag.String("listen", "localhost:8080", "<address>:<port> to listen on")
var configFile = flag.String("config", "", "Location of the config file")

func main() {
	flag.Parse()

	config, err := parseConfig()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Starting server on %s", config.ListenAddr)
	reg := registry.CreateRegistry()
	if err := server.Serve(config, reg); err != nil {
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

	if *listenAddr != "localhost:8080" || config.ListenAddr == "" {
		config.ListenAddr = *listenAddr
	}

	return config, nil
}
