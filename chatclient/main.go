package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
)

type MessageStruct struct {
	From    string
	To      string
	Message string
}

func main() {
	var ip, port, user, friend string
	fmt.Print("Enter IP of Server:")
	fmt.Scanf("%s\n", &ip)
	fmt.Print("Enter port of Server:")
	fmt.Scanf("%s\n", &port)

	c, e := net.Dial("tcp", ip+":"+port)
	if e != nil {
		fmt.Println(e.Error())
		return
	}
	w := bufio.NewWriter(c)
	r := bufio.NewReader(c)
	s, _ := r.ReadString(':')
	fmt.Println(s)

	fmt.Print("\nEnter user of Server:")
	fmt.Scanf("%s\n", &user)

	w.WriteString(user + "\n")
	w.Flush()
	fmt.Print("\nEnter friend name:")
	fmt.Scanf("%s\n", &friend)

	s, _ = r.ReadString('\n')
	fmt.Println(s)

	go func() {
		for {
			m := MessageStruct{}
			j := json.NewDecoder(c)
			err := j.Decode(&m)
			if err != nil {
				fmt.Println("Message error", err.Error())
				continue
			}
			fmt.Println("Message From", m.From, ":", m.Message)
		}
	}()

	for {
		m := MessageStruct{}
		in := bufio.NewReader(os.Stdin)
		m.Message, _ = in.ReadString('\n')
		m.From = user
		m.To = friend
		j := json.NewEncoder(c)
		j.Encode(&m)
	}

}
