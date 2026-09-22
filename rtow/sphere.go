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

func (s Sphere) Hit(ray Ray, rayT Interval) (HitRecord, bool) {
	oc := s.Center.Sub(ray.Origin)
	a := ray.Direction.LengthSquared()
	h := ray.Direction.Dot(oc)
	c := oc.LengthSquared() - s.Radius*s.Radius

	discriminant := h*h - a*c
	if discriminant < 0 {
		return HitRecord{}, false
	}

	sqrtd := math.Sqrt(discriminant)

	// Find the nearest root that lies in the acceptable range.
	root := (h - sqrtd) / a
	if !rayT.Surrounds(root) {
		root = (h + sqrtd) / a
		if !rayT.Surrounds(root) {
			return HitRecord{}, false
		}
	}

	t := root
	p := ray.At(t)
	outwardNormal := p.Sub(s.Center).Div(Full(s.Radius))
	record := HitRecord{P: p, T: t, Material: s.Material}
	record.SetFaceNormal(ray, outwardNormal)

	return record, true
}

func (s Sphere) BoundingBox() AABB {
	radius := Full(s.Radius)
	return NewAABB(s.Center.Sub(radius), s.Center.Add(radius))
}
