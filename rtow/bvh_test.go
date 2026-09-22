package main

import "testing"

func TestAABBHitWithParallelRay(t *testing.T) {
	box := NewAABB(NewPoint3(-1, -1, -1), NewPoint3(1, 1, 1))

	if !box.Hit(NewRay(NewPoint3(0, 0, -2), NewVec3(0, 0, 1)), NewInterval(0, 100)) {
		t.Fatal("expected ray inside parallel slabs to hit")
	}
	if box.Hit(NewRay(NewPoint3(2, 0, -2), NewVec3(0, 0, 1)), NewInterval(0, 100)) {
		t.Fatal("expected ray outside parallel slab to miss")
	}
}

func TestBVHMatchesHittableList(t *testing.T) {
	material := NewLambertian(NewColor(0.8, 0.2, 0.1))
	list := NewHittableList()
	list.Add(NewSphere(NewPoint3(-1, 0, -3), 0.5, material))
	list.Add(NewSphere(NewPoint3(0, 0, -4), 1.0, material))
	list.Add(NewSphere(NewPoint3(1, 0, -3), 0.5, material))
	bvh := NewBVHNode(list.Objects)

	rays := []Ray{
		NewRay(NewPoint3(0, 0, 0), NewVec3(0, 0, -1)),
		NewRay(NewPoint3(0, 0, 0), NewVec3(-0.3, 0, -1)),
		NewRay(NewPoint3(0, 0, 0), NewVec3(0.3, 0, -1)),
		NewRay(NewPoint3(0, 0, 0), NewVec3(0, 1, -1)),
	}

	for _, ray := range rays {
		listRecord, listHit := list.Hit(ray, NewInterval(0.001, 100))
		bvhRecord, bvhHit := bvh.Hit(ray, NewInterval(0.001, 100))
		if listHit != bvhHit {
			t.Fatalf("hit mismatch for ray %#v: list=%v bvh=%v", ray, listHit, bvhHit)
		}
		if listHit && listRecord.T != bvhRecord.T {
			t.Fatalf("nearest hit mismatch for ray %#v: list=%v bvh=%v", ray, listRecord.T, bvhRecord.T)
		}
	}
}
