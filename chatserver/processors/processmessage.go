package processors

import (
	"cryptolessons/chatserver/model"
	"encoding/json"
)



func ProcessMessage(m model.CommonMessage , encoder *json.Encoder) {



	if m.KeyExchg {

	}else if m.Msg {

	}else if m.Reg {

	}
}


