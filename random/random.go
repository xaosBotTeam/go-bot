package random

import (
	"math/rand"
	"time"
)

func RandomWaitTime() time.Duration {
	baseWaitTime := 50 * time.Millisecond
	randAdd := time.Duration(rand.Int31n(50)-25) * time.Millisecond
	return baseWaitTime + randAdd
}
