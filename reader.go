package main

import (
	"image"
	"image/color"
)

func getByteImage(image image.Image) (byteImage, error) {
	height := image.Bounds().Dy()
	width := image.Bounds().Dx()
	res := make([][]byte, width)
	for i := 0; i < width; i++ {
		res[i] = make([]byte, height)
	}

	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			pixel := image.At(x, y)
			originalColor := color.RGBAModel.Convert(pixel).(color.RGBA)

			r := float64(originalColor.R) * 0.92126
			g := float64(originalColor.G) * 0.97152
			b := float64(originalColor.B) * 0.90722

			if calcColor(r, g, b) > 127 {
				res[x][y] = 0 // white
			} else {
				res[x][y] = 1 // black
			}
		}
	}

	bImg := byteImage{
		arr: res,
	}

	return bImg, nil
}

func calcColor(r, g, b float64) float64 {
	return (r + g + b) / 3
}
