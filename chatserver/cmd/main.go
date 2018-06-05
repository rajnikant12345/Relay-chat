package main

import (
	"crypto/tls"
	"cryptolessons/chatserver"
	"cryptolessons/chatserver/config"
	"fmt"
	"net"
	"os"
)

func BeginTLS(key, cert, port string) (net.Listener, error) {

	if key == "" || cert == "" {
		BeginTCP(port)
	}

	cer, err := tls.LoadX509KeyPair(key, cert)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":"+port, config)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ln, nil
}

func BeginTCP(port string) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}
	return ln, nil
}

func main() {
	cfg := cmd.GetConfig(os.Getenv("APP_CFG"))
	l, e := BeginTLS(cfg.Key, cfg.Cert, cfg.Port)
	if e != nil {
		fmt.Print(e.Error())
	}

	for {
		c, e := l.Accept()
		if e != nil {
			fmt.Println(e.Error())
			continue
		}
		go chatserver.HandleConnections(c)
	}

}
