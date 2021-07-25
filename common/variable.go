package common

import (
	"math/rand"
	"time"
)

var (
	SeededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)
