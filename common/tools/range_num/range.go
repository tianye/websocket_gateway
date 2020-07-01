package range_num

import (
	"math/rand"
	"time"
)

func GenerateRangeNum(min int, max int) int {
	if min == max {
		return min
	}

	rand.Seed(time.Now().Unix())

	randNum := rand.Intn(max-min) + min

	return randNum
}
