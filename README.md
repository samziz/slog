# slog

### What is slog?
Slog is a small command-line tool that launches your process and sends your stdout and stderr streams to a remote machine over SSH. It's as simple as running `slog "echo hello world"`.

### Installation
The safest and easiest way to install slog is through the Makefile. This won't work for Windows since it relies on moving your package to `/usr/local/bin`. If you are brave enough to want to install this on Windows, or if you can't use the Makefile for any other reason, you should make sure you have the [Go compiler](https://golang.org/doc/install) installed, then clone this repo and run `go build && mv slog $HERE` where `$HERE` is somewhere in your path.

### Configuration
Slog requires a few variables: *host*, *port*, *user*, *outpath*, *errpath*, and either *auth:pass* or *auth:pem* (the absolute path to your .pem key). You can pass any of these as flags on the command line or in a `.slog` file in your home directory. If you choose to supply a `.slog` file, this should be in valid JSON. You can supply any of the following fields: 

```
{
	"host": "15.132.116.37",
	"port": "22",
	"user": "joe.bloggs",
	"outpath": "~/myproc/out",
	"errpath": "~/myproc/err",
	"loglevel": "verbose",
	"auth": {
		"pass": "hunter2",
		"pem": "~/secret/mysshkey.pem"
	}
}
```
All of these fields must be supplied except for `loglevel`, but they can be supplied as command line flags instead if you like. Additionally, after any flags, you should supply a command to execute. Any parameter expansion will be done on the remote machine rather than on yours. If you need to evaluate any variables in your command, you should do that beforehand and pipe it into slog using `xargs` (e.g. `parallel gcc grep */main.c | slog --outpath ~/myproc/$$-out --errpath ~/myproc/$$-err`). If you want to use slog as part of an automated workflow this is also the best way to achieve that.