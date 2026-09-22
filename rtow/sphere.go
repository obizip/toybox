package main

import (
	"math"
)

type Sphere struct {
	Center   Point3
	Radius   float64
	Material Material
}

func NewSphere(center Point3, radius float64, material Material) Sphere {
	// TODO: Initialize the material pointer
	return Sphere{
		Center:   center,
		Radius:   max(0, radius),
		Material: material,
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
	root := (h - sqrtd) / a
	if !rayT.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rayT.Surrounds(root) {
			return nil, false
		}
	}

	t := root
	p := ray.At(t)
	normal := p.Sub(s.Center).Div(Full(s.Radius))
	record := NewHitRecord(p, normal, t, false, s.Material)

	outwardNormal := p.Sub(s.Center).Div(Full(s.Radius))
	record.SetFaceNormal(ray, outwardNormal)

	return record, true
}
