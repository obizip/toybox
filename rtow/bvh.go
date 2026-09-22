package main

import "sort"

// BVHNode hierarchically groups objects by bounding box. A ray only visits
// branches whose boxes it intersects.
type BVHNode struct {
	Left, Right Hittable
	Box         AABB
}

func NewBVHNode(objects []Hittable) *BVHNode {
	if len(objects) == 0 {
		return nil
	}

	items := append([]Hittable(nil), objects...)
	return newBVHNode(items)
}

func newBVHNode(objects []Hittable) *BVHNode {
	if len(objects) == 1 {
		return &BVHNode{Left: objects[0], Box: objects[0].BoundingBox()}
	}

	box := objects[0].BoundingBox()
	for _, object := range objects[1:] {
		box = SurroundingAABB(box, object.BoundingBox())
	}

	axis := box.LongestAxis()
	sort.Slice(objects, func(i, j int) bool {
		return aabbAxisMin(objects[i].BoundingBox(), axis) < aabbAxisMin(objects[j].BoundingBox(), axis)
	})
	middle := len(objects) / 2
	left := newBVHNode(objects[:middle])
	right := newBVHNode(objects[middle:])

	return &BVHNode{
		Left:  left,
		Right: right,
		Box:   SurroundingAABB(left.Box, right.Box),
	}
}

func (n *BVHNode) Hit(ray Ray, rayT Interval) (HitRecord, bool) {
	if n == nil || !n.Box.Hit(ray, rayT) {
		return HitRecord{}, false
	}

	leftRecord, hitLeft := HitRecord{}, false
	if n.Left != nil {
		leftRecord, hitLeft = n.Left.Hit(ray, rayT)
	}
	if hitLeft {
		rayT.Max = leftRecord.T
	}
	if n.Right != nil {
		if rightRecord, hitRight := n.Right.Hit(ray, rayT); hitRight {
			return rightRecord, true
		}
	}
	return leftRecord, hitLeft
}

func (n *BVHNode) BoundingBox() AABB {
	return n.Box
}

func aabbAxisMin(box AABB, axis int) float64 {
	switch axis {
	case 0:
		return box.X.Min
	case 1:
		return box.Y.Min
	default:
		return box.Z.Min
	}
}
