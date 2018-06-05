package model

import (
	"net"
	"sync"
)

type CommonMessage struct {
	Conn     string       `json:"conn,omitempty"`
	Ref      string       `json:"ref,omitempty"`
	KeyExchg *KeyExchange `json:"key_exchg,omitempty"`
	Msg      *Message     `json:"message,omitempty"`
	Lgin     *Login       `json:"login,omitempty"`
}

type Login struct {
	UserName     string `json:"user_name,omitempty"`
	PasswordHash string `json:"password_hash,omitempty"`
}

type Message struct {
	To   string `json:"to,omitempty"`
	From string `json:"from,omitempty"`
	Data string `json:"data,omitempty"`
}

type KeyExchange struct {
	From string `json:"from,omitempty"`
	To   string `json:"to,omitempty"`
	Key  string `json:"key,omitempty"`
}

var mutex sync.RWMutex

var userMap map[string]Connection
var connMap map[string]string

type Connection struct {
	Connid string
	C      net.Conn
}

func init() {
	userMap = make(map[string]Connection)
	connMap = make(map[string]string)
}

func WriteMap(key string, value Connection) {
	mutex.Lock()
	defer mutex.Unlock()
	userMap[key] = value
	//connMap[value.Connid] = key
}

/*
func DeleteFromConnMap(key string) {
	mutex.Lock()
	defer mutex.Unlock()
	v, ok := connMap[key]
	if !ok {
		return
	}
	userMap[key].C.Close()
	delete(userMap, v)
	delete(connMap, key)
}
*/

func ReadKey(key string) (Connection, bool) {
	mutex.RLock()
	defer mutex.RUnlock()
	c, ok := userMap[key]
	return c, ok
}

func DeleteFromMap(key string) {
	mutex.Lock()
	defer mutex.Unlock()
	_, ok := userMap[key]
	if !ok {
		return
	}
	userMap[key].C.Close()
	delete(userMap, key)
}

func ClearMap() {
	mutex.Lock()
	defer mutex.Unlock()
	for k, v := range userMap {
		v.C.Close()
		delete(userMap, k)
	}

}
