// Package webhookListener server code
package webhookListener

import (
	"log"
	"net/http"
)

// Config of this server
type Config struct {
	ListenAddr string
	TLS        struct {
		Key  string
		Cert string
	}
	Endpoints []Endpoint `json:"endpoints"`
}

// Serve will register endpoints and create the http server
func Serve(config *Config, registry Registry) error {
	for _, endpoint := range config.Endpoints {
		registry.Add(endpoint)
	}

	http.Handle("/", registry)

	// Manage TLS
	if config.TLS.Key != "" && config.TLS.Cert != "" {
		log.Print("Starting with SSL")
		return http.ListenAndServeTLS(config.ListenAddr, config.TLS.Cert, config.TLS.Key, registry)
	}

	// No TLS configured
	log.Print("Warning: Server is starting without SSL, you should not pass any credentials using this configuration")
	log.Print("To use SSL, you must provide a config file with a `tls` section, and provide locations to a `key` file and a `cert` file")

	return http.ListenAndServe(config.ListenAddr, registry)
}
