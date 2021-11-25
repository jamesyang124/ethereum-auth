package utils

import (
	"math/rand"
	"time"
)

func RandNonce() int {
	return 100000 + rand.Intn(int(time.Now().UnixNano()%1000000))
}
