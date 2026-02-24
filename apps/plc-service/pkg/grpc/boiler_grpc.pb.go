package grpc

import (
	context "context"

	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

const _ = grpc.SupportPackageIsVersion9
const (
	BoilerPhysics_GetStatus_FullMethodName   = "/boiler.BoilerPhysics/GetStatus"
	BoilerPhysics_SetControls_FullMethodName = "/boiler.BoilerPhysics/SetControls"
)

type BoilerPhysicsClient interface {
	GetStatus(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*BoilerStatus, error)
	SetControls(ctx context.Context, in *ControlInput, opts ...grpc.CallOption) (*BoilerStatus, error)
}
type boilerPhysicsClient struct {
	cc grpc.ClientConnInterface
}

func NewBoilerPhysicsClient(cc grpc.ClientConnInterface) BoilerPhysicsClient {
	return &boilerPhysicsClient{cc}
}
func (c *boilerPhysicsClient) GetStatus(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*BoilerStatus, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BoilerStatus)
	err := c.cc.Invoke(ctx, BoilerPhysics_GetStatus_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *boilerPhysicsClient) SetControls(ctx context.Context, in *ControlInput, opts ...grpc.CallOption) (*BoilerStatus, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(BoilerStatus)
	err := c.cc.Invoke(ctx, BoilerPhysics_SetControls_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type BoilerPhysicsServer interface {
	GetStatus(context.Context, *Empty) (*BoilerStatus, error)
	SetControls(context.Context, *ControlInput) (*BoilerStatus, error)
	mustEmbedUnimplementedBoilerPhysicsServer()
}

type UnimplementedBoilerPhysicsServer struct{}

func (UnimplementedBoilerPhysicsServer) GetStatus(context.Context, *Empty) (*BoilerStatus, error) {
	return nil, status.Error(codes.Unimplemented, "method GetStatus not implemented")
}
func (UnimplementedBoilerPhysicsServer) SetControls(context.Context, *ControlInput) (*BoilerStatus, error) {
	return nil, status.Error(codes.Unimplemented, "method SetControls not implemented")
}
func (UnimplementedBoilerPhysicsServer) mustEmbedUnimplementedBoilerPhysicsServer() {}
func (UnimplementedBoilerPhysicsServer) testEmbeddedByValue()                       {}

type UnsafeBoilerPhysicsServer interface {
	mustEmbedUnimplementedBoilerPhysicsServer()
}

func RegisterBoilerPhysicsServer(s grpc.ServiceRegistrar, srv BoilerPhysicsServer) {
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&BoilerPhysics_ServiceDesc, srv)
}
func _BoilerPhysics_GetStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoilerPhysicsServer).GetStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BoilerPhysics_GetStatus_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoilerPhysicsServer).GetStatus(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}
func _BoilerPhysics_SetControls_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ControlInput)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BoilerPhysicsServer).SetControls(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: BoilerPhysics_SetControls_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BoilerPhysicsServer).SetControls(ctx, req.(*ControlInput))
	}
	return interceptor(ctx, in, info, handler)
}

var BoilerPhysics_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "boiler.BoilerPhysics",
	HandlerType: (*BoilerPhysicsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStatus",
			Handler:    _BoilerPhysics_GetStatus_Handler,
		},
		{
			MethodName: "SetControls",
			Handler:    _BoilerPhysics_SetControls_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "boiler.proto",
}
