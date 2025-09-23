package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/JuanPabloJimenez0250013/orbital-net/sim-service/model"
)

type NodeJSON struct {
	Name      string
	ID        string
	Speed     float64
	Altitude  float64
	XPos      float64
	YPos      float64
	Interface string
	CanView   []string
}

func NodesHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	var nodesJSON []NodeJSON

	for i := range model.Nodes {
		var visibleNodes []string
		for j := range model.Nodes {
			if i != j && model.Nodes[i].CanView(model.Nodes[j]) {
				visibleNodes = append(visibleNodes, model.Nodes[j].ID)
			}
		}

		xPos, yPos := model.Nodes[i].XYPosition()
		nodesJSON = append(nodesJSON, NodeJSON{
			Name:      model.Nodes[i].Name,
			ID:        model.Nodes[i].ID,
			Speed:     (model.Nodes[i].GetLinearSpeed()) / model.SIMULATION_SCALE_TO_ONE,
			Altitude:  (model.Nodes[i].Orbit.Radius - model.Nodes[i].ParentPlanet.Radius) / model.SIMULATION_SCALE_TO_ONE,
			XPos:      xPos / model.SIMULATION_SCALE_TO_ONE,
			YPos:      yPos / model.SIMULATION_SCALE_TO_ONE,
			Interface: fmt.Sprintf("%dxG%d", model.Nodes[i].Interface.PortQuantity, model.Nodes[i].Interface.PortGen),
			CanView:   visibleNodes,
		})
	}

	json.NewEncoder(w).Encode(nodesJSON)
}
