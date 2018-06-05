package model

type CommonMessage struct {
	Req Request
	Res Response
}

type Request struct {
	RequestType string `json:"request_type"`
	Reg     Register	`json:"reg"`
	Lgin    Login		`json:"login"`
	Msg     Message		`json:"message"`
	KeyExchg KeyExchange	`json:"key_exchg"`
}

type Register struct {
	UserName string		`json:"user_name"`
	PasswordHash string	`json:"password_hash"`
	ConfirmHash string	`json:"confirm_hash"`
}

type Login struct {
	UserName string		`json:"user_name"`
	PasswordHash string	`json:"password_hash"`
}

type Message struct {
	Jwt string		`json:"jwt"`
	To string		`json:"to"`
	From string		`json:"from"`
	DataType string		`json:"data_type"`
	Data string		`json:"data"`
}

type KeyExchange struct {
	Jwt string	`json:"jwt"`
	Key string	`json:"key"`
}

type Response struct {
	ResponseType string	`json:"response_type"`
	IncomingMessage Message	`json:"incoming_message"`
	KeyExchange	KeyExchange `json:"key_exchange"`
}
