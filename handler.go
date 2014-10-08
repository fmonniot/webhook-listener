// handler.go
package webhookListener

import "log"

type Handler interface {
	Call(GitlabPushMessage)
}

type Logger struct{}

func (l *Logger) Call(msg GitlabPushMessage) {
	log.Print(msg)
}

type Registry struct {
	entries []func(GitlabPushMessage)
}

func (r *Registry) Add(h func(msg GitlabPushMessage)) {
	r.entries = append(r.entries, h)
	return
}

func (r *Registry) Call(msg GitlabPushMessage) {
	for _, h := range r.entries {
		go h(msg)
	}
}

func MessageHandlers() Registry {
	var handlers Registry

	return handlers
}
