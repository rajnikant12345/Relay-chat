package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	//"strings"
	"encoding/json"
	"io"
	"strings"
	"sync"
)

var userConn map[string]net.Conn
var mu sync.RWMutex

func WriteMap(s string, c net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	userConn[s] = c

}

func DelMap(s string) {
	mu.Lock()
	defer mu.Unlock()
	delete(userConn, s)

}

func CleanSlate() {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range userConn {
		writeIt(v, "Server is cleaning slate...\n")
		v.Close()
		delete(userConn, k)
	}
}

func GetFromMap(s string) (net.Conn, bool) {
	mu.RLock()
	defer mu.RUnlock()
	c, ok := userConn[s]
	return c, ok
}

func writeIt(c net.Conn, s string) {
	w := bufio.NewWriter(c)
	w.WriteString(s)
	w.Flush()
}

func readit(c net.Conn) (string, error) {
	r := bufio.NewReader(c)
	s, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

type MessageStruct struct {
	From    string
	To      string
	Message string
}

func processMessage(c chan MessageStruct) {
	for m := range c {
		c1, ok := GetFromMap(m.From)
		if !ok {
			continue
		}
		c2, ok := GetFromMap(m.To)
		if !ok {
			j := json.NewEncoder(c1)
			m.Message = "Disconected"
			m.To = m.From
			m.From = "Server"
			j.Encode(&m)
			continue
		}
		fmt.Println(m.From , "=>",m.To ," : ", m.Message)
		j := json.NewEncoder(c2)
		j.Encode(&m)

	}
}

func main() {
	ch := make(chan MessageStruct, 100)
	mu = sync.RWMutex{}
	userConn = make(map[string]net.Conn)
	fmt.Println("Strating the chat server")

	port := os.Getenv("MY_SERVER_PORT")

	_, err := strconv.Atoi(port)

	if err != nil {
		fmt.Println("Enter a valid port")
		return
	}

	l, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	go func() {
		signalChan := make(chan os.Signal, 1)
		cleanupDone := make(chan bool)
		signal.Notify(signalChan, os.Interrupt)
		go func() {
			for _ = range signalChan {
				fmt.Println("\nReceived an interrupt, stopping services...\n")
				CleanSlate()
				cleanupDone <- true
			}
		}()
		<-cleanupDone
		os.Exit(0)
	}()

	for i := 0; i < 10; i++ {
		go processMessage(ch)
	}

	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error(), " Stopping the server..")
			c.Close()
			return
		}
		go func(net.Conn) {
			writeIt(c, "You are connected , please enter your user name: ")
			s, _ := readit(c)
			fmt.Println(s, " is connected")
			writeIt(c, s+": You are conected\n")
			WriteMap(s, c)
			for {
				m := MessageStruct{}
				j := json.NewDecoder(c)
				err := j.Decode(&m)
				if err != nil && err == io.EOF {
					fmt.Println(s, "disconnected")
					DelMap(s)
					break
				}
				if err != nil {
					writeIt(c, "Fail to decode message\n")
					continue
				}
				ch <- m
			}
		}(c)

	}
}
