package main

import (
	"fmt"
	"log/slog"
	"math"
	"os"
)

func hitSphere(center Point3, radius float64, ray Ray) float64 {
	oc := center.Sub(ray.Origin)
	a := ray.Direction.Dot(ray.Direction)
	b := -2.0 * ray.Direction.Dot(oc)
	c := oc.Dot(oc) - radius*radius
	// Quadratic Formula
	// \frac{-b \pm \sqrt{b^2 - 4ac}}{2a}
	// We can determine the number of solutions by calculating the discriminant:
	// positive means two real solutions,
	// negative means no real solutions,
	// and zero means one real solution.
	discriminant := b*b - 4.0*a*c
	if discriminant < 0 {
		return -1.0
	} else {
		return (-b - math.Sqrt(discriminant)) / (2.0 * a)
	}
}

func rayColor(ray Ray) Color {
	t := hitSphere(NewPoint3(0., 0., -1), 0.5, ray)
	if t > 0.0 {
		N := ray.At(t).Sub(NewVec3(0., 0., -1)).Unit()
		return Full(0.5).Mul(NewColor(N.X+1., N.Y+1., N.Z+1))
	}
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
	imageHeight := max(int(float64(imageWidth)/aspectRatio), 1)

	//
	// Camera
	//
	focal_length := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * float64(imageWidth) / float64(imageHeight)
	cameraCenter := NewPoint3(0.0, 0.0, 0.0)

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := NewVec3(viewportWidth, 0, 0)
	viewportV := NewVec3(0, -viewportHeight, 0)

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	pixelDeltaU := viewportU.Div(Full(float64(imageWidth)))
	pixelDeltaV := viewportV.Div(Full(float64(imageHeight)))

	// Calculate the location of the upper left pixel.
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
			ray := NewRay(cameraCenter, rayDirection)
			pixelColor := rayColor(ray)
			pixelColor.Write(os.Stdout)
		}
	}
}
