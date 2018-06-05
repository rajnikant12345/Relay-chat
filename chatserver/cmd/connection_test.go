package cmd

import (
	"cryptolessons/chatserver/model"
	"os"
	"testing"
)

func startserver() {
	os.Setenv("APP_CFG", "/Users/rajnikant/workspace/src/cryptolessons/app.json")
	go StartServer()
}

func TestHandleConnections(t *testing.T) {
	startserver()
	m := model.CommonMessage{}
	m.Ref = "asdf1234"

}
