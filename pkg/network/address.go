package network

import (
	"crypto/rand"
	"fmt"
)

func CreateIPAddress() string {
	byts := make([]byte, 2)
	rand.Read(byts)
	return fmt.Sprintf("172.29.%d.%d", byts[0], byts[1])
}
