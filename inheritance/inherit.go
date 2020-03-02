package main

import (
	"math"
)

type circle struct {
	radius float64
}

type shape interface {
	area() float64
}

func (c circle) area() float64 { //Circle inherits shape
	return math.Pi * c.radius * c.radius
}

func main() {
	c := circle{
		radius: 5,
	}
	info(c)
}
