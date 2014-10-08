// server.go
package webhookListener

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Config struct {
	ListenAddr string
	Tls        struct {
		Key  string
		Cert string
	}
	Apikeys struct {
		Key []string
	}
}

type GitlabPushMessage struct {
	Before            string
	After             string
	Ref               string
	UserId            int    `json:"user_id"`
	UserName          string `json:"user_name"`
	ProjectId         int    `json:"project_id"`
	TotalCommitsCount int    `json:"total_commits_count"`
	Repository        struct {
		Name        string
		Url         string
		Description string
		Homepage    string
	}
	Commits []struct {
		Id        string
		Message   string
		Timestamp string
		Url       string
		Author    struct {
			Name  string
			Email string
		}
	}
}

var ServerConfig *Config
var messageHandlers Registry

func Serve(config *Config) error {
	ServerConfig = config

	if len(ServerConfig.Apikeys.Key) == 0 {
		log.Print("Warning: The server is about to start without any authentication.  Anyone can trigger handlers to fire off")
		log.Print("To enable authentication, you must add an `apikeys` section with at least 1 `key`")
	}

	messageHandlers = MessageHandlers()
	messageHandlers.Add((&DnsHandler{}).Call)
	http.HandleFunc("/", reqHandler)

	// Manage TLS
	if config.Tls.Key != "" && config.Tls.Cert != "" {
		log.Print("Starting with SSL")
		return http.ListenAndServeTLS(config.ListenAddr, config.Tls.Cert, config.Tls.Key, Log(http.DefaultServeMux))
	}

	// No TLS configured
	log.Print("Warning: Server is starting without SSL, you should not pass any credentials using this configuration")
	log.Print("To use SSL, you must provide a config file with a [tls] section, and provide locations to a `key` file and a `cert` file")

	return http.ListenAndServe(config.ListenAddr, Log(http.DefaultServeMux))
}

func reqHandler(w http.ResponseWriter, r *http.Request) {
	if requestAuthenticated(r) {
		decoder := json.NewDecoder(r.Body)
		var message GitlabPushMessage

		err := decoder.Decode(&message)
		if err != nil {
			http.Error(w, "Could not decode json", 500)
			log.Print(err)
			return
		}

		fmt.Fprintf(w, "%+v", message)
		go handleMessage(message)
		return
	}
	http.Error(w, "Not Authorized", 401)
}

func requestAuthenticated(r *http.Request) bool {
	key := r.URL.Query().Get("key")

	for _, configKey := range ServerConfig.Apikeys.Key {
		if key == configKey {
			return true
		}
	}

	return (len(ServerConfig.Apikeys.Key) == 0)
}

func handleMessage(msg GitlabPushMessage) {
	messageHandlers.Call(msg)
}

func Log(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s", r.RemoteAddr, r.Method)
		handler.ServeHTTP(w, r)
	})
}
