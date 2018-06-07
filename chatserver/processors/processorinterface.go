package processors

import (
	"Relay-chat/chatserver/model"
	"net"
)

type Processor interface {
	ProcessMessage(m model.CommonMessage, c net.Conn)
}
