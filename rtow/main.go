package main

import (
	"fmt"
	"log/slog"
	"os"
)

func rayColor(ray Ray) Color {
	unitDirection := ray.Direction.Unit()
	a := 0.5 * (unitDirection.Y + 1.0)
	return Full(1.0 - a).Mul(NewColor(1.0, 1.0, 1.0)).Add(Full(a).Mul(NewColor(0.5, 0.7, 1.0)))
}

func main() {
	logger := slog.Default()

	//
	// Image
	//
	aspectRatio := 16.0 / 9.0
	imageWidth := 400
	imageHeight := max(int(float64(imageWidth) / aspectRatio), 1)

	//
	// Camera
	//
	focal_length := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * float64(imageWidth)/float64(imageHeight)
	cameraCenter := NewPoint3(0.0, 0.0, 0.0)

	// Caluculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := NewVec3(viewportWidth, 0, 0)
	viewportV := NewVec3(0, -viewportHeight, 0)

	// Caluculate the horizontal and vertical delta vectors from pixel to pixel.
	pixelDeltaU := viewportU.Div(Full(float64(imageWidth)))
	pixelDeltaV := viewportV.Div(Full(float64(imageHeight)))

	// Caluculate the location of the upper left pixel.
	viewportUpperLeft := cameraCenter.Sub(NewVec3(0., 0., focal_length)).Sub(viewportU.Div(Full(2))).Sub(viewportV.Div(Full(2)))
	pixel100Loc := viewportUpperLeft.Add(Full(0.5).Mul(pixelDeltaU.Add(pixelDeltaV)))

	//
	// Render
	//
	fmt.Print("P3\n", imageWidth, imageHeight, "\n255\n")

	for j := range imageHeight {
		logger.Info("Scaning lines", "remainingLines", imageHeight-j)
		for i := range imageWidth {
			pixelCenter := pixel100Loc.Add(Full(float64(i)).Mul(pixelDeltaU)).Add(Full(float64(j)).Mul(pixelDeltaV))
			rayDirection := pixelCenter.Sub(cameraCenter)
			ray := NewRay(pixelCenter, rayDirection)
			pixelColor := rayColor(ray)
			pixelColor.Write(os.Stdout)
		}
	}
}
