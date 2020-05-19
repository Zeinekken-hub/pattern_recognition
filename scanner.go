package main

import (
	"fmt"
)

type point struct {
	x int
	y int
}

type rectangle struct {
	start point
	end   point
}

type scanImage struct {
	data [][]byte
	lenX int
	lenY int
}

type possibleRect struct {
	sum  int
	rect rectangle
}

func scan(image [][]byte) []rectangle {
	lenX := len(image)
	lenY := len(image[0])
	cache := make([]rectangle, 0, 50)
	for x := 3; x < lenX-3; x++ {
		for y := 3; y < lenY-3; y++ {

			if len(cache) != 0 {
				cacheWork := false
				for _, elem := range cache {
					if elem.start.x < x && x < elem.end.x && elem.start.y < y && y < elem.end.y {
						cacheWork = true
						break
					}
				}
				if cacheWork || image[x][y] != 1 {
					continue
				} else {
					r := kkk(point{x: x, y: y}, image)
					cache = append(cache, r)
				}
			} else {
				if image[x][y] == 1 {
					r := kkk(point{x: x, y: y}, image)
					cache = append(cache, r)
				}
			}
		}
	}
	return cache
}

func kkk(start point, image [][]byte) rectangle {
	rect := rectangle{
		start: point{
			x: start.x - 5,
			y: start.y - 5,
		},
		end: point{
			x: start.x + 5,
			y: start.y + 5,
		},
	}
	for {
		prevRect := rectangle{
			start: point{
				x: rect.start.x,
				y: rect.start.y,
			},
			end: point{
				x: rect.end.x,
				y: rect.end.y,
			},
		}
		fmt.Printf("Rect: start(%d,%d)-end(%d,%d)\n", rect.start.x, rect.start.y, rect.end.x, rect.end.y)
		if rect.start.x-1 > 0 {
			rect.start.x--
		}
		if rect.start.y-1 > 0 {
			rect.start.y--
		}
		if rect.end.x+1 < 600 {
			rect.end.x++
		}
		if rect.end.y+1 < 600 {
			rect.end.y++
		}
		check(image, &rect)
		if prevRect.start.x == rect.start.x && prevRect.start.y == rect.start.y && prevRect.end.x == rect.end.x && prevRect.end.y == rect.end.y {
			break
		}
	}

	rect.start.x -= 10
	rect.start.y -= 10
	rect.end.x += 10
	rect.end.y += 10

	return rect
}

func check(image [][]byte, rect *rectangle) {
	if !checkYMin(image, *rect) {
		rect.start.x++
	}
	if !checkXMin(image, *rect) {
		rect.start.y++
	}
	if !checkYPlus(image, *rect) {
		rect.end.x--
	}
	if !checkXPlus(image, *rect) {
		rect.end.y--
	}
}

func checkXMin(image [][]byte, rect rectangle) bool {
	for x := rect.start.x; x <= rect.end.x; x++ {
		if image[x][rect.start.y-1] == 1 {
			return true
		}
	}

	return false
}

func checkXPlus(image [][]byte, rect rectangle) bool {
	for x := rect.start.x; x <= rect.end.x; x++ {
		if image[x][rect.end.y+1] == 1 {
			return true
		}
	}

	return false
}

func checkYMin(image [][]byte, rect rectangle) bool {
	for y := rect.start.y; y <= rect.end.y; y++ {
		if image[rect.start.x-1][y] == 1 {
			return true
		}
	}

	return false
}

func checkYPlus(image [][]byte, rect rectangle) bool {
	for y := rect.start.y; y <= rect.end.y; y++ {
		if image[rect.end.x+1][y] == 1 {
			return true
		}
	}
	return false
}
