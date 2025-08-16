package utils

import (
	"math/rand"
	"time"
)

func GenerateUniqueID() int64 {
	ts := time.Now().UnixMilli()
	return ts*1000 + rand.Int63n(1000)
}
