package processors

import (
	"Relay-chat/chatserver/applog"
	"Relay-chat/chatserver/model"
	"net"
)

type LoginProcessor struct {
}

func (l *LoginProcessor) ProcessMessage(m model.CommonMessage, c net.Conn) {
	applog.Info.Println("Processig Login Message")
	_, ok := model.ReadKey(m.Lgin.UserName)
	if ok {
		applog.Warning.Println("User", m.Lgin.UserName, "already connected")
		WriteMessage(c, m, fail, "chat-server", m.Lgin.UserName)
		return
	}
	applog.Info.Println("Login Success", m.Lgin.UserName)
	model.WriteMap(m.Lgin.UserName, model.Connection{Connid: m.Conn, Con: c})
	WriteMessage(c, m, loginsuccess, "chat-server", m.Lgin.UserName)
}
