package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io/ioutil"
	"log"
	"os/user"
)

type ConfigFile struct {
	Host string
	Port string
	User string
	Outpath string
	Errpath string
	Loglevel string
	Auth struct {
		Pass string
		Pem string
	}
}

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

	data, err := ioutil.ReadFile(u.HomeDir + "/.slog")
	if err != nil {
		log.Fatal(err)
	}

	var conf ConfigFile
	err = json.Unmarshal(data, &conf)
	
	if err != nil {
		log.Fatal(err)
	}

	// Load flags, using .slog file settings as defaults
	host := flag.String("host", conf.Host, "host address")
	port := flag.String("port", conf.Port, "port for remote (default of 22 should usually work)")
	user := flag.String("user", conf.User, "user for remote")
	outpath := flag.String("outpath", conf.Outpath, "absolute path for stdout")
	errpath := flag.String("errpath", conf.Errpath, "absolute path for stderr")
	authpem := flag.String("auth:pem", conf.Auth.Pass, "absolute path for .pem file")
	authpass := flag.String("auth:pass", conf.Auth.Pem, "ssh password for remote")
	loglevel := flag.String("loglevel", conf.Loglevel, "default is warning, else none or verbose")
	flag.Parse()

	// Get command to run
	args := flag.Args()
	if len(args) == 0 {
		return nil, errors.New("No command passed")
	}

	opts := map[string]string{
		"host": *host,
    	"port": *port,
		"user": *user,
		"outpath": *outpath,
		"errpath": *errpath,
		"auth:pem": *authpem,
		"auth:pass": *authpass,
		"loglevel": *loglevel,
		"cmd": args[0],
	}

	return opts, nil
}
