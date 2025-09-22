package model

import "math"

var (
	ParentPlanet *Planet
	Nodes        []*Node
)

func CreateSimObjects() {
	ParentPlanet = NewPlanet("Earth", EARTH_RADIUS, EARTH_MASS, EARTH_ROTATION_CYCLE)
	Nodes = []*Node{
		NewOrbitNode("Gio", ParentPlanet, 6371000, 1*math.Pi*2/3, 4, 4),
		NewOrbitNode("Lola", ParentPlanet, 6371000, 2*math.Pi*2/3, 4, 4),
		NewOrbitNode("Martina", ParentPlanet, 6371000, 3*math.Pi*2/3, 4, 4),
		NewGroundNode("Honey", ParentPlanet, 0, 4, 4),
	}
}
