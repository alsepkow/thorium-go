package main

import (
	"fmt"

	"code.google.com/p/intmath/intgr"
)

type Position2D struct {
	X float32
	Y float32
}

type Coordinate2D struct {
	X int
	Y int
}

func (c *Coordinate2D) ToString() string {
	return fmt.Sprintf("(%d,%d)", c.X, c.Y)
}

func (c *Coordinate2D) GetIndex() int {
	return 0
}

func IndexToCoordinate(index int) Coordinate2D {
	var coord Coordinate2D
	coord.X, coord.Y = position(index)
	return coord
}

func first(cycle int) int {
	x := 2*cycle - 1
	return x * x
}

func cycle(index int) int {
	return (intgr.Sqrt(index) + 1) / 2
}

func length(cycle int) int {
	return 8 * cycle
}

func sector(index int) int {
	c := cycle(index)
	offset := index - first(c)
	n := length(c)
	return 4 * offset / n
}

func position(index int) (int, int) {
	c := cycle(index)
	s := sector(index)
	offset := index - first(c) - s*length(c)/4

	switch s {
	case 0:
		return -c, -c + offset + 1
	case 1:
		return -c + offset + 1, c
	case 2:
		return c, c - offset - 1
	default:
		return c - offset - 1, -c
	}

}

func main() {
	var coord Coordinate2D
	var msg string
	var rev int
	for i := 1; i < 1000; i++ {
		coord = IndexToCoordinate(i)
		msg = coord.ToString()
		rev = coord.GetIndex()
		fmt.Printf("start=%d, calc=%s, reverse=%d\n", i, msg, rev)
	}

}
