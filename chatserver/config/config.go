package config

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type Config struct {
	Port        string
	Key         string
	Cert        string
	Loglevel    string
	ChannelSize int
	Workers     int
}

var CFG Config

func GetConfig(file string) (Config, error) {
	f, e := os.Open(file)
	if e != nil {
		log.Println("Unable to open config file", file, e.Error())
		return CFG, errors.New("Unable to open file.")
	}
	d := json.NewDecoder(f)
	e = d.Decode(&CFG)
	if e != nil {
		log.Println("Unable to decode config file", file, e.Error())
		return CFG, errors.New("Unable to open file.")
	}
	return CFG, nil
}

func init() {
	GetConfig(os.Getenv("APP_CFG"))
}
