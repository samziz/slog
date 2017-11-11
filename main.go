package main

import (
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os/user"
	"strings"
)

func main() {
	conf, err := getConfig()
	if err != nil {
		log.Fatal(err)
	}

	conn := Conn{
		Host: conf["host"],
    	Port: conf["port"],
    	User: conf["user"],
		OutPath: conf["outpath"],
		ErrPath: conf["errpath"],
		Cmd: conf["cmd"],
		Conf: conf,
	}

	err = conn.Dial()
	if err != nil {
		log.Fatal(err)
	}

	conn.ExecAndListen()
}

func getConfig() (map[string]string, error) {
	// Load .slog file 
	u, err := user.Current()
	if err != nil {
		return nil, err
	}

	var df map[string]string
	home := u.HomeDir
	data, err := ioutil.ReadFile(home + "/.slog")
	if err == nil {
		df, err = parseFileToMap(data)

		if err != nil {
			return nil, err
		}
	}

	// Load flags, using .slog file settings as defaults
	host := flag.String("host", df["host"], "host address")
	port := flag.String("port", "22", "port for remote (default of 22 should usually work)")
	user := flag.String("user", df["user"], "user for remote")
	outpath := flag.String("outpath", df["outpath"], "absolute path for stdout")
	errpath := flag.String("errpath", df["errpath"], "absolute path for stderr")
	authpem := flag.String("authpem", df["auth:pem"], "absolute path for .pem file")
	authpass := flag.String("authpass", df["auth:pass"], "ssh password for remote")
	loglevel := flag.String("loglevel", df["loglevel"], "default is warning, else none or verbose")
	flag.Parse()

	args := flag.Args()
	if len(args) == 0 {
		return nil, errors.New("No command passed")
	}
	
	cmd := args[0]

	conf := map[string]string{
		"host": *host,
    	"port": *port,
		"user": *user,
		"outpath": *outpath,
		"errpath": *errpath,
		"auth:pem": *authpem,
		"auth:pass": *authpass,
		"loglevel": *loglevel,
		"cmd": cmd,
	}

	return conf, nil
}

func parseFileToMap(data []byte) (map[string]string, error) {
	s := string(data)
	lns := strings.Split(s, "\n")

	opts := map[string]string{}

	for _, ln := range lns {
		if ln != "" {
			kv := strings.Split(ln, "=")
			opts[kv[0]] = kv[1]
		}
	}

	return opts, nil
}
