package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"sync"

	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"

	_ "golang.org/x/image/bmp"
)

var colorTolerance int

func colorsEqual(c1, c2 color.Color) bool {
	r1, g1, b1, _ := c1.RGBA()
	r2, g2, b2, _ := c2.RGBA()

	dr := int(r1>>8) - int(r2>>8)
	dg := int(g1>>8) - int(g2>>8)
	db := int(b1>>8) - int(b2>>8)

	distance := math.Sqrt(float64(dr*dr + dg*dg + db*db))
	return distance <= float64(colorTolerance)
}

func parseHexColor(s string) (color.Color, error) {
	s = strings.TrimPrefix(s, "#")
	if len(s) != 6 {
		return nil, fmt.Errorf("invalid HEX color format")
	}
	r, _ := strconv.ParseUint(s[0:2], 16, 8)
	g, _ := strconv.ParseUint(s[2:4], 16, 8)
	b, _ := strconv.ParseUint(s[4:6], 16, 8)
	return color.RGBA{uint8(r), uint8(g), uint8(b), 255}, nil
}

func loadImage(path string) image.Image {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open image: %v", err)
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		log.Fatalf("Failed to decode image: %v", err)
	}
	return img
}

func getCropBounds(img image.Image, transparent color.Color) (int, int, int, int) {
	bounds := img.Bounds()
	minX, minY := bounds.Max.X, bounds.Max.Y
	maxX, maxY := bounds.Min.X, bounds.Min.Y

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			if transparent != nil && !colorsEqual(img.At(x, y), transparent) {
				if x < minX {
					minX = x
				}
				if x > maxX {
					maxX = x
				}
				if y < minY {
					minY = y
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}
	return minX, minY, maxX + 1, maxY + 1
}

func getTransparentColor(img image.Image, spec string) color.Color {
	if spec == "" {
		return nil
	}

	bounds := img.Bounds()
	switch strings.ToLower(spec) {
	case "tl":
		return img.At(bounds.Min.X, bounds.Min.Y)
	case "tr":
		return img.At(bounds.Max.X-1, bounds.Min.Y)
	case "bl":
		return img.At(bounds.Min.X, bounds.Max.Y-1)
	case "br":
		return img.At(bounds.Max.X-1, bounds.Max.Y-1)
	default:
		c, err := parseHexColor(spec)
		if err != nil {
			log.Fatalf("Invalid color specification: %v", err)
		}
		return c
	}
}
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func processRows(img image.Image, startY, endY int, minWidthCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	width := img.Bounds().Max.X
	localMin := width

	for y := startY; y < endY; y++ {
		runLength := 1
		for x := 1; x < width; x++ {
			if colorsEqual(img.At(x, y), img.At(x-1, y)) {
				runLength++
			} else {
				if runLength < localMin {
					localMin = runLength
				}
				runLength = 1
			}
		}
		if runLength < localMin {
			localMin = runLength
		}
	}

	minWidthCh <- localMin
}

func processCols(img image.Image, startX, endX int, minHeightCh chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	height := img.Bounds().Max.Y
	localMin := height

	for x := startX; x < endX; x++ {
		runLength := 1
		for y := 1; y < height; y++ {
			if colorsEqual(img.At(x, y), img.At(x, y-1)) {
				runLength++
			} else {
				if runLength < localMin {
					localMin = runLength
				}
				runLength = 1
			}
		}
		if runLength < localMin {
			localMin = runLength
		}
	}

	minHeightCh <- localMin
}

func findSmallestBlockSizeParallel(img image.Image, workers int) (int, int) {
	width := img.Bounds().Max.X
	height := img.Bounds().Max.Y

	minWidthCh := make(chan int, workers)
	minHeightCh := make(chan int, workers)

	var wg sync.WaitGroup

	rowsPerWorker := height / workers
	for i := 0; i < workers; i++ {
		startY := i * rowsPerWorker
		endY := startY + rowsPerWorker
		if i == workers-1 {
			endY = height
		}
		wg.Add(1)
		go processRows(img, startY, endY, minWidthCh, &wg)
	}

	colsPerWorker := width / workers
	for i := 0; i < workers; i++ {
		startX := i * colsPerWorker
		endX := startX + colsPerWorker
		if i == workers-1 {
			endX = width
		}
		wg.Add(1)
		go processCols(img, startX, endX, minHeightCh, &wg)
	}

	wg.Wait()
	close(minWidthCh)
	close(minHeightCh)

	minWidth := width
	for w := range minWidthCh {
		minWidth = min(minWidth, w)
	}

	minHeight := height
	for h := range minHeightCh {
		minHeight = min(minHeight, h)
	}

	return minWidth, minHeight
}
