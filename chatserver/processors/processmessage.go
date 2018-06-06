package processors

import (
	"cryptolessons/chatserver/model"
	"encoding/json"
	"io"
	"net"
	"cryptolessons/chatserver/applog"
)

const fail = "-1"
const loginsuccess = "0"

func WriteMessage(c net.Conn, m model.CommonMessage, msg string, from, to string) {
	applog.Info.Println("WriteMessage","From:",from,"To:",to,"Message:",msg)
	encoder := json.NewEncoder(c)
	message := model.Message{}
	message.Data = msg
	message.From = from
	message.To = to
	m.Msg = &message
	e := encoder.Encode(&m)
	if e != nil && e == io.EOF {
		applog.Warning.Println("Encoder cannot send data")
		model.DeleteFromMap(to)
	}
}

func ProcessLoginMessage(m model.CommonMessage, c net.Conn) {
	applog.Info.Println("Processig Login Message")
	_, ok := model.ReadKey(m.Lgin.UserName)
	if ok {
		applog.Warning.Println("User",m.Lgin.UserName,"already connected")
		WriteMessage(c, m, fail, "chat-server", m.Lgin.UserName)
		return
	}
	applog.Info.Println("Login Success",m.Lgin.UserName)
	model.WriteMap(m.Lgin.UserName, model.Connection{m.Conn, c})
	WriteMessage(c, m, loginsuccess, "chat-server", m.Lgin.UserName)
}

func ProcessKeyExchange(m model.CommonMessage, c net.Conn) {
	applog.Info.Println("Processig Key Exchange Message")
	_, ok := model.ReadKey(m.KeyExchg.From)
	if !ok {
		applog.Warning.Println("User",m.KeyExchg.From,"not connected")
		WriteMessage(c, m, fail, "chat-server", m.KeyExchg.From)
		return
	}
	val, ok := model.ReadKey(m.KeyExchg.To)
	if !ok {
		applog.Warning.Println("User",m.KeyExchg.To,"not connected")
		WriteMessage(c, m, fail, "chat-server", m.KeyExchg.From)
		return
	}
	encoder := json.NewEncoder(val.C)
	e := encoder.Encode(&m)
	if e != nil && e == io.EOF {
		applog.Warning.Println("Encoder cannot send data")
		model.DeleteFromMap(m.KeyExchg.To)
	}
	applog.Info.Println("Key Exchange", "From:",m.KeyExchg.From,"To:",m.KeyExchg.To,"Message:",m.KeyExchg.Key)

}

func ProcessMessage(m model.CommonMessage, c net.Conn) {

	p := 0
	applog.Info.Println("Processig Message")

	if m.Lgin != nil {
		p++
		ProcessLoginMessage(m, c)
	}
	if m.KeyExchg != nil {
		if p != 0 {
			applog.Warning.Println("Already processed Login Message")
			return
		}
		p++
		ProcessKeyExchange(m, c)
	}
	if m.Msg != nil {
		if p != 0 {
			applog.Warning.Println("Already processed Login Message or Key exchange")
			return
		}
		_, ok := model.ReadKey(m.Msg.From)
		if !ok {
			applog.Warning.Println("User",m.Msg.From,"not connected")
			WriteMessage(c, m, fail, "chat-server", m.Msg.From)
			return
		}
		v, ok := model.ReadKey(m.Msg.To)
		if !ok {
			applog.Warning.Println("User",m.Msg.To,"not connected")
			WriteMessage(c, m, fail, "chat-server", m.Msg.From)
			return
		}
		encoder := json.NewEncoder(v.C)
		e := encoder.Encode(&m)
		if e != nil && e == io.EOF {
			applog.Warning.Println("Encoder cannot send data")
			model.DeleteFromMap(m.Msg.To)
		}
		applog.Info.Println("From:",m.Msg.From,"To:",m.Msg.To,"Message:",m.Msg.Data)
	}
}
