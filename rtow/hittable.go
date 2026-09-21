package main

type HitRecord struct {
	P         Point3
	Normal    Vec3
	Material  Material
	T         float64
	FrontFace bool
}

type Hittable interface {
	Hit(ray Ray, rayT Interval) (*HitRecord, bool)
}

func NewHitRecord(p Point3, normal Vec3, t float64, frontFace bool, material Material) *HitRecord {
	return &HitRecord{
		P:         p,
		Normal:    normal,
		Material:  material,
		T:         t,
		FrontFace: frontFace,
	}
}

func (r *HitRecord) SetFaceNormal(ray Ray, outwardNormal Vec3) {
	// Set the hit record normal vector.
	// NOTE: the parameter `outward_normal` is assumed to have unit length.

	r.FrontFace = ray.Direction.Dot(outwardNormal) < 0
	if r.FrontFace {
		r.Normal = outwardNormal
	} else {
		r.Normal = outwardNormal.Neg()
	}
}

type HittableList struct {
	Objects []Hittable
}

func NewHittableList() *HittableList {
	return &HittableList{}
}

func (l *HittableList) Clear() {
	l.Objects = []Hittable{}
}

func (l *HittableList) Add(object Hittable) {
	l.Objects = append(l.Objects, object)
}

func (l *HittableList) Hit(ray Ray, rayT Interval) (*HitRecord, bool) {
	hitAnything := false
	closestSoFar := rayT.Max
	var lastRecord *HitRecord = nil

	for _, object := range l.Objects {
		if record, ok := object.Hit(ray, NewInterval(rayT.Min, closestSoFar)); ok {
			hitAnything = true
			closestSoFar = record.T
			lastRecord = record
		}
	}

	return lastRecord, hitAnything
}
