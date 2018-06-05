package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Port string
	Key  string
	Cert string
}

func GetConfig(file string) Config {
	f, e := os.Open(file)
	if e != nil {
		fmt.Print("Unable to open config file", file, e.Error())
		return Config{}
	}
	d := json.NewDecoder(f)
	c := Config{}
	e = d.Decode(&c)
	if e != nil {
		fmt.Print("Unable to decode config file", file, e.Error())
		return Config{}
	}
	return c
}
