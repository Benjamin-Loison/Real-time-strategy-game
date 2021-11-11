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

func location_type_from_str(s string) location_type {
	switch (s) {
		case "Floor":
			return Floor
		case "ElfBuilding":
			return ElfBuilding
		case "OrcBuilding":
			return OrcBuilding
		case "HumanBuilding":
			return HumanBuilding
		default:
			return Floor
	}
}

func random_id(id int) string {
	t := time.Now().UnixMilli()	// get the time to the millisecond
	s := fmt.Sprintf("%d%d", id, t)
	return s
}

