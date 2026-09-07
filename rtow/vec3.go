package main

import (
	"math"
	"strconv"
)

type Vec3 struct {
	X, Y, Z float64
}

func NewVec3(e0, e1, e2 float64) Vec3 {
	return Vec3{
		X: e0,
		Y: e1,
		Z: e2,
	}
}

func NewRandVec3() Vec3 {
	return Vec3{
		X: Rand(),
		Y: Rand(),
		Z: Rand(),
	}
}

func NewRandVec3From(min, max float64) Vec3 {
	return Vec3{
		X: RandFrom(min, max),
		Y: RandFrom(min, max),
		Z: RandFrom(min, max),
	}
}

func NewRandUnitVec3() Vec3 {
	for {
		p := NewRandVec3From(-1, 1)
		lensq := p.LengthSquared()
		if 1e-160 < lensq && lensq <= 1 {
			return p.Div(Full(math.Sqrt(lensq)))
		}
	}
}

func NewRandVec3OnHemisphere(normal Vec3) Vec3 {
	onUnitSphere := NewRandUnitVec3()
	if onUnitSphere.Dot(normal) > 0.0 {
		// In the same hemisphere as the normal.
		return onUnitSphere
	} else {
		return onUnitSphere.Neg()
	}
}

func Full(f float64) Vec3 {
	return NewVec3(f, f, f)
}

func (v Vec3) Neg() Vec3 {
	return NewVec3(-v.X, -v.Y, -v.Z)
}

func (v Vec3) Add(other Vec3) Vec3 {
	v.X += other.X
	v.Y += other.Y
	v.Z += other.Z
	return v
}

func (v Vec3) Sub(other Vec3) Vec3 {
	v.X -= other.X
	v.Y -= other.Y
	v.Z -= other.Z
	return v
}

func (v Vec3) Mul(other Vec3) Vec3 {
	v.X *= other.X
	v.Y *= other.Y
	v.Z *= other.Z
	return v
}

func (v Vec3) Div(other Vec3) Vec3 {
	v.X /= other.X
	v.Y /= other.Y
	v.Z /= other.Z
	return v
}

func (v Vec3) LengthSquared() float64 {
	return v.X*v.X + v.Y*v.Y + v.Z*v.Z
}

func (v Vec3) Length() float64 {
	return math.Sqrt(v.LengthSquared())
}

// Vector Utility Functions
func (v Vec3) String() string {
	prec := 5
	floatToString := func(f float64) string { return strconv.FormatFloat(f, 'f', prec, 64) }
	return floatToString(v.X) + " " + floatToString(v.Y) + " " + floatToString(v.Z)
}

func (v Vec3) Dot(other Vec3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

func (v Vec3) Unit() Vec3 {
	return v.Div(Full(v.Length()))
}

// Point3 is just alias for Vec3, but useful for geometric clarity in the code.
type Point3 = Vec3

func NewPoint3(x, y, z float64) Point3 {
	return NewVec3(x, y, z)
}
