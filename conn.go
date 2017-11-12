package main

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io"
	"io/ioutil"
	"log"
	"os/exec"
	"path"
	"strings"
)

type Conn struct {
	Host     string
	Port	 string
	User     string
	OutPath  string
	ErrPath  string
	Cmd      string
	Conf     map[string]string
	Client 	 *ssh.Client
	CmdSession *ssh.Session
}


/* Exported funcs */

func (c *Conn) Log(errlv string, errs ...string) {
	err := strings.Join(errs, "")

	lv, exists := c.Conf["loglevel"]
	if !exists {
		lv = "warning"
	}

	switch errlv {
	case "warning":
		if lv != "none" {
			fmt.Println("WARNING:", err)
		}
	case "verbose":
		if lv == "verbose" {
			fmt.Println("VERBOSE:", err)
		}
	}
}

func (c *Conn) Dial() error {
	auth, err := c.GetAuthMethods()
	if err != nil {
		return AuthError{Body: err.Error()}
	}

	cc := ssh.ClientConfig{
		User: c.User,
		Timeout: 0,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //FIXME
		Auth: auth,
	}

	cl, err := ssh.Dial("tcp", c.Host + ":" + c.Port, &cc)
	if err != nil {
		return AuthError{Body: err.Error()}
	}
	
	c.Log("verbose", "Connected to host:", c.Host)
	c.Client = cl
	return nil
}

func (c *Conn) GetAuthMethods() ([]ssh.AuthMethod, error) {
	var arr []ssh.AuthMethod

	if val, exists := c.Conf["auth:pass"]; exists {
		m := ssh.Password(val)
		arr = append(arr, m)
	}

	if val, exists := c.Conf["auth:pem"]; exists {
		s, err := c.getSignerFromPem(val)
		if err != nil {
			return nil, err
		}

		m := ssh.PublicKeys(s)
    	arr = append(arr, m)

    	c.Log("verbose", "Loaded .pem file from:", val)
	}

	if len(arr) == 0 {
		return nil, AuthError{Body: "No auth methods supplied"}
	}

	return arr, nil
}

func (c *Conn) ExecAndListen() error {	
		
	cmd := exec.Command("sh", "-c", c.Cmd)
	stdout, err := cmd.StdoutPipe()
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return ConnectionError{Body: err.Error()}
	}

	c.Log("verbose", "Executing command:", c.Cmd)
	
	err = cmd.Start()
	
	if err != nil {
		return ConnectionError{Body: err.Error()}
	}

	go c.streamToRemoteFile(stdout, c.OutPath)
	go c.streamToRemoteFile(stderr, c.ErrPath)

	for {
		
	}
}


/* Unexported funcs */

// Load .pem file at path and return Signer
func (c *Conn) getSignerFromPem(path string) (ssh.Signer, error) {
	f, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, OSError{Body: "Could not find .pem file at path:" + path}
	}

	return ssh.ParsePrivateKey(f)
}

// Read from Reader at loc and write data to remote file over SSH
func (c *Conn) streamToRemoteFile(r io.Reader, dst string) {
	zeros := make([]byte, 256)

	// Init directory and file
	dir := path.Dir(dst)
	s, err := c.Client.NewSession()

	// $ mkdir '$dir' && touch $dist
	err = s.Run("mkdir -p " + dir + "&& touch " + dst)
	
	if err != nil {
		cerr := CopyError{Body: err.Error()}
		log.Fatal(cerr)
	}

	s.Close()
	
	// Stream data
	for {
		buf := make([]byte, 256)
		_, err := r.Read(buf)

		if err != io.EOF {
			str := string(buf)
			s, _ := c.Client.NewSession()

			if !bytes.Equal(buf, zeros) {
				out, err := s.CombinedOutput("cat << SSHEOF >> " + dst + "\n" + str + "\nSSHEOF")
				
				if err != nil {
					cerr := CopyError{Body: err.Error()}
					c.Log("verbose", cerr.Error(), string(out))
				}
			}
			
			s.Close()
		}
	}
}