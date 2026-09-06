package main

import (
	"io"
	"fmt"
)

type Color = Vec3

func NewColor(r, g, b float64) Color {
	return NewVec3(r, g, b)
}

func (c Color) Write(writer io.Writer) {
	r := c.X
	g := c.Y
	b := c.Z

	intensity := NewInterval(0.0, 0.999)
	ir := int(256 * intensity.Clamp(r));
	ig := int(256 * intensity.Clamp(g));
	ib := int(256 * intensity.Clamp(b));

	fmt.Fprintf(writer, "%d %d %d\n", ir, ig, ib)
}
