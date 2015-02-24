// Package webhookListener server code
package webhookListener

import (
	"log"
	"net/http"
)

// Config TODO
type Config struct {
	ListenAddr string
	TLS        struct {
		Key  string
		Cert string
	}
	Endpoints map[string]Endpoint `json:"endpoints"`
}

var serverConfig *Config
var registry Registry

// Serve TODO
func Serve(config *Config) error {
	serverConfig = config

	registry = CreateRegistry()

	for key, endpoint := range serverConfig.Endpoints {
		registry.Add(key, endpoint)
	}

	http.Handle("/", &registry)

	// Manage TLS
	if config.TLS.Key != "" && config.TLS.Cert != "" {
		log.Print("Starting with SSL")
		return http.ListenAndServeTLS(config.ListenAddr, config.TLS.Cert, config.TLS.Key, &registry)
	}

	// No TLS configured
	log.Print("Warning: Server is starting without SSL, you should not pass any credentials using this configuration")
	log.Print("To use SSL, you must provide a config file with a `tls` section, and provide locations to a `key` file and a `cert` file")

	return http.ListenAndServe(config.ListenAddr, &registry)
}
