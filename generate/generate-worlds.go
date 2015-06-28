package generate

import (
	"fmt"
	"log"
	"math/rand"
)

// Factions to be removed later into different file
const (
	Red = 1 << iota
	Blue
)

type Region struct {
	Location Coordinate2D
	Fortress *Fortress
}

type Fortress struct {
	Location Coordinate2D
	Parent   *Region
	Towns    []Town

	FactionScore int
}

type Town struct {
	Location Coordinate2D
	Parent   *Fortress
	Outposts []Outpost

	FactionScore int
}

type Outpost struct {
	Location Coordinate2D
	Parent   *Town

	FactionScore int // score 1-200: basic quests // 200-1000: intermediate quests // 1000+ advanced quests
}

var loop_size int
var region_width int
var min_towns int
var max_towns int

// outposts per town
var min_outposts int
var max_outposts int

func init() {
	fmt.Println("Init")

	loop_size = 8
	region_width = 20
	min_towns = 2
	max_towns = 4
	min_outposts = 2
	max_outposts = 4

	// can't divide by zero!
	var zero Coordinate2D
	region_zero := GenerateRegion(zero, region_width)
	printRegion(region_zero)

	// building loop
	var coord Coordinate2D
	for i := 1; i < loop_size; i++ {
		coord = IndexToCoordinate(i)

		var megacoord Coordinate2D
		megacoord.X = coord.X * region_width
		megacoord.Y = coord.Y * region_width

		//GenerateRegion(megacoord, region_width)
	}
}

func GenerateRegion(coord Coordinate2D, region_width int) *Region {
	log.Print("Generating Region @ ", coord)

	region_data := new(Region)
	// Generate A Fortress
	GenerateFortress(region_data)

	// Generate Towns + Update Fortress Struct
	// Generate Outposts + Update Fortress Struct

	return region_data
}

func GenerateFortress(region_data *Region) {
	log.Print("Generating Fortress")

	// generate location
	max_displacement := 2
	x_offset := rand.Intn(max_displacement*2+1) - max_displacement
	y_offset := rand.Intn(max_displacement*2+1) - max_displacement

	var location Coordinate2D
	location.X = x_offset
	location.Y = y_offset

	region_data.Fortress.Location = location

	var town_count int
	town_count = rand.Intn(max_towns-min_towns+1) + min_towns
	log.Print("Town Count: ", town_count)
	for i := 0; i < town_count; i++ {
		GenerateTown(region_data)
	}
}

func GenerateTown(region_data *Region) {
	log.Print("Generating Town")

	// Choose town location
	// Choose starting FactionScore

	// Choose number of outposts
	// Generate Outposts
	var outpost_count int

	outpost_count = rand.Intn(max_outposts-min_outposts+1) + min_outposts
	log.Print("Outpost Count: ", outpost_count)

	for i := 0; i < outpost_count; i++ {
		GenerateOutpost(region_data)
	}

	// Add to region data
}

func GenerateOutpost(region_data *Region) {
	log.Print("Generating Outpost")
}

func printRegion(region_data *Region) {
	message := fmt.Sprintf("Region:%s Fortress:%s\n", region_data.Location, region_data.Fortress.Location)
	log.Print(message)
}
