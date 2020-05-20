package main

import (
	"log"
	"sync"
)

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
	log.Printf("Scan starting.\n")
	lenX := len(image)
	lenY := len(image[0])
	cache := make([]rectangle, 0, cacheSize)

	for x := 0; x < lenX; x++ {
		for y := 0; y < lenY; y++ {
			p := point{x: x, y: y}
			if image[x][y] != 1 || inCache(p, cache) {
				continue
			} else {
				rect := dynamicScan(p, image)
				cache = append(cache, rect)
			}
		}
	}

	return cache
}

//imgCh chan [][]byte
func getNewByteImage(xm, ym int) [][]byte {
	log.Printf("Getting clear byte image.\n")
	arr := make([][]byte, xm)
	for i := 0; i < xm; i++ {
		arr[i] = make([]byte, ym)
	}
	return arr
	// imgCh <- arr
}

func fillByRecursive(start point, image [][]byte) [][]byte {
	log.Printf("Getting figure points recursively, start point(%d, %d)\n", start.x, start.y)
	//newImageCh := make(chan [][]byte)
	// go getNewByteImage(len(image), len(image[0]), newImageCh)
	newImage := getNewByteImage(len(image), len(image[0]))
	points := recFindNearestPoints(image, start)
	// newImage := <-newImageCh
	for _, elem := range points {
		newImage[elem.x][elem.y] = 1
	}
	return newImage
}

func dynamicScan(start point, byteImage [][]byte) rectangle {
	image := fillByRecursive(start, byteImage)
	xm := len(byteImage)
	ym := len(byteImage[0])
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
		if rect.start.x-1 > 0 {
			rect.start.x--
		}
		if rect.start.y-1 > 0 {
			rect.start.y--
		}
		if rect.end.x+1 < xm {
			rect.end.x++
		}
		if rect.end.y+1 < ym {
			rect.end.y++
		}
		checkBorders(image, &rect)
		if prevRect.start.x == rect.start.x && prevRect.start.y == rect.start.y && prevRect.end.x == rect.end.x && prevRect.end.y == rect.end.y {
			break
		}
	}
	formatRect(&rect, xm, ym, 5)
	return rect
}

func checkBorders(image [][]byte, rect *rectangle) {
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

func formatRect(rect *rectangle, xMax, yMax, pixelBound int) {
	log.Printf("Formatting rect %v\n", *rect)
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

type rec struct {
	points []point
}

func recFindNearestPoints(image [][]byte, start point) []point {
	points := make([]point, 0, 32000)
	rec := &rec{points: points}
	wg := &sync.WaitGroup{}
	wg.Add(1)
	go recursiveFindNearestPoints(image, start, rec, wg)
	wg.Wait()
	log.Printf("points len:%d\n", len(rec.points))
	return rec.points
}

//Concurrency
func recursiveFindNearestPoints(image [][]byte, p point, rec *rec, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, elem := range rec.points {
		if elem.x == p.x && elem.y == p.y {
			return
		}
	}
	rec.points = append(rec.points, p)
	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 || p.x+dx < 0 || p.y+dy < 0 || p.x+dx > len(image) || p.y+dy > len(image[0]) {
				continue
			}
			if image[p.x+dx][p.y+dy] == 1 {
				wg.Add(1)
				go recursiveFindNearestPoints(image, point{x: p.x + dx, y: p.y + dy}, rec, wg)
			}
		}
	}

}
