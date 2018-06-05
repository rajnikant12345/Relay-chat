package cmd

import (
	"crypto/tls"
	"cryptolessons/chatserver/config"
	"net"
	"os"
	"log"
	"cryptolessons/chatserver/processors"
	"cryptolessons/chatserver/model"
	"encoding/json"
)

func BeginTLS(key, cert, port string) (net.Listener, error) {

	if key == "" || cert == "" {
		return BeginTCP(port)
	}

	cer, err := tls.LoadX509KeyPair(key, cert)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	config := &tls.Config{Certificates: []tls.Certificate{cer}}
	ln, err := tls.Listen("tcp", ":"+port, config)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ln, nil
}

func BeginTCP(port string) (net.Listener, error) {
	ln, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return ln, nil
}

func StartServer() {
	cfg, e := config.GetConfig(os.Getenv("APP_CFG"))
	if e != nil {
		log.Println(e.Error())
		return
	}
	l, e := BeginTLS(cfg.Key, cfg.Cert, cfg.Port)
	if e != nil {
		log.Println(e.Error())
	}

	StartWorkers()

	for {
		c, e := l.Accept()
		if e != nil {
			log.Println(e.Error())
			continue
		}

		m := model.CommonMessage{}
		m.Ref = processors.GenerateReference()
		m.Conn = processors.GenerateConnectinId()

		enc := json.NewEncoder(c)
		e = enc.Encode(&m)
		if e != nil {
			log.Println(e.Error())
			continue
		}
		go HandleConnections(c , m.Conn , m.Ref )
	}
}

