package main

import (
	"image"
	"image/color"
	"log"
)

type byteImage struct {
	arr [][]byte
}

func (arr *byteImage) Set(x, y int, b byte) {
	arr.arr[x][y] = b
}

func (arr *byteImage) SetRectangle(x1, y1, x2, y2 int) {
	for x := x1; x <= x2; x++ {
		arr.Set(x, y1, 2)
		arr.Set(x, y2, 2)
	}
	for y := y1; y <= y2; y++ {
		arr.Set(x1, y, 2)
		arr.Set(x2, y, 2)
	}
}

func (arr *byteImage) SetCenters(rects []rectangle) {
	for _, rect := range rects {
		for dx := -3; dx <= 3; dx++ {
			for dy := -3; dy <= 3; dy++ {
				arr.Set(rect.center.x+dx, rect.center.y+dy, 3)
			}
		}
	}
}

func (arr *byteImage) NewRGBAImage() *image.RGBA {
	log.Printf("Converting byte image to jpg\n")
	image := image.NewRGBA(image.Rect(0, 0, 600, 600))
	red := color.RGBA{200, 30, 30, 255}
	green := color.RGBA{0, 230, 64, 1}
	lenX := len(arr.arr)
	lenY := len(arr.arr[0])

	for x := 0; x < lenX; x++ {
		for y := 0; y < lenY; y++ {
			switch arr.arr[x][y] {
			case 1:
				image.Set(x, y, color.Black)
			case 2:
				image.Set(x, y, red)
			case 3:
				image.Set(x, y, green)
			default:
				image.Set(x, y, color.White)
			}
		}
	}

	return image
}
