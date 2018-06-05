package cmd

import (
	"cryptolessons/chatserver/model"
	"os"
	"testing"
	"net"
	"encoding/json"
	"fmt"

	"cryptolessons/chatserver/processors"

)

func startserver() {
	os.Setenv("APP_CFG", "/Users/rajnikant/workspace/src/cryptolessons/app.json")
	go StartServer()
}




func TestHandleConnections(t *testing.T) {

	users := make([]string,100)
	cid := make([]string,10000)
	conn := make([]net.Conn,10000)
//	var wg sync.WaitGroup
	//wg.Add(100)
	for i,_ := range users {


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

			m.Lgin = &model.Login{}
			m.Lgin.UserName = users[i]
			enc.Encode(&m)
			dec.Decode(&m)

		}(i)
	}



	//wg.Wait()

	//wg.Add(10000)

	for i,_ := range users {
		for j,_ := range users {
			//wg.Add(1)
			func(i , j int) {
				//defer wg.Done()
				//fmt.Println(i,j)
				if i == j {
					return
				}
				enc := json.NewEncoder(conn[i])
				m := model.CommonMessage{}
				m.Conn = cid[i]
				//fmt.Println(cid[i])
				m.Msg = &model.Message{}
				m.Msg.Data = "hello"
				m.Msg.From = users[i]
				m.Msg.To = users[j]
				//mm,_ :=  json.Marshal(&m)
				//fmt.Println(string(mm))
				enc.Encode(&m)
				dec := json.NewDecoder(conn[j])
				dec.Decode(&m)
				//fmt.Println(j)
			}(i,j)

		}

	}

	//wg.Wait()

}
