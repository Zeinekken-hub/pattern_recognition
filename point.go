package main

import "fmt"

type point struct {
	x int
	y int
}

func (p *point) String() string {
	return fmt.Sprintf("(x:%d;y:%d)", p.x, p.y)
}

func (p *point) equal(p2 point) bool {
	return p.x == p2.x && p.y == p2.y
}
