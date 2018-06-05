package model

type CommonMessage struct {
	Conn string `json:"conn,omitempty"`
	Ref string `json:"ref,omitempty"`
	KeyExchg *KeyExchange	`json:"key_exchg,omitempty"`
	Msg     *Message		`json:"message,omitempty"`
	Lgin    *Login		`json:"login,omitempty"`
	Reg     *Register	`json:"reg,omitempty"`
}


type Register struct {
	UserName string		`json:"user_name,omitempty"`
	PasswordHash string	`json:"password_hash,omitempty"`
	ConfirmHash string	`json:"confirm_hash,omitempty"`
}

type Login struct {
	UserName string		`json:"user_name,omitempty"`
	PasswordHash string	`json:"password_hash,omitempty"`
}

type Message struct {
	Jwt string		`json:"jwt,omitempty"`
	To string		`json:"to,omitempty"`
	From string		`json:"from,omitempty"`
	DataType string		`json:"data_type,omitempty"`
	Data string		`json:"data,omitempty"`
}

type KeyExchange struct {
	Jwt string	`json:"jwt,omitempty"`
	Key string	`json:"key,omitempty"`
}

func GetUserFromConnId(id string) {
	
}