package main

import (
	"image"
	"image/color"
)

type setByteImage interface {
	Set(x, y int, b byte)
}

type byteImage struct {
	arr [][]byte
}

func (arr *byteImage) Set(x, y int, b byte) {
	arr.arr[x][y] = b
}

func (arr *byteImage) SetRectangle(x1, y1, x2, y2 int, b byte) {
	for x := x1; x <= x2; x++ {
		arr.Set(x, y1, b)
		arr.Set(x, y2, b)
	}
	for y := y1; y <= y2; y++ {
		arr.Set(x1, y, b)
		arr.Set(x2, y, b)
	}
}

func (arr *byteImage) convertOwnBytesToImage() *image.RGBA {
	image := image.NewRGBA(image.Rect(0, 0, 600, 600))
	red := color.RGBA{200, 30, 30, 255}
	lenX := len(arr.arr)
	lenY := len(arr.arr[0])

	for x := 0; x < lenX; x++ {
		for y := 0; y < lenY; y++ {
			switch arr.arr[x][y] {
			case 1:
				image.Set(x, y, color.Black)
			case 2:
				image.Set(x, y, red)
			default:
				image.Set(x, y, color.White)
			}
		}
	}

	return image
}
