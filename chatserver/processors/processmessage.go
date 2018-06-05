package processors

import (
	"cryptolessons/chatserver/model"
	"encoding/json"
	"io"
	"log"
	"net"
)

const fail = "-1"
const loginsuccess = "0"

func WriteMessage(c net.Conn, m model.CommonMessage, msg string, from, to string) {
	encoder := json.NewEncoder(c)
	message := model.Message{}
	message.Data = msg
	message.From = from
	message.To = to
	m.Msg = &message
	e := encoder.Encode(&m)
	if e != nil && e == io.EOF {
		model.DeleteFromMap(to)
	}
}

func ProcessLoginMessage(m model.CommonMessage, c net.Conn) {
	log.Println("ProcessLoginMessage")
	_, ok := model.ReadKey(m.Lgin.UserName)
	if ok {
		log.Println("User already connected.")
		WriteMessage(c, m, fail, "chat-server", m.Lgin.UserName)
		return
	}
	log.Println("Login Success.")
	model.WriteMap(m.Lgin.UserName, model.Connection{m.Conn, c})
	WriteMessage(c, m, loginsuccess, "chat-server", m.Lgin.UserName)
}

func ProcessKeyExchange(m model.CommonMessage, c net.Conn) {
	_, ok := model.ReadKey(m.KeyExchg.From)
	if !ok {
		log.Println("User not logged in:", m.KeyExchg.From)
		WriteMessage(c, m, fail, "chat-server", m.KeyExchg.From)
		return
	}
	val, ok := model.ReadKey(m.KeyExchg.To)
	if !ok {
		log.Println("Second User not logged in:", m.KeyExchg.To)
		WriteMessage(c, m, fail, "chat-server", m.KeyExchg.From)
		return
	}
	encoder := json.NewEncoder(val.C)
	e := encoder.Encode(&m)
	if e != nil && e == io.EOF {
		model.DeleteFromMap(m.KeyExchg.To)
	}

}

func ProcessMessage(m model.CommonMessage, c net.Conn) {

	p := 0
	log.Println("ProcessMessage")

	if m.Lgin != nil {
		p++

		ProcessLoginMessage(m, c)
	}
	if m.KeyExchg != nil {
		if p != 0 {
			return
		}
		p++
		ProcessKeyExchange(m, c)
	}
	if m.Msg != nil {
		if p != 0 {
			return
		}
		_, ok := model.ReadKey(m.Msg.From)
		if !ok {
			log.Println("User not logged in:", m.Msg.From)
			WriteMessage(c, m, fail, "chat-server", m.Msg.From)
			return
		}
		v, ok := model.ReadKey(m.Msg.To)
		if !ok {
			log.Println("Second User not logged in:", m.Msg.To)
			WriteMessage(c, m, fail, "chat-server", m.Msg.From)
			return
		}
		encoder := json.NewEncoder(v.C)
		e := encoder.Encode(&m)
		if e != nil && e == io.EOF {
			model.DeleteFromMap(m.Msg.To)
		}
	}
}
