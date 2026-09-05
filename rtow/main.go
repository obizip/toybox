package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	logger := slog.Default()

	// Image

	imageWidth := 256
	imageHeight := 256

	// Render

	fmt.Print("P3\n", imageWidth, imageHeight, "\n255\n")

	for j := range imageHeight {
		logger.Info("Scaning lines", "remainingLines", imageHeight-j)
		for i := range imageWidth {
			r := float64(i) / float64(imageWidth-1)
			g := float64(j) / float64(imageHeight-1)
			b := 0.0

			color := NewColor(r, g, b)
			color.Write(os.Stdout)
		}
	}
}
