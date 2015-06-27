package main

import "fmt"
import "math"
import "math/rand"
import "time"
import "strconv"

type Coord2D struct {
	X int
	Y int
}

type Area struct {
	Width    int
	Location Coord2D
	Nodes    map[Coord2D]AreaNode
}

func NewArea() *Area {
	area := new(Area)
	area.Nodes = make(map[Coord2D]AreaNode)
	return area
}

func (area *Area) AddNode(node AreaNode) {
	area.Nodes[node.Location] = node
}

type NodeType int

const (
	_ NodeType = iota
	City
	Town
	Outpost
)

type AreaNode struct {
	Location Coord2D
	Type     NodeType
}

func main() {
	fmt.Printf("Hello world\n")

	timenow := time.Now()
	fmt.Printf(strconv.Itoa(timenow.Nanosecond()))
	rand.Seed(int64(time.Now().Nanosecond()))

	area := NewArea()
	area.Location = Coord2D{1, 2}
	area.Width = 20

	genCity(area)
	town_count := 4
	genTowns(area, town_count)

	fmt.Println("\n")
	printArea(area)
}

func genCity(area *Area) {

	mid := area.Width / 2
	rand_x := rand.Intn(3) - 1 + mid
	rand_y := rand.Intn(3) - 1 + mid

	var node AreaNode
	node.Location = Coord2D{rand_x, rand_y}
	node.Type = City
	area.AddNode(node)
	fmt.Printf("City@%d,%d\n", rand_x, rand_y)
}

func genTowns(area *Area, count int) {
	for i := 0; i < count; i++ {
		rand_x := rand.Intn(area.Width)
		rand_y := rand.Intn(area.Width)
		location := Coord2D{rand_x, rand_y}
		for !checkTownLoc(area, location) {
			rand_x = rand.Intn(area.Width)
			rand_y = rand.Intn(area.Width)
			location = Coord2D{rand_x, rand_y}
		}

		var node AreaNode
		node.Location = Coord2D{rand_x, rand_y}
		node.Type = Town
		area.AddNode(node)

		genOutposts(area, node.Location, count)
	}

}

func genOutposts(area *Area, town_coord Coord2D, count int) {
	for i := 0; i < count; i++ {
		rand_x := rand.Intn(7) - 3
		rand_y := rand.Intn(7) - 3
		location := Coord2D{rand_x + town_coord.X, rand_y + town_coord.Y}
		for !checkOutpostLoc(area, town_coord, location) {
			rand_x = rand.Intn(area.Width)
			rand_y = rand.Intn(area.Width)
			location = Coord2D{rand_x, rand_y}
		}

		var node AreaNode
		node.Location = Coord2D{rand_x, rand_y}
		node.Type = Outpost
		area.AddNode(node)

	}

}

func checkOutpostLoc(area *Area, town Coord2D, outpost Coord2D) bool {

	if outpost.X < 0 || outpost.X >= area.Width {
		return false
	}

	if outpost.Y < 0 || outpost.Y >= area.Width {
		return false
	}

	for k, v := range area.Nodes {
		fmt.Println()
		distance := dist(town, outpost)
		switch v.Type {
		case City:
			if distance <= 4 {
				return false
			}
		case Town:
			if distance <= 5 {
				return false
			}
		}
	}

	return true
}

func checkTownLoc(area *Area, coord Coord2D) bool {

	for k, v := range area.Nodes {
		fmt.Println(coord.X, coord.Y)
		distance := dist(coord, k)
		switch v.Type {
		case City:
			if distance <= 4 {
				return false
			}
		case Town:
			if distance <= 5 {
				return false
			}
		}
	}

	if coord.X <= 1 || area.Width-coord.X <= 2 {
		return false
	}

	if coord.Y <= 1 || area.Width-coord.Y <= 2 {
		return false
	}

	return true
}

func dist(a Coord2D, b Coord2D) float64 {
	return math.Sqrt(math.Pow((float64(a.X)-float64(b.X)), 2) + (math.Pow(float64(a.Y)-float64(b.Y), 2)))
}

func printArea(area *Area) {

	msg := fmt.Sprintf("Area@%d,%d:\n", area.Location.X, area.Location.Y)

	for x := 0; x < area.Width; x++ {
		for y := 0; y < area.Width; y++ {
			coord := Coord2D{x, y}
			if node, ok := area.Nodes[coord]; ok {
				switch node.Type {
				case City:
					msg += "X"
				case Town:
					msg += "t"
				case Outpost:
					msg += "^"
				}
			} else {
				msg += "."
			}
		}
		msg += "\n"
	}

	fmt.Println(msg)
}
