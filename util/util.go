package util

import (
	"math/rand"
	"time"
)

func GetRandomCellState() bool {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(100)%3 == 0
}
