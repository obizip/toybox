package main

type HitRecord struct {
	P         Point3
	Normal    Vec3
	Material  Material
	T         float64
	FrontFace bool
}

type Hittable interface {
	// Hit returns the nearest acceptable intersection by value.
	Hit(ray Ray, rayT Interval) (HitRecord, bool)
	BoundingBox() AABB
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

func (l *HittableList) Hit(ray Ray, rayT Interval) (HitRecord, bool) {
	hitAnything := false
	closestSoFar := rayT.Max
	var result HitRecord

	for _, object := range l.Objects {
		if candidate, ok := object.Hit(ray, NewInterval(rayT.Min, closestSoFar)); ok {
			hitAnything = true
			closestSoFar = candidate.T
			result = candidate
		}
	}

	return result, hitAnything
}

func (l *HittableList) BoundingBox() AABB {
	if len(l.Objects) == 0 {
		return EmptyAABB()
	}

	box := l.Objects[0].BoundingBox()
	for _, object := range l.Objects[1:] {
		box = SurroundingAABB(box, object.BoundingBox())
	}
	return box
}
