package main

import (
	"math"
)

/// xi + yj + zk + w
type Quaternion struct {
	x, y, z, w float64
}

func (q Quaternion) conjugate() Quaternion {
	return Quaternion{
		-q.x,
		-q.y,
		-q.z,
		q.w,
	}
}

func (q Quaternion) mul(other Quaternion) Quaternion {
	return Quaternion{
		q.x*other.w + q.y*other.z - q.z*other.y + q.w*other.x,
		-q.x*other.z + q.y*other.w + q.z*other.x + q.w*other.y,
		q.x*other.y - q.y*other.x + q.z*other.w + q.w*other.z,
		-q.x*other.x - q.y*other.y - q.z*other.z + q.w*other.w,
	}
}

type Vector3 struct {
	x, y, z float64
}

func (v Vector3) Length() float64 {
	return math.Sqrt(v.x*v.x + v.y*v.y + v.z*v.z)
}

func (v Vector3) Normalize() Vector3 {
	l := v.Length()
	if l == 0 {
		return Vector3{}
	}
	return Vector3{x: v.x / l, y: v.y / l, z: v.z / l}
}

func (v Vector3) Rotate(asix Vector3, angle float64) Vector3 {
	if asix.Length() == 0 {
		return v // 軸が存在しない場合は回転せずそのまま返す
	}
	radian := float64(angle) * math.Pi / 180

	u := asix.Normalize()
	sin := math.Sin(radian / 2)
	cos := math.Cos(radian / 2)
	q := Quaternion{
		sin * u.x,
		sin * u.y,
		sin * u.z,
		cos,
	}

	vq := Quaternion{
		v.x,
		v.y,
		v.z,
		0,
	}

	ret := q.mul(vq).mul(q.conjugate())
	return Vector3{
		ret.x,
		ret.y,
		ret.z,
	}
}
