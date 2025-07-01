package main

import (
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/spf13/cobra"
)

var VERSION string

func img2Ansi(img image.Image, blockWidth, blockHeight int, transparent color.Color, crop bool) {
	bounds := img.Bounds()
	minX, minY, maxX, maxY := bounds.Min.X, bounds.Min.Y, bounds.Max.X, bounds.Max.Y

	if crop && transparent != nil {
		minX, minY, maxX, maxY = getCropBounds(img, transparent)
	}

	height := maxY - minY
	width := maxX - minX

	for yOffset := 0; yOffset < height; yOffset += blockHeight {
		for xOffset := 0; xOffset < width; xOffset += blockWidth {
			x := minX + xOffset
			y := minY + yOffset

			r, g, b, _ := img.At(x, y).RGBA()
			pixelColor := img.At(x, y)
			if transparent != nil && colorsEqual(pixelColor, transparent) {
				fmt.Print("  ")
				continue
			}
			r8 := r >> 8
			g8 := g >> 8
			b8 := b >> 8
			fmt.Printf("\x1b[38;2;%d;%d;%dm██", r8, g8, b8)
		}
		fmt.Print("\x1b[0m\n")
	}
}

func main() {
	var workers int
	var transparentColor string
	var crop bool

	var rootCmd = &cobra.Command{
		Use: "pixel2ansi",
	}
	var aboutCmd = &cobra.Command{
		Use:   "about",
		Short: "Show about info",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(`pixel2ansi is a command-line tool for analyzing and converting pixel art
images. It can detect the smallest repeating pixel unit, handle color
variations with adjustable tolerance, and render the image as ANSI art in
the terminal. Features include transparent color detection, optional
cropping, and parallel processing for speed.

Prefer lossless formats, like PNG or BMP. If you use JPEG, experiment with
different values for tolerence (a good value to start is something like 150).

Source: https://github.com/vrypan/moomsay`)
		},
	}
	var inspectCmd = &cobra.Command{
		Use:   "inspect [file]",
		Short: "Get unit pixel and grid dimensions",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			img := loadImage(args[0])
			imgWidth := img.Bounds().Max.X
			imgHeight := img.Bounds().Max.Y
			minWidth, minHeight := findSmallestBlockSizeParallel(img, workers)
			gridWidth := imgWidth / minWidth
			gridHeight := imgHeight / minHeight
			fmt.Printf("Block size: %dx%d pixels\n", minWidth, minHeight)
			fmt.Printf("Grid size: %dx%d pixels\n", gridWidth, gridHeight)
		},
	}

	var convertCmd = &cobra.Command{
		Use:   "print [file]",
		Short: "Print image as ANSI blocks",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			img := loadImage(args[0])
			minWidth, minHeight := findSmallestBlockSizeParallel(img, workers)
			tColor := getTransparentColor(img, transparentColor)
			img2Ansi(img, minWidth, minHeight, tColor, crop)
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "Get the current version number",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(VERSION)
		},
	}

	inspectCmd.Flags().IntVarP(&workers, "workers", "w", 4, "Number of concurrent workers")
	inspectCmd.Flags().IntVar(&colorTolerance, "tolerance", 0, "Color tolerance for grouping similar colors")

	convertCmd.Flags().IntVarP(&workers, "workers", "w", 4, "Number of concurrent workers")
	convertCmd.Flags().StringVarP(&transparentColor, "transparent", "t", "", "Transparent color (HEX or tl, tr, bl, br)")
	convertCmd.Flags().BoolVar(&crop, "crop", false, "Crop transparent pixels from the output")
	convertCmd.Flags().IntVar(&colorTolerance, "tolerance", 0, "Color tolerance for grouping similar colors")

	rootCmd.AddCommand(inspectCmd)
	rootCmd.AddCommand(convertCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(aboutCmd)

	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
