package handler

import (
	"context"
	"fmt"

	pb "github.com/JuanPabloJimenez0250013/orbital-net/sim-service/handler/proto"
	"github.com/JuanPabloJimenez0250013/orbital-net/sim-service/model"
)

type SimServer struct {
	pb.UnimplementedSimServiceServer
}

func (s *SimServer) GetNodes(ctx context.Context, in *pb.Empty) (*pb.NodesResponse, error) {
	var nodes []*pb.Node

	for i := range model.Nodes {
		var visibleNodes []string
		for j := range model.Nodes {
			if i != j && model.Nodes[i].CanView(model.Nodes[j]) {
				visibleNodes = append(visibleNodes, model.Nodes[j].ID)
			}
		}

		xPos, yPos := model.Nodes[i].XYPosition()
		nodes = append(nodes, &pb.Node{
			Name:      model.Nodes[i].Name,
			Id:        model.Nodes[i].ID,
			Speed:     model.Nodes[i].GetLinearSpeed() / model.SIMULATION_SCALE_TO_ONE,
			Altitude:  (model.Nodes[i].Orbit.Radius - model.Nodes[i].ParentPlanet.Radius) / model.SIMULATION_SCALE_TO_ONE,
			XPos:      xPos / model.SIMULATION_SCALE_TO_ONE,
			YPos:      yPos / model.SIMULATION_SCALE_TO_ONE,
			Interface: fmt.Sprintf("%dxG%d", model.Nodes[i].Interface.PortQuantity, model.Nodes[i].Interface.PortGen),
			CanView:   visibleNodes,
		})
	}

	return &pb.NodesResponse{Nodes: nodes}, nil
}
