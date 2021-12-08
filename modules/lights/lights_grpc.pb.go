// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package lights

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// LightsClient is the client API for Lights service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LightsClient interface {
	SetState(ctx context.Context, in *ZoneStateRequest, opts ...grpc.CallOption) (*ZoneStateResponse, error)
}

type lightsClient struct {
	cc grpc.ClientConnInterface
}

func NewLightsClient(cc grpc.ClientConnInterface) LightsClient {
	return &lightsClient{cc}
}

func (c *lightsClient) SetState(ctx context.Context, in *ZoneStateRequest, opts ...grpc.CallOption) (*ZoneStateResponse, error) {
	out := new(ZoneStateResponse)
	err := c.cc.Invoke(ctx, "/lights.Lights/SetState", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LightsServer is the server API for Lights service.
// All implementations must embed UnimplementedLightsServer
// for forward compatibility
type LightsServer interface {
	SetState(context.Context, *ZoneStateRequest) (*ZoneStateResponse, error)
	mustEmbedUnimplementedLightsServer()
}

// UnimplementedLightsServer must be embedded to have forward compatible implementations.
type UnimplementedLightsServer struct {
}

func (UnimplementedLightsServer) SetState(context.Context, *ZoneStateRequest) (*ZoneStateResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetState not implemented")
}
func (UnimplementedLightsServer) mustEmbedUnimplementedLightsServer() {}

// UnsafeLightsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LightsServer will
// result in compilation errors.
type UnsafeLightsServer interface {
	mustEmbedUnimplementedLightsServer()
}

func RegisterLightsServer(s grpc.ServiceRegistrar, srv LightsServer) {
	s.RegisterService(&Lights_ServiceDesc, srv)
}

func _Lights_SetState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ZoneStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LightsServer).SetState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lights.Lights/SetState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LightsServer).SetState(ctx, req.(*ZoneStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Lights_ServiceDesc is the grpc.ServiceDesc for Lights service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Lights_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "lights.Lights",
	HandlerType: (*LightsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SetState",
			Handler:    _Lights_SetState_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "modules/lights/lights.proto",
}
