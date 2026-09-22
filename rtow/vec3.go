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

func NewRandInUnitDisk() Vec3 {
	for {
		p := NewVec3(RandFrom(-1, 1), RandFrom(-1, 1), 0)
		if p.LengthSquared() < 1 {
			return p
		}
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

func (v Vec3) NearZero() bool {
	// Return true if the vector is close to zero in all dimention
	s := 1e-8
	return math.Abs(v.X) < s && (math.Abs(v.Y) < s) && (math.Abs(v.Z) < s)
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

func (u Vec3) Cross(v Vec3) Vec3 {
	return NewVec3(u.Y*v.Z-u.Z*v.Y, u.Z*v.X-u.X*v.Z, u.X*v.Y-u.Y*v.X)
}

func (v Vec3) Unit() Vec3 {
	return v.Div(Full(v.Length()))
}

func (v Vec3) Reflect(n Vec3) Vec3 {
	return v.Sub(Full(2. * v.Dot(n)).Mul(n))
}

func (uv Vec3) Refract(n Vec3, etaiOverEtat float64) Vec3 {
	cosTheta := min(uv.Neg().Dot(n), 1.0)
	rOutPerp := Full(etaiOverEtat).Mul(uv.Add(Full(cosTheta).Mul(n)))
	rOutParallel := Full(-math.Sqrt(math.Abs(1.0 - rOutPerp.LengthSquared()))).Mul(n)
	return rOutPerp.Add(rOutParallel)
}

// Point3 is just alias for Vec3, but useful for geometric clarity in the code.
type Point3 = Vec3

func NewPoint3(x, y, z float64) Point3 {
	return NewVec3(x, y, z)
}
