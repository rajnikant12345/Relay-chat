package cmd

import (
	"crypto/tls"
	"cryptolessons/chatserver/config"
	"cryptolessons/chatserver/model"
	"cryptolessons/chatserver/processors"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"cryptolessons/chatserver/applog"
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
	cfg := config.CFG

	l, e := BeginTLS(cfg.Key, cfg.Cert, cfg.Port)
	if e != nil {
		log.Println(e.Error())
	}

	go func() {
		signalChan := make(chan os.Signal, 1)
		cleanupDone := make(chan bool)
		signal.Notify(signalChan, os.Interrupt)
		go func() {
			for _ = range signalChan {
				applog.Warning.Println("Received an interrupt, stopping services...")
				model.ClearMap()
				cleanupDone <- true
			}
		}()
		<-cleanupDone
		os.Exit(0)
	}()

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

			applog.Error.Println("Accept Loop",e.Error())
			continue
		}
		go HandleConnections(c, m.Conn)
	}
}
