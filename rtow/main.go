package main

import (
	"log/slog"
)

func main() {
	logger := slog.Default()

	world := NewHittableList()

	materialGround := NewLambertial(NewColor(0.8, 0.8, 0.0))
	materialCenter := NewLambertial(NewColor(0.1, 0.2, 0.5))
	materialLeft := NewDielectric(1.50)
	materialRight := NewMetal(NewColor(0.8, 0.6, 0.2), 1.0)

	world.Add(NewSphere(NewPoint3(0.0, -100.5, -1.0), 100.0, &materialGround))
	world.Add(NewSphere(NewPoint3(0.0, 0.0, -1.2), 0.5, &materialCenter))
	world.Add(NewSphere(NewPoint3(-1.0, 0.0, -1.0), 0.5, &materialLeft))
	world.Add(NewSphere(NewPoint3(1.0, 0.0, -1.0), 0.5, &materialRight))

	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	samplesPerPixel := 100
	maxDepth := 50

	camera := NewCamera(logger, world, aspectRatio, imageWidth, samplesPerPixel, maxDepth)
	camera.Render(world)
}
