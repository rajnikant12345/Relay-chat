package applog

import (
	"log"
	"io/ioutil"
	"os"
	"cryptolessons/chatserver/config"
)

var (
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func init() {

	level := config.CFG.Loglevel

	flag := log.Ldate|log.Ltime|log.Llongfile

	Warning = log.New(os.Stdout,
		"WARNING: ",
		flag)

	Error = log.New(os.Stderr,
		"ERROR: ",
		flag)

	if level == "DEBUG" {
		Trace = log.New(os.Stdout,
			"TRACE: ",
			flag)

		Info = log.New(os.Stdout,
			"INFO: ",
			flag)
	} else {
		Trace = log.New(ioutil.Discard,
			"TRACE: ",
			flag)

		Info = log.New(ioutil.Discard,
			"INFO: ",
			flag)

	}

}


