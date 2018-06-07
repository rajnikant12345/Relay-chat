package processors

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/google/uuid"
)

func GenerateReference() string {
	c := 20
	b := make([]byte, c)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func GenerateConnectinId() string {
	p, _ := uuid.NewRandom()
	return p.String()
}
