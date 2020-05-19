package main

const (
	cacheSize = 50
)

type point struct {
	x int
	y int
}

type rectangle struct {
	start point
	end   point
}

func inCache(point point, cache []rectangle) bool {
	for _, elem := range cache {
		if elem.start.x < point.x && point.x < elem.end.x && elem.start.y < point.y && point.y < elem.end.y {
			return true
		}
	}
	return false
}

func scan(image [][]byte) []rectangle {
	lenX := len(image)
	lenY := len(image[0])
	cache := make([]rectangle, 0, cacheSize)

	for x := 0; x < lenX-0; x++ {
		for y := 0; y < lenY-0; y++ {
			p := point{x: x, y: y}

			if inCache(p, cache) || image[x][y] != 1 {
				continue
			} else {
				rect := dynamicScan(p, image)
				cache = append(cache, rect)
			}
		}
	}
	return cache
}

func dynamicScan(start point, image [][]byte) rectangle {
	lenX := len(image)
	lenY := len(image[0])
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
		//fmt.Printf("Rect: start(%d,%d)-end(%d,%d)\n", rect.start.x, rect.start.y, rect.end.x, rect.end.y)
		if rect.start.x-1 > 0 {
			rect.start.x--
		}
		if rect.start.y-1 > 0 {
			rect.start.y--
		}
		if rect.end.x+1 < lenX {
			rect.end.x++
		}
		if rect.end.y+1 < lenY {
			rect.end.y++
		}
		check(image, &rect)
		if prevRect.start.x == rect.start.x && prevRect.start.y == rect.start.y && prevRect.end.x == rect.end.x && prevRect.end.y == rect.end.y {
			break
		}
	}

	formatRect(&rect, lenX, lenY, 10)

	return rect
}

func formatRect(rect *rectangle, xMax, yMax, pixelBound int) {
	if rect.start.x-pixelBound < 1 {
		rect.start.x = 1
	} else {
		rect.start.x -= pixelBound
	}
	if rect.start.y-pixelBound < 1 {
		rect.start.y = 1
	} else {
		rect.start.y -= pixelBound
	}
	if rect.end.x+pixelBound > xMax-1 {
		rect.end.x = xMax - 1
	} else {
		rect.end.x += pixelBound
	}
	if rect.end.y+pixelBound > yMax-1 {
		rect.end.y = yMax - 1
	} else {
		rect.end.y += pixelBound
	}
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
