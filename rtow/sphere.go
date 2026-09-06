package main

import (
	"math"
)

type Sphere struct {
	Center Point3
	Radius float64
}

func NewSphere(center Point3, radius float64) Sphere {
	return Sphere {
		Center: center,
		Radius: max(0, radius),
	}
}

func (s Sphere) Hit(ray Ray, rayT Interval) (*HitRecord, bool) {
	oc := s.Center.Sub(ray.Origin)
	a := ray.Direction.LengthSquared()
	h := ray.Direction.Dot(oc)
	c := oc.LengthSquared() - s.Radius*s.Radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return nil, false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range.
	root := (h - sqrtd) / a;
	if !rayT.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rayT.Surrounds(root) {
			return nil, false
		}
	}

	t := root
	p := ray.At(t)
	normal := p.Sub(s.Center).Div(Full(s.Radius))
	record := NewHitRecord(p, normal, t, false)

	outwardNormal := p.Sub(s.Center).Div(Full(s.Radius))
	record.SetFaceNormal(ray, outwardNormal)

	return &record, true
}
