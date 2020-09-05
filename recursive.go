package main

import (
	"strconv"
	"sync"
)

type rec struct {
	points []point
	uniq   map[string]bool
	mu     *sync.Mutex
	wg     *sync.WaitGroup
}

func (r *rec) recFindNearestPoints(image [][]byte, start point) []point {
	r.points = make([]point, 0, 6000)
	r.uniq = make(map[string]bool)
	r.mu = &sync.Mutex{}
	r.wg = &sync.WaitGroup{}
	// rec := &rec{points: points, uniq: uniq, mu: &sync.Mutex{}}
	r.wg.Add(1)
	go recursiveFindNearestPoints(image, start, r)
	r.wg.Wait()
	return r.points
}

func recursiveFindNearestPoints(image [][]byte, p point, rec *rec) {
	defer rec.wg.Done()
	xString := strconv.Itoa(p.x)
	yString := strconv.Itoa(p.y)
	rec.mu.Lock()
	if _, ok := rec.uniq[xString+yString]; ok {
		rec.mu.Unlock()
		return
	}
	rec.uniq[xString+yString] = true
	rec.mu.Unlock()
	rec.points = append(rec.points, p)

	for dx := -1; dx <= 1; dx++ {
		for dy := -1; dy <= 1; dy++ {
			if dx == 0 && dy == 0 || p.x+dx < 0 || p.y+dy < 0 || p.x+dx >= len(image) || p.y+dy >= len(image[0]) {
				continue
			}
			if image[p.x+dx][p.y+dy] == 1 {
				rec.wg.Add(1)
				go recursiveFindNearestPoints(image, point{x: p.x + dx, y: p.y + dy}, rec)
			}
		}
	}
}
