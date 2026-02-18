package physics

import (
	"context"
	pb "plc-service/pkg/grpc"
)

type IPhysicsClient interface {
	GetStatus(ctx context.Context) (*pb.BoilerStatus, error)
	SetControls(ctx context.Context, input *pb.ControlInput) (*pb.BoilerStatus, error)
	Close() error
}
