package utils

import (
	"math/rand"
	"time"
)

func RandSelect(values []string) string {
	index := rand.Intn(len(values))
	return values[index]
}

func RandRangeFloat64(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}

func RandEvent(probability float32) bool {
	if probability >= 1.0 {
		return true
	}
	if probability <= 0.0 {
		return false
	}
	return rand.Float32() <= probability
}

func InitRandSeed() {
	rand.Seed(time.Now().UTC().UnixNano())
}
