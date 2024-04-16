package base

import (
	"math/rand"
	"time"
)

func Random[T ~float32](min, max T) T {
	return min + (max-min)*T(rand.Float64())
}

type RandType interface {
	~int | ~int32 | ~int64 | ~uint | ~uint32 | ~uint64
}

// [min, max)
func RandomInt(min, max int) int {
	if min > max {
		min, max = max, min
	} else if min == max {
		return min
	}

	return min + rand.Int()%(max-min)
}

func RandomInt32(min, max int32) int32 {
	if min > max {
		min, max = max, min
	} else if min == max {
		return min
	}

	return min + rand.Int31()%(max-min)
}

// [min,max]
func RandomIntA(min, max int) int {
	if min >= max {
		min, max = max, min
	}

	return min + rand.Int()%(max-min+1)
}

var Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
