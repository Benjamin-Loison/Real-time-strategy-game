package main
import (
	//"crypto/rand"
	//"encoding/hex"
	"fmt"
	"time"
)

/*
func Bytes(n int) []byte {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		panic(err)
	}
	return b
}
*/

func random_id(id int) string {
	t := time.Now().UnixMilli()	// get the time to the millisecond
	s := fmt.Sprintf("%d.%d", id, t)
	return s
}

