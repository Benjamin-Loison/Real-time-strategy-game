package main
import (
	"crypto/rand"
	"encoding/hex"
)

func Bytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}

func random_id(n int) string {
	return hex.EncodeToString(Bytes(n))
}

