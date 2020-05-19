package main

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

func scan(image [][]byte) rectangle {
	lenX := len(image)
	lenY := len(image[0])
	captWindH := 180
	captWindL := 180
	imgScan := scanImage{
		data: image,
		lenY: lenY,
		lenX: lenX,
	}

	possibleRes := make([]possibleRect, 0, 6000)

	for x := 3; x < lenX-captWindH; x++ {
		for y := 3; y < lenY-captWindL; y++ {
			rect := rectangle{
				start: point{
					x: x,
					y: y,
				},
				end: point{
					x: x + captWindH,
					y: y + captWindL,
				},
			}
			r := capturePossibleRect(imgScan, rect)
			possibleRes = append(possibleRes, r)
		}
	}

	return getMaxPossibleRect(possibleRes).rect
}

func capturePossibleRect(image scanImage, rect rectangle) possibleRect {
	sum := 0
	for x := rect.start.x - 3; x <= rect.end.x-3; x++ {
		for y := rect.start.y - 3; y <= rect.end.y-3; y++ {
			sum += int(image.data[x][y])
		}
	}
	return possibleRect{
		sum:  sum,
		rect: rect,
	}
}

func getMaxPossibleRect(arr []possibleRect) possibleRect {
	max := possibleRect{
		sum:  0,
		rect: rectangle{},
	}

	for _, elem := range arr {
		if elem.sum > max.sum {
			max = elem
		}
	}

	return max
}
