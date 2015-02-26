// Package webhookListener Handler support
package webhookListener

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"text/template"
)

// Endpoint is the configuration of a path of our HTTP server including security,
// type of message and commands to execute
type Endpoint struct {
	MessageType string `json:"messageType"`
	Apikey      string `json:"apiKey"`
	Path        string
	CommandDir  string `json:"commandDir"`
	Commands    [][]string
}

// Registry represent a collection of endpoint and there http.Handler
type Registry interface {
	Add(endpoint Endpoint)
	ServeHTTP(res http.ResponseWriter, req *http.Request)
}

type registryStruct struct {
	endpoints map[string]Endpoint
}

type httpError struct {
	code    int
	message string
}

// CreateRegistry will allocate a new registry and its endpoints
func CreateRegistry() Registry {
	var reg registryStruct

	reg.endpoints = make(map[string]Endpoint)

	return &reg
}

// Add an endpoint to the registry
func (reg *registryStruct) Add(endpoint Endpoint) {
	reg.endpoints[endpoint.Path] = endpoint

	return
}

func (reg *registryStruct) ServeHTTP(res http.ResponseWriter, req *http.Request) {
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

	log.Printf("%+v", message)

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

		for key := range params {
			params[key] = parseTemplate(params[key], msg)
		}

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

func parseTemplate(tpl string, value interface{}) string {
	commandTpl := template.Must(template.New("letter").Parse(tpl))
	var b bytes.Buffer
	commandTpl.Execute(&b, value)
	command := b.String()

	return command
}

func run(dir string, cmd *exec.Cmd) {
	cmd.Dir = dir
	err := cmd.Run()
	if err != nil {
		fmt.Print(err)
	}
}
