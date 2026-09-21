package main

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
}

func NewMetal(albedo Color) Metal {
	return Metal{
		Albedo: albedo,
	}
}

func (m *Metal) Scatter(rayIn Ray, record *HitRecord) (attenuation Color, scattered Ray, ok bool) {
	reflected := rayIn.Direction.Reflect(record.Normal)
	scattered = NewRay(record.P, reflected)
	attenuation = m.Albedo
	ok = true
	return
}


