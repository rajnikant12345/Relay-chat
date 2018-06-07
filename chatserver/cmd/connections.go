package cmd

import (
	"Relay-chat/chatserver/model"
	"Relay-chat/chatserver/processors"
	"encoding/json"
	"io"
	"net"
	"Relay-chat/chatserver/applog"
	"Relay-chat/chatserver/config"
)

type channelData struct {
	m       model.CommonMessage
	encoder net.Conn
}

func processMessage(ch chan channelData) {
	for m := range ch {
		p := processors.GetProcessor(m.m)
		if p!= nil {
			p.ProcessMessage(m.m , m.encoder)
		}
	}
}

func StartWorkers(ch chan channelData) {
	for i := 0; i < config.CFG.Workers; i++ {
		go processMessage(ch)
	}
}

//TODO: Handle timeout of an idle connection

func HandleConnections(c net.Conn, conn string) {

	applog.Info.Println("Starting Handle Connection for", conn)
	var ch chan channelData
	ch = make(chan channelData, config.CFG.ChannelSize)

	StartWorkers(ch)

	j := json.NewDecoder(c)

	for {
		m := model.CommonMessage{}
		e := j.Decode(&m)
		if e != nil {
			if e == io.EOF {
				applog.Warning.Println("Terminating connection for", conn, "ERROR:", e.Error())
				model.DeleteFromConnMap(conn)
				return
			} else {
				applog.Warning.Println("Terminating connection for", conn,"ERROR:", e.Error())
				model.DeleteFromConnMap(conn)
				return
			}
		}
		//validate conection id
		if m.Conn != conn {
			applog.Error.Println("Expected conid:",conn, "Receive Connid:",m.Conn)
			c.Close()
			break
		}

		cdata := channelData{m, c}

		ch <- cdata
	}
}
