package opcua

import (
	"context"
	"fmt"
	"log"

	"github.com/gopcua/opcua/server"
	"github.com/gopcua/opcua/ua"
)

type Server struct {
	srv          *server.Server
	pressNode    *ua.NodeID
	tempNode     *ua.NodeID
	fuelNode     *ua.NodeID
	setpointNode *ua.NodeID
}

func NewServer(port int) *Server {
	endpoint := fmt.Sprintf("opc.tcp://localhost:%d", port)

	s := server.New(
		server.ListenAddr(fmt.Sprintf(":%d", port)),
		server.EnableReflection(true),
		server.MaxMessageSize(65535),
	)

	return &Server{srv: s}
}

func (s *Server) Start(ctx context.Context) error {
	nm := s.srv.NamespaceManager()
	ns, err := nm.Register("urn:virtual-plc:boiler")
	if err != nil {
		return err
	}

	boilerID := ua.NewNumericNodeID(ns, 1000)
	nm.AddObject(boilerID, "Boiler", ua.NewNumericNodeID(0, 85))

	s.pressNode = ua.NewNumericNodeID(ns, 1001)
	nm.AddVariable(s.pressNode, "Pressure", 0.0, boilerID)

	s.tempNode = ua.NewNumericNodeID(ns, 1002)
	nm.AddVariable(s.tempNode, "Temperature", 0.0, boilerID)

	s.fuelNode = ua.NewNumericNodeID(ns, 1003)
	nm.AddVariable(s.fuelNode, "FuelValve", 0.0, boilerID)

	s.setpointNode = ua.NewNumericNodeID(ns, 1004)
	nm.AddVariable(s.setpointNode, "Setpoint", 60.0, boilerID)

	log.Printf("📡 OPC UA Server starting on port 4840...")
	return s.srv.Start(ctx)
}

func (s *Server) Stop() {
	if s.srv != nil {
		s.srv.Stop()
	}
}

func (s *Server) UpdateData(pressure, temp, fuel, setpoint float64) {
	nm := s.srv.NamespaceManager()
	nm.SetVariableValue(s.pressNode, ua.MustVariant(pressure))
	nm.SetVariableValue(s.tempNode, ua.MustVariant(temp))
	nm.SetVariableValue(s.fuelNode, ua.MustVariant(fuel))
	nm.SetVariableValue(s.setpointNode, ua.MustVariant(setpoint))
}
