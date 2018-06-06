package main


import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"bufio"
	"strings"
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



func main() {


	if len(os.Args) < 5 {
		fmt.Println("Please enter ip port user friend key(optional) ")
		return
	}

	ip := os.Args[1]
	port := os.Args[2]
	user := os.Args[3]
	file := os.Args[4]

	c, e := net.Dial("tcp", ip+":"+port)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	m := &CommonMessage{}

	decoder := json.NewDecoder(c)
	encoder := json.NewEncoder(c)

	err := decoder.Decode(&m)

	if err != nil {
		if err != nil {
			fmt.Println("Decode Main",err.Error())
			os.Exit(-1)
		}
	}
	cid := m.Conn

	m.Lgin = &Login{}
	m.Conn = cid
	m.Lgin.UserName = user
	err = encoder.Encode(&m)

	if err != nil {
		fmt.Println("Decode Main",err.Error())
		os.Exit(-1)
	}

	m1 := &CommonMessage{}
	err = decoder.Decode(&m1)

	if err != nil {
		fmt.Println("message",err.Error())
		os.Exit(-1)
	}
	fmt.Println(m1.Msg.From , m1.Msg.Data)

	inFile, err := os.Open(file)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(-1)
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		m := CommonMessage{}
		command := strings.Split(scanner.Text()," ")
		switch command[0] {
		case "message":
			m.Conn = cid
			m.Msg = &Message{}
			m.Msg.Data = command[2]
			m.Msg.To = command[1]
			m.Msg.From = user
			err = encoder.Encode(&m)
			if err != nil {
				fmt.Println("message",err.Error())
				os.Exit(-1)
			}
		//	time.Sleep(time.Millisecond*2)

		/*	m1 := &CommonMessage{}
			err := decoder.Decode(&m1)

			if err != nil {
				fmt.Println("message",err.Error())
				os.Exit(-1)
			}
			fmt.Println(m1.Msg.From , m1.Msg.Data)*/
		}
	}

	c.Close()

}
