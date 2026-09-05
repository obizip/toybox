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
	r := c.X()
	g := c.Y()
	b := c.Z()

	ir := int(255.999 * r);
	ig := int(255.999 * g);
	ib := int(255.999 * b);

	fmt.Fprintf(writer, "%d %d %d\n", ir, ig, ib)
}
