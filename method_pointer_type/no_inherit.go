package main

// (t T)           T and *T
// (t *T)          *T

import (
	"fmt"
	"math"
)

type circle struct {
	radius float64
}

type shape interface {
	area() float64
}

func (c *circle) area() float64 {
	//Circle does not inherit shape correctly her as param is pointer
	return math.Pi * c.radius * c.radius
}

func info(s shape) {
	fmt.Println("area :", s.area())
}

func main() {
	c := circle{
		radius: 5,
	}
	//info(c) //This wont work as info accepts pointer
	fmt.Println("area", c.area())
}
