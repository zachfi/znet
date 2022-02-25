// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package pluginutil

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

// PluginMultiplexingClient is the client API for PluginMultiplexing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PluginMultiplexingClient interface {
	MultiplexingSupport(ctx context.Context, in *MultiplexingSupportRequest, opts ...grpc.CallOption) (*MultiplexingSupportResponse, error)
}

type pluginMultiplexingClient struct {
	cc grpc.ClientConnInterface
}

func NewPluginMultiplexingClient(cc grpc.ClientConnInterface) PluginMultiplexingClient {
	return &pluginMultiplexingClient{cc}
}

func (c *pluginMultiplexingClient) MultiplexingSupport(ctx context.Context, in *MultiplexingSupportRequest, opts ...grpc.CallOption) (*MultiplexingSupportResponse, error) {
	out := new(MultiplexingSupportResponse)
	err := c.cc.Invoke(ctx, "/pluginutil.multiplexing.PluginMultiplexing/MultiplexingSupport", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PluginMultiplexingServer is the server API for PluginMultiplexing service.
// All implementations must embed UnimplementedPluginMultiplexingServer
// for forward compatibility
type PluginMultiplexingServer interface {
	MultiplexingSupport(context.Context, *MultiplexingSupportRequest) (*MultiplexingSupportResponse, error)
	mustEmbedUnimplementedPluginMultiplexingServer()
}

// UnimplementedPluginMultiplexingServer must be embedded to have forward compatible implementations.
type UnimplementedPluginMultiplexingServer struct {
}

func (UnimplementedPluginMultiplexingServer) MultiplexingSupport(context.Context, *MultiplexingSupportRequest) (*MultiplexingSupportResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method MultiplexingSupport not implemented")
}
func (UnimplementedPluginMultiplexingServer) mustEmbedUnimplementedPluginMultiplexingServer() {}

// UnsafePluginMultiplexingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PluginMultiplexingServer will
// result in compilation errors.
type UnsafePluginMultiplexingServer interface {
	mustEmbedUnimplementedPluginMultiplexingServer()
}

func RegisterPluginMultiplexingServer(s grpc.ServiceRegistrar, srv PluginMultiplexingServer) {
	s.RegisterService(&PluginMultiplexing_ServiceDesc, srv)
}

func _PluginMultiplexing_MultiplexingSupport_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MultiplexingSupportRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PluginMultiplexingServer).MultiplexingSupport(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pluginutil.multiplexing.PluginMultiplexing/MultiplexingSupport",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PluginMultiplexingServer).MultiplexingSupport(ctx, req.(*MultiplexingSupportRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PluginMultiplexing_ServiceDesc is the grpc.ServiceDesc for PluginMultiplexing service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PluginMultiplexing_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pluginutil.multiplexing.PluginMultiplexing",
	HandlerType: (*PluginMultiplexingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "MultiplexingSupport",
			Handler:    _PluginMultiplexing_MultiplexingSupport_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "sdk/helper/pluginutil/multiplexing.proto",
}
