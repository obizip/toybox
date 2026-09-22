package main

import "math"

// AABB is an axis-aligned bounding box used to skip objects a ray cannot hit.
type AABB struct {
	X, Y, Z Interval
}

func NewAABB(a, b Point3) AABB {
	return AABB{
		X: NewInterval(min(a.X, b.X), max(a.X, b.X)),
		Y: NewInterval(min(a.Y, b.Y), max(a.Y, b.Y)),
		Z: NewInterval(min(a.Z, b.Z), max(a.Z, b.Z)),
	}
}

func EmptyAABB() AABB {
	return AABB{
		X: NewInterval(math.Inf(1), math.Inf(-1)),
		Y: NewInterval(math.Inf(1), math.Inf(-1)),
		Z: NewInterval(math.Inf(1), math.Inf(-1)),
	}
}

func SurroundingAABB(a, b AABB) AABB {
	return AABB{
		X: NewInterval(min(a.X.Min, b.X.Min), max(a.X.Max, b.X.Max)),
		Y: NewInterval(min(a.Y.Min, b.Y.Min), max(a.Y.Max, b.Y.Max)),
		Z: NewInterval(min(a.Z.Min, b.Z.Min), max(a.Z.Max, b.Z.Max)),
	}
}

func (b AABB) Hit(ray Ray, rayT Interval) bool {
	origins := [3]float64{ray.Origin.X, ray.Origin.Y, ray.Origin.Z}
	directions := [3]float64{ray.Direction.X, ray.Direction.Y, ray.Direction.Z}
	intervals := [3]Interval{b.X, b.Y, b.Z}

	for axis := range intervals {
		if directions[axis] == 0 {
			if origins[axis] < intervals[axis].Min || origins[axis] > intervals[axis].Max {
				return false
			}
			continue
		}
		invDirection := 1.0 / directions[axis]
		t0 := (intervals[axis].Min - origins[axis]) * invDirection
		t1 := (intervals[axis].Max - origins[axis]) * invDirection
		if invDirection < 0 {
			t0, t1 = t1, t0
		}
		rayT.Min = max(rayT.Min, t0)
		rayT.Max = min(rayT.Max, t1)
		if rayT.Max <= rayT.Min {
			return false
		}
	}
	return true
}

func (b AABB) LongestAxis() int {
	xSize := b.X.Size()
	ySize := b.Y.Size()
	zSize := b.Z.Size()
	if xSize > ySize && xSize > zSize {
		return 0
	}
	if ySize > zSize {
		return 1
	}
	return 2
}
