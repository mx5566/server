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

func RandomInt[T RandType](min, max T) T {
	return min + (max-min)*T(rand.Uint64())
}

var Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
