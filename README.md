# slog

### What is slog?
Slog is a small command-line tool that launches your process and streams STDOUT and STDERR to a remote machine. It's as simple as running `slog "echo hello world"`.

### Installation
First make sure you have the [Go compiler](https://golang.org/doc/install) installed, then clone this repo and run `go build && mv slog /usr/bin/local/slog`.

### Configuration
Slog requires a few variables: *host*, *port*, *user*, *outpath*, *errpath*, and either *auth:pass* or *auth:pem* (the absolute path to your .pem key). You can pass any of these as flags on the command line or in a `.slog` file in your home directory.
