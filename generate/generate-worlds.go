package generate

import (
	"log"
	"math/rand"
)

type Megacity struct {
	Location Coordinate2D
	Outposts []Outpost
}

type Town struct {
	Parent   Megacity
	Outposts []Outpost
	Location Coordinate2D
}

type Outpost struct {
	Parent   Town
	Location Coordinate2D
}

var loop_size int
var city_span int
var min_subbases int
var max_subbases int

func Start() {

	loop_size = 8
	city_span = 8
	min_subbases = 3
	max_subbases = 8

	var coord Coordinate2D

	var zero Coordinate2D
	GenerateMegaCity(zero)

	for i := 1; i < loop_size; i++ {
		coord = IndexToCoordinate(i)

		var megacoord Coordinate2D
		megacoord.X = coord.X * city_span
		megacoord.Y = coord.Y * city_span

		GenerateMegaCity(megacoord)
	}
}

func GenerateMegaCity(coord Coordinate2D) {
	log.Print("Generating Mega City @ ", coord)

	num_subbases := rand.Intn(max_subbases-min_subbases) + min_subbases
	for i := 0; i < num_subbases; i++ {
		GenerateSubBase(coord, i)
	}
}

func GenerateSubBase(coord Coordinate2D, index int) {
	log.Print("Generating Sub Base ", index, " @ ", coord)
}
