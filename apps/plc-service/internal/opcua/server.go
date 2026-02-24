package opcua

import (
	"context"
	"log"

	"github.com/gopcua/opcua/server"
	"github.com/gopcua/opcua/ua"
)

type Server struct {
	srv          *server.Server
	ns           server.NameSpace
	pressNode    *ua.NodeID
	tempNode     *ua.NodeID
	fuelNode     *ua.NodeID
	setpointNode *ua.NodeID
	levelNode    *ua.NodeID
	flowNode     *ua.NodeID
}

func NewServer(port int) *Server {
	srv := server.New(
		server.EndPoint("opc.tcp://0.0.0.0", port),
		server.ServerName("Virtual Boiler PLC"),
	)
	return &Server{srv: srv}
}

func (s *Server) Start(ctx context.Context) error {
	s.ns = server.NewNameSpace("urn:virtual-plc:boiler")
	nsIdx := s.srv.AddNamespace(s.ns)

	boilerID := ua.NewNumericNodeID(uint16(nsIdx), 1000)
	boilerNode := server.NewFolderNode(boilerID, "Boiler")
	s.ns.AddNode(boilerNode)

	s.pressNode = ua.NewNumericNodeID(uint16(nsIdx), 1001)
	boilerNode.AddVariable(server.NewVariableNode(s.pressNode, "Pressure", 0.0))

	s.tempNode = ua.NewNumericNodeID(uint16(nsIdx), 1002)
	boilerNode.AddVariable(server.NewVariableNode(s.tempNode, "Temperature", 0.0))

	s.fuelNode = ua.NewNumericNodeID(uint16(nsIdx), 1003)
	boilerNode.AddVariable(server.NewVariableNode(s.fuelNode, "FuelValve", 0.0))

	s.setpointNode = ua.NewNumericNodeID(uint16(nsIdx), 1004)
	boilerNode.AddVariable(server.NewVariableNode(s.setpointNode, "Setpoint", 60.0))

	s.levelNode = ua.NewNumericNodeID(uint16(nsIdx), 1005)
	boilerNode.AddVariable(server.NewVariableNode(s.levelNode, "DrumLevel", 500.0))

	s.flowNode = ua.NewNumericNodeID(uint16(nsIdx), 1006)
	boilerNode.AddVariable(server.NewVariableNode(s.flowNode, "SteamFlow", 0.0))

	log.Printf("OPC UA Server starting on %v...", s.srv.URLs())
	return s.srv.Start(ctx)
}

func (s *Server) Stop() {
	if s.srv != nil {
		s.srv.Close()
	}
}

func (s *Server) GetSetpoint() float64 {
	if s.setpointNode == nil {
		return 60.0
	}
	if n := s.srv.Node(s.setpointNode); n != nil {
		if dv := n.Value(); dv != nil && dv.Value != nil {
			if val, ok := dv.Value.Value().(float64); ok {
				return val
			}
		}
	}
	return 60.0
}

func (s *Server) UpdateData(pressure, temp, fuel, level, flow float64) {
	update := func(nid *ua.NodeID, value float64) {
		if nid == nil {
			return
		}
		if n := s.srv.Node(nid); n != nil {
			dv := &ua.DataValue{
				Value: ua.MustVariant(value),
			}
			n.SetAttribute(ua.AttributeIDValue, dv)
		}
	}

	update(s.pressNode, pressure)
	update(s.tempNode, temp)
	update(s.fuelNode, fuel)
	update(s.levelNode, level)
	update(s.flowNode, flow)
}
