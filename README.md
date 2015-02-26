# WebHook Listener

A simple HTTP(s) server that execute configured action when a message is received.

## Build Instruction

```sh
go get github.com/mitchellh/mapstructure
go get github.com/fmonniot/webhook-listener
go install github.com/fmonniot/webhook-listener/cli
```

## Usage

Launch it with `$GOPATH/bin/webhook-listener`. There are two optional arguments:

- `-config=""`: location of the config file
- `-listen="localhost:8080"`: `<address>:<port>` to listen on

## Configuration

You can configure the server by passing a configuration file written in json via
the argument `config`.

This file is composed of three section:

- `listenAddr`: where the server should listen, format: `addr:port`
- `tls` TLS information if you want to use HTTPS
  - `key`: path to your certificate key
  - `cert`: path to your certificate cert
- `endpoints`: an array of all endpoints of the server

An endpoint consist of:

- A `messageType` that represent the type of message it can receive (currently only
  `GitlabPushMessage` is accepted).
- A `path` which is the path used by the server to know what endpoint it is
- An `apiKey` which will be needed (either as a GET params or a header field)
- A `commandDir` which represent the working directory of the commands executed
- A `commands` which consist of arrays of commands (the first entry is the
  command and the other are the arguments). These commands are processed and you
  can provide message parameters following [the go text.template][1] syntax.

An example of configuration can be found in the root directory of this repository.

[1]: https://golang.org/pkg/text/template/
