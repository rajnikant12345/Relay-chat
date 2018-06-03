package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
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

func MessageEncryptor(k, s string) string {
	h := sha256.New()
	b := h.Sum([]byte(k))
	fmt.Println()
	c, e := aes.NewCipher(b[:32])
	if e != nil {
		fmt.Println(e.Error())
		return s
	}

	m := cipher.NewCFBEncrypter(c, b[:c.BlockSize()])
	out := make([]byte, len(s))
	m.XORKeyStream(out, []byte(s))

	return hex.EncodeToString(out)
}

func MessageDecryptor(k, s string) string {
	d, _ := hex.DecodeString(s)
	h := sha256.New()
	b := h.Sum([]byte(k))
	fmt.Println()
	c, e := aes.NewCipher(b[:32])
	if e != nil {
		fmt.Println(e.Error())
		return s
	}
	m := cipher.NewCFBDecrypter(c, b[:c.BlockSize()])
	out := make([]byte, len(d))
	m.XORKeyStream(out, d)
	return string(out)

}

func main() {
	var ip, port, user, friend, key string
	fmt.Print("Enter IP of Server:")
	fmt.Scanf("%s\n", &ip)
	fmt.Print("Enter port of Server:")
	fmt.Scanf("%s\n", &port)
	fmt.Print("Enter encryption key for messages.:")
	fmt.Scanf("%s\n", &key)

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
			if key != "" {
				m.Message = MessageDecryptor(key, m.Message)
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
		if key != "" {
			m.Message = MessageEncryptor(key, m.Message)
		}
		j := json.NewEncoder(c)
		j.Encode(&m)
	}

}
