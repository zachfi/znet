// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package astro

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

// AstroClient is the client API for Astro service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AstroClient interface {
	Sunrise(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	Sunset(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
	PreSunset(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error)
}

type astroClient struct {
	cc grpc.ClientConnInterface
}

func NewAstroClient(cc grpc.ClientConnInterface) AstroClient {
	return &astroClient{cc}
}

func (c *astroClient) Sunrise(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/astro.Astro/Sunrise", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *astroClient) Sunset(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/astro.Astro/Sunset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *astroClient) PreSunset(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/astro.Astro/PreSunset", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AstroServer is the server API for Astro service.
// All implementations must embed UnimplementedAstroServer
// for forward compatibility
type AstroServer interface {
	Sunrise(context.Context, *Empty) (*Empty, error)
	Sunset(context.Context, *Empty) (*Empty, error)
	PreSunset(context.Context, *Empty) (*Empty, error)
	mustEmbedUnimplementedAstroServer()
}

// UnimplementedAstroServer must be embedded to have forward compatible implementations.
type UnimplementedAstroServer struct {
}

func (UnimplementedAstroServer) Sunrise(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sunrise not implemented")
}
func (UnimplementedAstroServer) Sunset(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Sunset not implemented")
}
func (UnimplementedAstroServer) PreSunset(context.Context, *Empty) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PreSunset not implemented")
}
func (UnimplementedAstroServer) mustEmbedUnimplementedAstroServer() {}

// UnsafeAstroServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AstroServer will
// result in compilation errors.
type UnsafeAstroServer interface {
	mustEmbedUnimplementedAstroServer()
}

func RegisterAstroServer(s grpc.ServiceRegistrar, srv AstroServer) {
	s.RegisterService(&Astro_ServiceDesc, srv)
}

func _Astro_Sunrise_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AstroServer).Sunrise(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/astro.Astro/Sunrise",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AstroServer).Sunrise(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Astro_Sunset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AstroServer).Sunset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/astro.Astro/Sunset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AstroServer).Sunset(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Astro_PreSunset_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AstroServer).PreSunset(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/astro.Astro/PreSunset",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AstroServer).PreSunset(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Astro_ServiceDesc is the grpc.ServiceDesc for Astro service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Astro_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "astro.Astro",
	HandlerType: (*AstroServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Sunrise",
			Handler:    _Astro_Sunrise_Handler,
		},
		{
			MethodName: "Sunset",
			Handler:    _Astro_Sunset_Handler,
		},
		{
			MethodName: "PreSunset",
			Handler:    _Astro_PreSunset_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "internal/astro/astro.proto",
}
