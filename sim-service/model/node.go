package model

import (
	"fmt"
	"math"
	"sim-service/misc"
)

type Node struct {
	ID           string
	Name         string
	ParentPlanet *Planet
	Orbit        *Orbit
	Interface    *Interface
}

func NewOrbitNode(name string, parentPlanet *Planet, orbitRadius, orbitThetaPos float64, portQuantity int, portGen int) *Node {
	return &Node{ID: "on_" + misc.HashString(name)[:ID_HASH_LEN], Name: name, ParentPlanet: parentPlanet, Orbit: &Orbit{Radius: orbitRadius + parentPlanet.Radius, ThetaPos: orbitThetaPos, ThetaSpeed: math.Sqrt(G * parentPlanet.Mass / math.Pow(parentPlanet.Radius+orbitRadius, 3))}, Interface: &Interface{PortQuantity: portQuantity, PortGen: portGen}}
}

func NewGroundNode(name string, parentPlanet *Planet, positionTheta float64, portQuantity int, portGen int) *Node {
	return &Node{ID: "gn_" + misc.HashString(name)[:ID_HASH_LEN], Name: name, ParentPlanet: parentPlanet, Orbit: &Orbit{Radius: parentPlanet.Radius, ThetaPos: positionTheta, ThetaSpeed: parentPlanet.RotationThetaSpeed}, Interface: &Interface{PortQuantity: portQuantity, PortGen: portGen}}
}

func (node *Node) Print() {
	fmt.Println("------------------------------------")
	fmt.Printf("%s@%s\n", node.Name, node.ID)
	fmt.Printf(" * Speed:\t%.2f m/s\n", node.GetLinearSpeed())
	altitude := node.Orbit.Radius - node.ParentPlanet.Radius
	if altitude <= 0 {
		fmt.Println(" * Altitude:\tGround")
	} else if altitude < 1000 {
		fmt.Printf(" * Altitude:\t%.0f m\n", altitude)
	} else {
		fmt.Printf(" * Altitude:\t%.0f km\n", altitude/1000)
	}
	fmt.Printf(" * Angle:\t%.0f°\n", 180*math.Mod(node.Orbit.ThetaPos, math.Pi*2)/math.Pi)
	xPos, yPos := node.XYPosition()
	fmt.Printf(" * Position:\t[%.2f,%.2f]\n", xPos, yPos)
	fmt.Printf(" * Interface:\t%dxG%d\n", node.Interface.PortQuantity, node.Interface.PortGen)
	fmt.Println("------------------------------------")
}

func (node *Node) GetLinearSpeed() float64 {
	return node.Orbit.ThetaSpeed * node.Orbit.Radius
}

func (node *Node) Move() {
	node.Orbit.ThetaPos += (node.Orbit.ThetaSpeed * TICK_DELAY / 1000) * SIMULATION_SPEED
}

func (node *Node) XYPosition() (float64, float64) {
	return node.Orbit.Radius * math.Cos(node.Orbit.ThetaPos), node.Orbit.Radius * math.Sin(node.Orbit.ThetaPos)
}

func (node1 *Node) CanView(node2 *Node) bool {
	x1, y1 := node1.XYPosition()
	x2, y2 := node2.XYPosition()

	A := x2 - x1
	B := y2 - y1
	C := x2*y1 - x1*y2

	// Distance from planet center to the infinite line
	D := math.Abs(C) / math.Sqrt(A*A+B*B)

	// Projection factor (0..1 if inside segment)
	T := (-(x1*A + y1*B)) / (A*A + B*B)

	// Visible if: line misses the planet, OR closest point is outside the segment
	return D >= node1.ParentPlanet.Radius || T < 0 || T > 1
}
