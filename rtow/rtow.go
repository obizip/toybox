package main

import (
	"math"
	"math/rand"
)

func DegreeToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180.0
}

func Rand() float64 {
	return rand.Float64()
}

func RandFrom(min, max float64) float64 {
	return rand.Float64()*(max-min) + min
}
