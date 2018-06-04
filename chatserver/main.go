package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strconv"
	"encoding/json"
	"io"
	"strings"
	"sync"
)

//This map saves existing user and connection mapping
var userConn map[string]net.Conn

//This is a read write lock, for userConn map
var mu sync.RWMutex

// Write to map
func WriteMap(s string, c net.Conn) {
	mu.Lock()
	defer mu.Unlock()
	userConn[s] = c

}

//Delete entry from map
func DelMap(s string) {
	mu.Lock()
	defer mu.Unlock()
	delete(userConn, s)

}

// Clean all the connections from map
func CleanSlate() {
	mu.Lock()
	defer mu.Unlock()
	for k, v := range userConn {
		writeIt(v, "Server is cleaning slate...\n")
		v.Close()
		delete(userConn, k)
	}
}

// Read value from map
func GetFromMap(s string) (net.Conn, bool) {
	mu.RLock()
	defer mu.RUnlock()
	c, ok := userConn[s]
	return c, ok
}

//Write value to connection
func writeIt(c net.Conn, s string) {
	w := bufio.NewWriter(c)
	w.WriteString(s)
	w.Flush()
}

//Read value from connection
func readit(c net.Conn) (string, error) {
	r := bufio.NewReader(c)
	s, err := r.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(s), nil
}

//This struct translates to JSON and server expects same fromat from clients. refer README.md
type MessageStruct struct {
	From    string
	To      string
	Message string
}


// This is worker go routine
func processMessage(c chan MessageStruct) {
	// While bufferd channel c is nor closed, iterate channel and m has iterated value m is of type MessageStruct
	for m := range c {
		//Read sender's name from m
		c1, ok := GetFromMap(m.From)
		if !ok {
			continue
		}
		//Read receiver's name from m
		c2, ok := GetFromMap(m.To)
		if !ok {
			fmt.Println("socker", c1)
			j := json.NewEncoder(c1)
			m.Message = "Disconected"
			m.To = m.From
			m.From = "Server"
			err := j.Encode(&m)
			if err != nil {
				fmt.Println(err.Error())
			}
			continue
		}
		//send message to receiver
		fmt.Println(m.From, "=>", m.To, " : ", m.Message)
		j := json.NewEncoder(c2)
		j.Encode(&m)

	}
}

//This is main function
func main() {
	//create a buffered channel
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

	//start listening
	l, err := net.Listen("tcp", "0.0.0.0:"+port)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// this go routine handles ctrl + c
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

	// I am launching 10 go routines which act as worker, TODO: this value can be configurable.
	for i := 0; i < 10; i++ {
		go processMessage(ch)
	}

	for {
		//accept connection
		c, err := l.Accept()
		if err != nil {
			fmt.Println(err.Error(), " Stopping the server..")
			c.Close()
			return
		}
		
		// Launch a go routine per connection.
		go func(net.Conn) {
			writeIt(c, "You are connected , please enter your user name: ")
			s, _ := readit(c)

			_, ok := GetFromMap(s)

			if ok {
				writeIt(c, "User already connected.")
				c.Close()
				return
			}

			fmt.Println(s, " is connected")
			writeIt(c, s+": You are conected\n")
			WriteMap(s, c)
			for {
				// read message
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
				//write message to bufferd channel
				ch <- m
			}
		}(c)

	}
}
