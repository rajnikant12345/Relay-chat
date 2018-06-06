package main

import (
	"crypto/rand"
	"encoding/hex"

	"os/exec"
	"sync"
	"fmt"

)

func main() {

	var wg sync.WaitGroup
	bname := "C:/Users/rkant/workspace/src/main.exe"
	fname := "C:/Users/rkant/workspace/src/input.txt"

	for i:=0;i<100;i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			b := make([]byte,64)
			rand.Read(b)
			name := hex.EncodeToString(b)

			cmd := exec.Command(bname, "localhost","6789",name ,fname )

			err := cmd.Run()
			if err != nil {
				fmt.Println(err.Error())
			}
		}()

	}

	wg.Wait()

}
