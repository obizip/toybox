package main

import (
	"io"
	"math"
	"strconv"
)

type Color = Vec3

func NewColor(r, g, b float64) Color {
	return NewVec3(r, g, b)
}

func NewRandColor() Color {
	return NewColor(Rand(), Rand(), Rand())
}

func NewRandColorFrom(min, max float64) Color {
	return NewColor(RandFrom(min, max), RandFrom(min, max), RandFrom(min, max))
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

	var buffer [32]byte
	line := buffer[:0]
	line = strconv.AppendInt(line, int64(ir), 10)
	line = append(line, ' ')
	line = strconv.AppendInt(line, int64(ig), 10)
	line = append(line, ' ')
	line = strconv.AppendInt(line, int64(ib), 10)
	line = append(line, '\n')
	_, _ = writer.Write(line)
}
