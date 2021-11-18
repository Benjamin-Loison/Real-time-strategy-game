package main
import (
	//"crypto/rand"
	//"encoding/hex"
	"math/rand"
	"time"
)

func RandomString(len int) string {
      bytes := make([]byte, len)
     for i := 0; i < len; i++ {
          bytes[i] = byte(65 + rand.Intn(25))  //A=65 and Z = 65+25
      }
      return string(bytes)
}

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

func random_id(l int) string {
	t := string(time.Now().UnixMilli())// get the time to the millisecond
	s := RandomString(l - len(t))
	return t + s
}

