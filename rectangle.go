package main

type rectangle struct {
	start point
	end   point
}

func (r *rectangle) equal(r2 *rectangle) bool {
	return r2.start.x == r.start.x && r2.start.y == r.start.y && r2.end.x == r.end.x && r2.end.y == r.end.y
}

func (r *rectangle) formatRect(xMax, yMax, pixelBound int) {
	if r.start.x-pixelBound < 1 {
		r.start.x = 1
	} else {
		r.start.x -= pixelBound
	}
	if r.start.y-pixelBound < 1 {
		r.start.y = 1
	} else {
		r.start.y -= pixelBound
	}
	if r.end.x+pixelBound > xMax-1 {
		r.end.x = xMax - 1
	} else {
		r.end.x += pixelBound
	}
	if r.end.y+pixelBound > yMax-1 {
		r.end.y = yMax - 1
	} else {
		r.end.y += pixelBound
	}
}

func (r *rectangle) center() point {
	return point{x: (r.start.x + r.end.x) / 2, y: (r.start.y + r.end.y) / 2}
}
