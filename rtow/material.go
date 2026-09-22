package main

import (
	"math"
)

type Material interface {
	Scatter(rayIn Ray, record *HitRecord) (attenuation Color, scattered Ray, ok bool)
}

type Lambertian struct {
	Albedo Color
}

func NewLambertial(albedo Color) Lambertian {
	return Lambertian{
		Albedo: albedo,
	}
}

func (l *Lambertian) Scatter(rayIn Ray, record *HitRecord) (attenuation Color, scattered Ray, ok bool) {
	scatterDirection := record.Normal.Add(NewRandUnitVec3())

	// Catch degenerate scatter direction
	if scatterDirection.NearZero() {
		scatterDirection = record.Normal
	}

	scattered = NewRay(record.P, scatterDirection)
	attenuation = l.Albedo
	ok = true
	return
}

type Metal struct {
	Albedo Color
	Fuzz  float64
}

func NewMetal(albedo Color, fuzz float64) Metal {
	return Metal{
		Albedo: albedo,
		Fuzz: fuzz,
	}
}

func (m *Metal) Scatter(rayIn Ray, record *HitRecord) (attenuation Color, scattered Ray, ok bool) {
	reflected := rayIn.Direction.Reflect(record.Normal)
	reflected = reflected.Unit().Add(Full(m.Fuzz).Mul(NewRandUnitVec3()))
	scattered = NewRay(record.P, reflected)
	attenuation = m.Albedo
	ok = scattered.Direction.Dot(record.Normal) > 0
	return
}

type Dielectric struct {
	// Refractive index in vacuum or air, or the ratio of the material's refractive index
	// over the refractive index of the enclosing media.
	RefractionIndex float64
}

func NewDielectric(refractionIndex float64) Dielectric {
	return Dielectric {
		RefractionIndex: refractionIndex,
	}
}

func (d *Dielectric) Scatter(rayIn Ray, record *HitRecord) (attenuation Color, scattered Ray, ok bool) {
	attenuation = NewColor(1.0, 1.0, 1.0)
	var ri float64
	if record.FrontFace {
		ri = 1.0 / d.RefractionIndex
	} else {
		ri = d.RefractionIndex
	}
	unitDirection := rayIn.Direction.Unit()
	cosTheta := min(unitDirection.Neg().Dot(record.Normal), 1.0)
	sinTheta := math.Sqrt(1.0 - cosTheta*cosTheta)

	cannotRefract := ri * sinTheta > 1.0
	var direction Vec3
	if cannotRefract {
		direction = unitDirection.Reflect(record.Normal)
	} else {
		direction = unitDirection.Refract(record.Normal, ri)
	}

	scattered = NewRay(record.P, direction)
	ok = true
	return
}
