package main

import (
	"image"
	"image/color"
	"log"
	"math"
)

type bImage struct {
	src    [][]byte
	width  int
	height int
	square int
	center point
	cache  []rectangle
}

func (b *bImage) setBounds() {
	b.width = len(b.src)
	b.height = len(b.src[0])
}

func (b *bImage) scan() {
	log.Printf("Scan starting.\n")
	log.Printf("WIDTH: %d, HEIGHT: %d\n", b.width, b.height)

	b.cache = make([]rectangle, 0, 50)

	for x := 5; x < b.width-5; x++ {
		for y := 5; y < b.height-5; y++ {
			p := point{x: x, y: y}
			if b.src[x][y] != 1 || b.inCache(p) {
				continue
			} else {
				rect := b.dynamicScan(p)
				b.cache = append(b.cache, rect)
			}
		}
	}
}

func (b *bImage) inCache(point point) bool {
	for _, elem := range b.cache {
		if elem.start.x < point.x && point.x < elem.end.x && elem.start.y < point.y && point.y < elem.end.y {
			return true
		}
	}
	return false
}

func (b *bImage) dynamicScan(start point) rectangle {
	rec := rec{}
	points := rec.recFindNearestPoints(b.src, start)
	rect := rectangle{start: point{x: math.MaxInt32, y: math.MaxInt32}}

	for i := 0; i < len(points); i++ {
		if rect.start.x > points[i].x {
			rect.start.x = points[i].x
		}
		if rect.start.y > points[i].y {
			rect.start.y = points[i].y
		}
		if rect.end.x < points[i].x {
			rect.end.x = points[i].x
		}
		if rect.end.y < points[i].y {
			rect.end.y = points[i].y
		}
	}
	rect.formatRect(b.width, b.height, 4)
	return rect
}

func (b *bImage) Set(x, y int, c byte) {
	b.src[x][y] = c
}

func (b *bImage) SetRectangle(rect rectangle) {
	for x := rect.start.x; x <= rect.end.x; x++ {
		b.Set(x, rect.start.y, 2)
		b.Set(x, rect.end.y, 2)
	}
	for y := rect.start.y; y <= rect.end.y; y++ {
		b.Set(rect.start.x, y, 2)
		b.Set(rect.end.x, y, 2)
	}
}

func (b *bImage) SetRectangles() {
	for _, rect := range b.cache {
		b.SetRectangle(rect)
	}
}

func (b *bImage) SetCenters() {
	for _, rect := range b.cache {
		for dx := -3; dx <= 3; dx++ {
			for dy := -3; dy <= 3; dy++ {
				c := rect.center()
				b.Set(c.x+dx, c.y+dy, 3)
			}
		}
	}
}

func (b *bImage) NewRGBAImage() *image.RGBA {
	log.Printf("Converting byte image to jpg\n")
	image := image.NewRGBA(image.Rect(0, 0, 600, 600))
	red := color.RGBA{200, 30, 30, 255}
	green := color.RGBA{0, 230, 64, 1}
	lenX := len(b.src)
	lenY := len(b.src[0])

	for x := 0; x < lenX; x++ {
		for y := 0; y < lenY; y++ {
			switch b.src[x][y] {
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

func (b *bImage) getFigureSquare(rect rectangle) int {
	sum := 0
	for x := rect.start.x; x <= rect.end.x; x++ {
		for y := rect.start.y; y <= rect.end.y; y++ {
			if b.src[x][y] != 0 {
				sum++
			}
		}
	}
	return sum
}

func (b *bImage) maxMinFigureSquare() (int, int) {
	maxSquare, minSquare := math.MinInt32, math.MaxInt32

	for _, elem := range b.cache {
		square := b.getFigureSquare(elem)
		maxSquare = int(math.Max(float64(maxSquare), float64(square)))
		minSquare = int(math.Min(float64(minSquare), float64(square)))
	}

	return maxSquare, minSquare
}
