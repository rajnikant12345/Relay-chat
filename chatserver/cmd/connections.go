package cmd

import (
	"cryptolessons/chatserver/model"
	"cryptolessons/chatserver/processors"
	"encoding/json"
	"io"
	"log"
	"net"
)

type channelData struct {
	m       model.CommonMessage
	encoder net.Conn
}

func processMessage(ch chan channelData) {
	for m := range ch {
		processors.ProcessMessage(m.m, m.encoder)
	}
}

func StartWorkers(ch chan channelData) {
	for i := 0; i < 10; i++ {
		go processMessage(ch)
	}
}

//TODO: Handle timeout of an idle connection

func HandleConnections(c net.Conn, conn string) {

	var ch chan channelData
	ch = make(chan channelData, 1000)

	StartWorkers(ch)

	j := json.NewDecoder(c)

	for {
		m := model.CommonMessage{}
		e := j.Decode(&m)
		if e != nil {
			if e == io.EOF {
				model.DeleteFromConnMap(conn)
				log.Println("Conection terminated.")
				return
			} else {
				model.DeleteFromConnMap(conn)
				log.Println("Handle Connection, closing ", e.Error())
				return
			}
		}
		//validate conection id
		if m.Conn != conn {
			c.Close()
			break
		}

		cdata := channelData{m, c}

		ch <- cdata
	}
}
