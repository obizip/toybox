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
	// Rendered image width in pixel count
	ImageWidth int
	// Rendered image height
	ImageHeight int
	// Count of random samples for each pixel
	SamplePerPixel int
	// Color scale factor for a sum of pixel samples
	PixelSamplesScale float64
	// Camera center
	Center Point3
	// Location of pixel 0, 0
	Pixel00Loc Point3
	// Offset to pixel to the right
	PixelDeltaU Vec3
	// Offset to pixel below
	PixelDeltaV Vec3
	// Maximum number of ray bounces into scene
	MaxDepth int
	// Vertical view angle (field of view)
	VerticalFov float64
	// Point camera is looking from
	LookFrom Point3
	// Point camera is looking at
	LookAt Point3
	// Camera-relative "up" direction
	VerticalUp Vec3
	// Camera frame basis vectors
	U, V, W Vec3
	// Variation angle fo rays through each pixel
	DefocusAngle float64
	// Distance from camera lookfrom point to plane of perfect focus
	FocusDistance float64
	// Defocus disk horizontal radius
	DefocusDiskU Vec3
	// Defocus disk vertical radisu
	DefocusDiskV Vec3
}

func NewCamera(logger *slog.Logger, world Hittable, aspectRatio float64, imageWidth int, samplesPerPixel int, maxDepth int, verticalFov float64, lookFrom, lookAt Point3, verticalUp Vec3, defocusAngle, focusDistance float64) Camera {
	//
	// Image
	//
	imageHeight := max(int(float64(imageWidth)/aspectRatio), 1)
	pixelSamplesScale := 1.0 / float64(samplesPerPixel)

	//
	// Camera
	//
	center := lookFrom
	theta := DegreeToRadians(verticalFov)
	h := math.Tan(theta / 2)
	viewportHeight := 2.0 * h * focusDistance
	viewportWidth := viewportHeight * float64(imageWidth) / float64(imageHeight)

	// Calculate the u,v,w unit basis vectors for the camera coordinate frame.
	w := lookFrom.Sub(lookAt).Unit()
	u := verticalUp.Cross(w).Unit()
	v := w.Cross(u)

	// Calculate the vectors across the horizontal and down the vertical viewport edges.
	viewportU := Full(viewportWidth).Mul(u)
	viewportV := Full(viewportHeight).Mul(v.Neg())

	// Calculate the horizontal and vertical delta vectors from pixel to pixel.
	pixelDeltaU := viewportU.Div(Full(float64(imageWidth)))
	pixelDeltaV := viewportV.Div(Full(float64(imageHeight)))

	// Calculate the location of the upper left pixel.
	viewportUpperLeft := center.Sub(Full(focusDistance).Mul(w)).Sub(viewportU.Div(Full(2))).Sub(viewportV.Div(Full(2)))
	pixel00Loc := viewportUpperLeft.Add(Full(0.5).Mul(pixelDeltaU.Add(pixelDeltaV)))

	// Calculate the camera defocus disk basis vecotrs.
	defocusRadius := focusDistance * math.Tan(DegreeToRadians(defocusAngle/2))
	defocusDiskU := u.Mul(Full(defocusRadius))
	defocusDiskV := v.Mul(Full(defocusRadius))

	return Camera{
		Logger:            logger,
		AspectRatio:       aspectRatio,
		ImageWidth:        imageWidth,
		ImageHeight:       imageHeight,
		SamplePerPixel:    samplesPerPixel,
		PixelSamplesScale: pixelSamplesScale,
		Center:            center,
		Pixel00Loc:        pixel00Loc,
		PixelDeltaU:       pixelDeltaU,
		PixelDeltaV:       pixelDeltaV,
		MaxDepth:          maxDepth,
		VerticalFov:       verticalFov,
		LookFrom:          lookFrom,
		LookAt:            lookAt,
		VerticalUp:        verticalUp,
		U:                 u,
		V:                 v,
		W:                 w,
		DefocusAngle:      defocusAngle,
		FocusDistance:     focusDistance,
		DefocusDiskU:      defocusDiskU,
		DefocusDiskV:      defocusDiskV,
	}
}

func (c Camera) Render(world Hittable) {
	fmt.Print("P3\n", c.ImageWidth, c.ImageHeight, "\n255\n")

	for j := range c.ImageHeight {
		c.Logger.Info("Scanning lines", "remainingLines", c.ImageHeight-j)
		for i := range c.ImageWidth {
			pixelColor := NewColor(0., 0., 0.)
			for range c.SamplePerPixel {
				ray := c.getRay(i, j)
				pixelColor = pixelColor.Add(c.rayColor(ray, world, c.MaxDepth))
			}
			pixelColor = Full(c.PixelSamplesScale).Mul(pixelColor)
			pixelColor.Write(os.Stdout)
		}
	}
	c.Logger.Info("Finished")
}

func (c Camera) getRay(i, j int) Ray {
	// Construct a camera ray originating from the defocus disk
	// and directed at a randomly sampled point around the pixel location i, j
	offset := c.sampleSquare()
	pixelSample := c.Pixel00Loc.
		Add(Full(float64(i) + offset.X).Mul(c.PixelDeltaU)).
		Add(Full(float64(j) + offset.Y).Mul(c.PixelDeltaV))
	var rayOrigin Vec3
	if c.DefocusAngle <= 0 {
		rayOrigin = c.Center
	} else {
		rayOrigin = c.defocusDiskSample()
	}
	rayDirection := pixelSample.Sub(rayOrigin)

	return NewRay(rayOrigin, rayDirection)

}

func (c Camera) sampleSquare() Vec3 {
	return NewVec3(Rand()-0.5, Rand()-0.5, 0)
}

func (c Camera) defocusDiskSample() Vec3 {
	p := NewRandInUnitDisk()
	return c.Center.Add(Full(p.X).Mul(c.DefocusDiskU)).Add(Full(p.Y).Mul(c.DefocusDiskV))
}

func (c Camera) rayColor(ray Ray, world Hittable, depth int) Color {
	// If we've exceeded the ray bounce limit, no more light is gathered.
	if depth <= 0 {
		return NewColor(0., 0., 0.)
	}

	if record, ok := world.Hit(ray, NewInterval(0.001, math.Inf(1))); ok {
		if attenuation, scattered, ok := record.Material.Scatter(ray, record); ok {
			return attenuation.Mul(c.rayColor(scattered, world, depth-1))
		}
		return NewColor(0., 0., 0.)
	}

	unitDirection := ray.Direction.Unit()
	a := 0.5 * (unitDirection.Y + 1.0)
	return Full(1.0 - a).Mul(NewColor(1.0, 1.0, 1.0)).Add(Full(a).Mul(NewColor(0.5, 0.7, 1.0)))
}
