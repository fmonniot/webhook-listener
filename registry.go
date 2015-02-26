// Package webhookListener Handler support
package webhookListener

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

// Endpoint TODO
type Endpoint struct {
	MessageType string `json:"messageType"`
	Apikey      string `json:"apiKey"`
	Path        string
	CommandDir  string `json:"commandDir"`
	Commands    [][]string
}

// Registry TODO
type Registry struct {
	endpoints map[string]Endpoint
}

type httpError struct {
	code    int
	message string
}

// CreateRegistry TODO
func CreateRegistry() Registry {
	var reg Registry

	reg.endpoints = make(map[string]Endpoint)

	return reg
}

// Add a handler for message
func (reg *Registry) Add(key string, endpoint Endpoint) {
	reg.endpoints[endpoint.Path] = endpoint

	log.Printf("%#v", reg)
	return
}

func (reg *Registry) ServeHTTP(res http.ResponseWriter, req *http.Request) {
	log.Printf("%s %s %s", req.RemoteAddr, req.Method, req.URL.Path)

	// Route to correct endpoint
	endpoint, exist := reg.endpoints[req.URL.Path]

	if !exist {
		http.NotFound(res, req)
		return
	}

	// Authenticate correct APIKey (either via Header or Get param)
	authErr := endpoint.authenticate(req)
	if authErr != nil {
		http.Error(res, authErr.message, authErr.code)
	}

	// Convert body to generic map
	var message map[string]interface{}
	httpErr := decodeJSON(&message, req.Body)
	if httpErr != nil {
		http.Error(res, httpErr.message, httpErr.code)
		return
	}

	// Convert generic map to struct
	msg, err := DecodeMessage(message)
	if err != nil {
		http.Error(res, "Wrong json format", 400)
		return
	}

	// Use msg, will be used in commands later
	log.Printf("Parsed JSON: %v", msg)

	// Exec	 commands
	for _, params := range endpoint.Commands {
		// TODO parse command and potentialy inject data from message
		cmd := exec.Command(params[0], params[1:]...)
		run(endpoint.CommandDir, cmd)
	}

	fmt.Fprint(res, "OK")
}

func (endpoint Endpoint) authenticate(req *http.Request) *httpError {
	key := req.Header.Get("apiKey")

	if key == "" {
		key = req.URL.Query().Get("apiKey")
	}

	if key == "" {
		return &httpError{403, "No API Key given"}
	}

	if key != endpoint.Apikey {
		return &httpError{403, "Wrong API Key"}
	}

	return nil
}

func decodeJSON(message *map[string]interface{}, body io.Reader) *httpError {
	decoder := json.NewDecoder(body)
	err := decoder.Decode(&message)
	if err != nil {
		return &httpError{500, "Could not decode json"}
	}
	return nil
}

func run(dir string, cmd *exec.Cmd) {
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		fmt.Print(err)
	}
}