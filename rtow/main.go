package main

import (
	"log/slog"
)

func main() {
	logger := slog.Default()

	world := NewHittableList()

	groundMaterial := NewLambertian(NewColor(0.5, 0.5, 0.5))
	world.Add(NewSphere(NewPoint3(0, -1000, 0), 1000, groundMaterial))

	for a := -11; a < 11; a++ {
		for b := -11; b < 11; b++ {
			chooseMat := Rand()
			center := NewVec3(float64(a)+0.9*Rand(), 0.2, float64(b)+0.9*Rand())

			if center.Sub(NewPoint3(4, 0.2, 0)).Length() > 0.9 {
				var sphereMaterial Material

				if chooseMat < 0.8 {
					// diffuse
					albedo := NewRandColor().Mul(NewRandColor())
					sphereMaterial = NewLambertian(albedo)
					world.Add(NewSphere(center, 0.2, sphereMaterial))
				} else if chooseMat < 0.95 {
					// metal
					albedo := NewRandColorFrom(0.5, 1)
					fuzz := RandFrom(0, 0.5)
					sphereMaterial = NewMetal(albedo, fuzz)
					world.Add(NewSphere(center, 0.2, sphereMaterial))
				} else {
					// glass
					sphereMaterial = NewDielectric(1.5)
					world.Add(NewSphere(center, 0.2, sphereMaterial))
				}
			}
		}
	}

	material1 := NewDielectric(1.5)
	world.Add(NewSphere(NewPoint3(0., 1., 0.), 1.0, material1))

	material2 := NewLambertian(NewColor(0.4, 0.2, 0.1))
	world.Add(NewSphere(NewPoint3(-4., 1., 0.), 1.0, material2))

	material3 := NewMetal(NewColor(0.7, 0.6, 0.5), 0.0)
	world.Add(NewSphere(NewPoint3(4., 1., 0.), 1.0, material3))

	aspectRatio := 16.0 / 9.0
	imageWidth := 1200
	samplesPerPixel := 500
	maxDepth := 50
	verticalFov := 20.0
	lookFrom := NewPoint3(13., 2., 3.)
	lookAt := NewPoint3(0., 0., 0.)
	verticalUp := NewVec3(0., 1., 0.)

	defocusAngle := 0.6
	focusDistance := 10.0

	worldBVH := NewBVHNode(world.Objects)
	camera := NewCamera(
		logger, aspectRatio, imageWidth, samplesPerPixel, maxDepth, verticalFov, lookFrom, lookAt, verticalUp, defocusAngle, focusDistance)

	camera.Render(worldBVH)
}
