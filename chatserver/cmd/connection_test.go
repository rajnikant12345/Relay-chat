package cmd

import (
	"cryptolessons/chatserver/model"
	"testing"
	"net"
	"encoding/json"
	"fmt"
	"cryptolessons/chatserver/processors"

	"sync"
	"time"
)


func Decoder(c net.Conn) {
	dec := json.NewDecoder(c)
	for {
		m := model.CommonMessage{}
		e := dec.Decode(&m)
		if e == nil {
			fmt.Println(m.Msg.From , m.Msg.To , m.Msg.Data)
		}

	}
}




func TestHandleConnections(t *testing.T) {

	users := make([]string,2)
	cid := make([]string,10000)
	conn := make([]net.Conn,10000)
	var wg sync.WaitGroup

	for i,_ := range users {

		//wg.Add(1)
		 func(i int) {
			//defer wg.Done()
			users[i] = processors.GenerateConnectinId()
			c,e := net.Dial("tcp","localhost:6789")
			if e != nil {
				fmt.Print(e.Error())
				t.Fail()
				return
			}
			conn[i] = c
			dec := json.NewDecoder(c)
			enc := json.NewEncoder(c)
			m := model.CommonMessage{}
			dec.Decode(&m)

			cid[i] = m.Conn

			 fmt.Println(m.Conn)
			m.Lgin = &model.Login{}
			m.Lgin.UserName = users[i]
			enc.Encode(&m)

			go Decoder(c)
		}(i)
	}

	for i,_ := range users {

		for j,_ := range users {
			wg.Add(1)
			 go func(i , j int) {
				defer wg.Done()
				if i == j {
					return
				}
				enc := json.NewEncoder(conn[i])
				m := model.CommonMessage{}
				m.Conn = cid[i]
				m.Msg = &model.Message{}
				m.Msg.Data = "hello "
				m.Msg.From = users[i]
				m.Msg.To = users[j]
				enc.Encode(&m)
			}(i,j)

		}

	}
	wg.Wait()

	time.Sleep(time.Second*5)

	for i,_ := range users {
		if conn[i] != nil {
			conn[i].Close()
		}

	}


}
