package physics

import (
	"context"
	"fmt"
	pb "plc-service/pkg/grpc"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type BoilerClient struct {
	conn   *grpc.ClientConn
	client pb.BoilerPhysicsClient
}

func NewClient(address string) (*BoilerClient, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))

	if err != nil {
		return nil, fmt.Errorf("failed to connect to physics service: %w", err)
	}

	client := pb.NewBoilerPhysicsClient(conn)

	return &BoilerClient{
		conn:   conn,
		client: client,
	}, nil
}

func (c *BoilerClient) Close() error {
	return c.conn.Close()
}

func (c *BoilerClient) GetStatus(ctx context.Context) (*pb.BoilerStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return c.client.GetStatus(ctx, &pb.Empty{})
}

func (c *BoilerClient) SetControls(ctx context.Context, fuel, water, steam float64) (*pb.BoilerStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()

	req := &pb.ControlInput{
		FuelValve:      fuel,
		FeedwaterValve: water,
		SteamValve:     steam,
	}
	return c.client.SetControls(ctx, req)
}
