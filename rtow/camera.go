package main

import (
	"fmt"
	"log/slog"
	"math"
	"os"
)

type Camera struct {
	Logger *slog.Logger
	// Ratio of image width over height
	AspectRatio float64
	// Rendererd image width in pixel count
	ImageWidth  int
	ImageHeight int
	Center      Point3
	Pixel100Loc Point3
	PixelDeltaU Vec3
	PixelDeltaV Vec3
}

func NewCamera(logger *slog.Logger, aspectRatio float64, imageWidth int, world Hittable) Camera {
	//
	// Image
	//
	imageHeight := max(int(float64(imageWidth)/aspectRatio), 1)

	//
	// Camera
	//
	focal_length := 1.0
	viewportHeight := 2.0
	viewportWidth := viewportHeight * float64(imageWidth) / float64(imageHeight)
	center := NewPoint3(0.0, 0.0, 0.0)

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := NewVec3(viewportWidth, 0, 0)
	viewportV := NewVec3(0, -viewportHeight, 0)

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	pixelDeltaU := viewportU.Div(Full(float64(imageWidth)))
	pixelDeltaV := viewportV.Div(Full(float64(imageHeight)))

	// Calculate the location of the upper left pixel.
	viewportUpperLeft := center.Sub(NewVec3(0., 0., focal_length)).Sub(viewportU.Div(Full(2))).Sub(viewportV.Div(Full(2)))
	pixel100Loc := viewportUpperLeft.Add(Full(0.5).Mul(pixelDeltaU.Add(pixelDeltaV)))

	return Camera {
		Logger      : logger,
		AspectRatio : aspectRatio,
		ImageWidth  : imageWidth,
		ImageHeight : imageHeight,
		Center      : center,
		Pixel100Loc : pixel100Loc,
		PixelDeltaU : pixelDeltaU,
		PixelDeltaV : pixelDeltaV,
	}
}

func (c Camera) Render(world Hittable) {
	fmt.Print("P3\n", c.ImageWidth, c.ImageHeight, "\n255\n")

	for j := range c.ImageHeight {
		c.Logger.Info("Scaning lines", "remainingLines", c.ImageHeight-j)
		for i := range c.ImageWidth {
			pixelCenter := c.Pixel100Loc.Add(Full(float64(i)).Mul(c.PixelDeltaU)).Add(Full(float64(j)).Mul(c.PixelDeltaV))
			rayDirection := pixelCenter.Sub(c.Center)
			ray := NewRay(c.Center, rayDirection)
			pixelColor := c.rayColor(ray, world)
			pixelColor.Write(os.Stdout)
		}
	}
	c.Logger.Info("Finished")
}

func (c Camera) rayColor(ray Ray, world Hittable) Color {
	if record, ok := world.Hit(ray, NewInterval(0, math.Inf(1))); ok {
		return Full(0.5).Mul(record.Normal.Add(NewColor(1., 1., 1.)))
	}

	unitDirection := ray.Direction.Unit()
	a := 0.5 * (unitDirection.Y + 1.0)
	return Full(1.0 - a).Mul(NewColor(1.0, 1.0, 1.0)).Add(Full(a).Mul(NewColor(0.5, 0.7, 1.0)))
}
