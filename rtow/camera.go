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
	ImageWidth        int
	ImageHeight       int
	SamplePerPixel    int
	PixelSamplesScale float64
	Center            Point3
	Pixel00Loc        Point3
	PixelDeltaU       Vec3
	PixelDeltaV       Vec3
	MaxDepth          int
}

func NewCamera(logger *slog.Logger, world Hittable, aspectRatio float64, imageWidth int, samplesPerPixel int, maxDepth int) Camera {
	//
	// Image
	//
	imageHeight := max(int(float64(imageWidth)/aspectRatio), 1)
	pixelSamplesScale := 1.0 / float64(samplesPerPixel)

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

	return Camera{
		Logger:            logger,
		AspectRatio:       aspectRatio,
		ImageWidth:        imageWidth,
		ImageHeight:       imageHeight,
		SamplePerPixel:    samplesPerPixel,
		PixelSamplesScale: pixelSamplesScale,
		Center:            center,
		Pixel00Loc:        pixel100Loc,
		PixelDeltaU:       pixelDeltaU,
		PixelDeltaV:       pixelDeltaV,
		MaxDepth:          maxDepth,
	}
}

func (c Camera) Render(world Hittable) {
	fmt.Print("P3\n", c.ImageWidth, c.ImageHeight, "\n255\n")

	for j := range c.ImageHeight {
		c.Logger.Info("Scaning lines", "remainingLines", c.ImageHeight-j)
		for i := range c.ImageWidth {
			pixelColor := NewColor(0., 0., 0.)
			for range c.SamplePerPixel {
				ray := c.getRay(i, j)
				pixelColor = pixelColor.Add(c.rayColor(ray, world, c.MaxDepth))
			}
			// pixelCenter := c.Pixel100Loc.Add(Full(float64(i)).Mul(c.PixelDeltaU)).Add(Full(float64(j)).Mul(c.PixelDeltaV))
			// rayDirection := pixelCenter.Sub(c.Center)
			// ray := NewRay(c.Center, rayDirection)
			// pixelColor := c.rayColor(ray, world)
			pixelColor = Full(c.PixelSamplesScale).Mul(pixelColor)
			pixelColor.Write(os.Stdout)
		}
	}
	c.Logger.Info("Finished")
}

func (c Camera) getRay(i, j int) Ray {
	offset := c.sampleSquare()
	pixelSample := c.Pixel00Loc.
		Add(Full(float64(i) + offset.X).Mul(c.PixelDeltaU)).
		Add(Full(float64(j) + offset.Y).Mul(c.PixelDeltaV))
	rayOrigin := c.Center
	rayDirection := pixelSample.Sub(rayOrigin)

	return NewRay(rayOrigin, rayDirection)

}

func (c Camera) sampleSquare() Vec3 {
	return NewVec3(Rand()-0.5, Rand()-0.5, 0)
}

func (c Camera) rayColor(ray Ray, world Hittable, depth int) Color {
	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return NewColor(0., 0., 0.)
	}

	if record, ok := world.Hit(ray, NewInterval(0.001, math.Inf(1))); ok {
		direction := record.Normal.Add(NewRandUnitVec3())
		return Full(0.5).Mul(c.rayColor(NewRay(record.P, direction), world, depth-1))
	}

	unitDirection := ray.Direction.Unit()
	a := 0.5 * (unitDirection.Y + 1.0)
	return Full(1.0 - a).Mul(NewColor(1.0, 1.0, 1.0)).Add(Full(a).Mul(NewColor(0.5, 0.7, 1.0)))
}
