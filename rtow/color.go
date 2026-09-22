package main

import (
	"fmt"
	"io"
	"math"
)

type Color = Vec3

func NewColor(r, g, b float64) Color {
	return NewVec3(r, g, b)
}

func linearToGamma(linearComponent float64) float64 {
	if linearComponent > 0 {
		return math.Sqrt(linearComponent)
	}

	return 0
}

func (c Color) Write(writer io.Writer) {
	r := linearToGamma(c.X)
	g := linearToGamma(c.Y)
	b := linearToGamma(c.Z)

	intensity := NewInterval(0.0, 0.999)
	ir := int(256 * intensity.Clamp(r))
	ig := int(256 * intensity.Clamp(g))
	ib := int(256 * intensity.Clamp(b))

	fmt.Fprintf(writer, "%d %d %d\n", ir, ig, ib)
}
