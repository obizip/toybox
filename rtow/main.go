package main

import (
	"log/slog"
)

func main() {
	logger := slog.Default()

	world := NewHittableList()
	world.Add(NewSphere(NewPoint3(0., 0., -1.), 0.5))
	world.Add(NewSphere(NewPoint3(0., -100.5, -1.), 100))

	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	samplesPerPixel := 100
	maxDepth := 50

	camera := NewCamera(logger, world, aspectRatio, imageWidth, samplesPerPixel, maxDepth)
	camera.Render(world)
}
