package main

import (
	"math"
)

type Interval struct {
	Min, Max float64
}

func NewInterval(min, max float64) Interval {
	return Interval{
		Min: min,
		Max: max,
	}
}

func NewEmptyInterval(min, max float64) Interval {
	return Interval{
		Min: math.Inf(1), 
		Max: math.Inf(-1),
	}
}

func NewUniverseInterval(min, max float64) Interval {
	return Interval{
		Min: math.Inf(-1), 
		Max: math.Inf(1),
	}
}

func (i *Interval) Size() float64 {
	return i.Max - i.Min
}

func (i *Interval) Contains(x float64) bool {
	return i.Min <= x && x <= i.Max
}

func (i *Interval) Surrounds(x float64) bool {
	return i.Min < x && x < i.Max
}
