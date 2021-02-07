// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.24.0
// 	protoc        v3.12.4
// source: pkg/iot/iot.proto

package iot

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_iot_iot_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_iot_iot_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_pkg_iot_iot_proto_rawDescGZIP(), []int{0}
}

// <discovery_prefix>/<component>/[<node_id>/]<object_id>/config
type DeviceDiscovery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	DiscoveryPrefix string   `protobuf:"bytes,1,opt,name=discovery_prefix,json=discoveryPrefix,proto3" json:"discovery_prefix,omitempty"`
	Component       string   `protobuf:"bytes,2,opt,name=component,proto3" json:"component,omitempty"`
	NodeId          string   `protobuf:"bytes,3,opt,name=node_id,json=nodeId,proto3" json:"node_id,omitempty"`
	ObjectId        string   `protobuf:"bytes,4,opt,name=object_id,json=objectId,proto3" json:"object_id,omitempty"`
	Endpoint        []string `protobuf:"bytes,5,rep,name=endpoint,proto3" json:"endpoint,omitempty"`
	Message         []byte   `protobuf:"bytes,6,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *DeviceDiscovery) Reset() {
	*x = DeviceDiscovery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_iot_iot_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *DeviceDiscovery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*DeviceDiscovery) ProtoMessage() {}

func (x *DeviceDiscovery) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_iot_iot_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use DeviceDiscovery.ProtoReflect.Descriptor instead.
func (*DeviceDiscovery) Descriptor() ([]byte, []int) {
	return file_pkg_iot_iot_proto_rawDescGZIP(), []int{1}
}

func (x *DeviceDiscovery) GetDiscoveryPrefix() string {
	if x != nil {
		return x.DiscoveryPrefix
	}
	return ""
}

func (x *DeviceDiscovery) GetComponent() string {
	if x != nil {
		return x.Component
	}
	return ""
}

func (x *DeviceDiscovery) GetNodeId() string {
	if x != nil {
		return x.NodeId
	}
	return ""
}

func (x *DeviceDiscovery) GetObjectId() string {
	if x != nil {
		return x.ObjectId
	}
	return ""
}

func (x *DeviceDiscovery) GetEndpoint() []string {
	if x != nil {
		return x.Endpoint
	}
	return nil
}

func (x *DeviceDiscovery) GetMessage() []byte {
	if x != nil {
		return x.Message
	}
	return nil
}

type Action struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Event  string `protobuf:"bytes,1,opt,name=event,proto3" json:"event,omitempty"`
	Device string `protobuf:"bytes,2,opt,name=device,proto3" json:"device,omitempty"`
	Zone   string `protobuf:"bytes,3,opt,name=zone,proto3" json:"zone,omitempty"`
}

func (x *Action) Reset() {
	*x = Action{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_iot_iot_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Action) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Action) ProtoMessage() {}

func (x *Action) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_iot_iot_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Action.ProtoReflect.Descriptor instead.
func (*Action) Descriptor() ([]byte, []int) {
	return file_pkg_iot_iot_proto_rawDescGZIP(), []int{2}
}

func (x *Action) GetEvent() string {
	if x != nil {
		return x.Event
	}
	return ""
}

func (x *Action) GetDevice() string {
	if x != nil {
		return x.Device
	}
	return ""
}

func (x *Action) GetZone() string {
	if x != nil {
		return x.Zone
	}
	return ""
}

type UpdateRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Device string `protobuf:"bytes,1,opt,name=device,proto3" json:"device,omitempty"`
}

func (x *UpdateRequest) Reset() {
	*x = UpdateRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_pkg_iot_iot_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *UpdateRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*UpdateRequest) ProtoMessage() {}

func (x *UpdateRequest) ProtoReflect() protoreflect.Message {
	mi := &file_pkg_iot_iot_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use UpdateRequest.ProtoReflect.Descriptor instead.
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return file_pkg_iot_iot_proto_rawDescGZIP(), []int{3}
}

func (x *UpdateRequest) GetDevice() string {
	if x != nil {
		return x.Device
	}
	return ""
}

var File_pkg_iot_iot_proto protoreflect.FileDescriptor

var file_pkg_iot_iot_proto_rawDesc = []byte{
	0x0a, 0x11, 0x70, 0x6b, 0x67, 0x2f, 0x69, 0x6f, 0x74, 0x2f, 0x69, 0x6f, 0x74, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x03, 0x69, 0x6f, 0x74, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x22, 0xc6, 0x01, 0x0a, 0x0f, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x44, 0x69, 0x73, 0x63,
	0x6f, 0x76, 0x65, 0x72, 0x79, 0x12, 0x29, 0x0a, 0x10, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65,
	0x72, 0x79, 0x5f, 0x70, 0x72, 0x65, 0x66, 0x69, 0x78, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0f, 0x64, 0x69, 0x73, 0x63, 0x6f, 0x76, 0x65, 0x72, 0x79, 0x50, 0x72, 0x65, 0x66, 0x69, 0x78,
	0x12, 0x1c, 0x0a, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x09, 0x63, 0x6f, 0x6d, 0x70, 0x6f, 0x6e, 0x65, 0x6e, 0x74, 0x12, 0x17,
	0x0a, 0x07, 0x6e, 0x6f, 0x64, 0x65, 0x5f, 0x69, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x06, 0x6e, 0x6f, 0x64, 0x65, 0x49, 0x64, 0x12, 0x1b, 0x0a, 0x09, 0x6f, 0x62, 0x6a, 0x65, 0x63,
	0x74, 0x5f, 0x69, 0x64, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x6f, 0x62, 0x6a, 0x65,
	0x63, 0x74, 0x49, 0x64, 0x12, 0x1a, 0x0a, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x18, 0x05, 0x20, 0x03, 0x28, 0x09, 0x52, 0x08, 0x65, 0x6e, 0x64, 0x70, 0x6f, 0x69, 0x6e, 0x74,
	0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x0c, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x22, 0x4a, 0x0a, 0x06, 0x41, 0x63,
	0x74, 0x69, 0x6f, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x05, 0x65, 0x76, 0x65, 0x6e, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65,
	0x76, 0x69, 0x63, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x7a, 0x6f, 0x6e, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x7a, 0x6f, 0x6e, 0x65, 0x22, 0x27, 0x0a, 0x0d, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63,
	0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x64, 0x65, 0x76, 0x69, 0x63, 0x65, 0x32,
	0x35, 0x0a, 0x03, 0x49, 0x4f, 0x54, 0x12, 0x2e, 0x0a, 0x0c, 0x55, 0x70, 0x64, 0x61, 0x74, 0x65,
	0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x12, 0x12, 0x2e, 0x69, 0x6f, 0x74, 0x2e, 0x55, 0x70, 0x64,
	0x61, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0a, 0x2e, 0x69, 0x6f, 0x74,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x22, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62,
	0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x61, 0x71, 0x75, 0x65, 0x32, 0x30, 0x38, 0x2f, 0x7a, 0x6e,
	0x65, 0x74, 0x2f, 0x70, 0x6b, 0x67, 0x2f, 0x69, 0x6f, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x33,
}

var (
	file_pkg_iot_iot_proto_rawDescOnce sync.Once
	file_pkg_iot_iot_proto_rawDescData = file_pkg_iot_iot_proto_rawDesc
)

func file_pkg_iot_iot_proto_rawDescGZIP() []byte {
	file_pkg_iot_iot_proto_rawDescOnce.Do(func() {
		file_pkg_iot_iot_proto_rawDescData = protoimpl.X.CompressGZIP(file_pkg_iot_iot_proto_rawDescData)
	})
	return file_pkg_iot_iot_proto_rawDescData
}

var file_pkg_iot_iot_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_pkg_iot_iot_proto_goTypes = []interface{}{
	(*Empty)(nil),           // 0: iot.Empty
	(*DeviceDiscovery)(nil), // 1: iot.DeviceDiscovery
	(*Action)(nil),          // 2: iot.Action
	(*UpdateRequest)(nil),   // 3: iot.UpdateRequest
}
var file_pkg_iot_iot_proto_depIdxs = []int32{
	3, // 0: iot.IOT.UpdateDevice:input_type -> iot.UpdateRequest
	0, // 1: iot.IOT.UpdateDevice:output_type -> iot.Empty
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_pkg_iot_iot_proto_init() }
func file_pkg_iot_iot_proto_init() {
	if File_pkg_iot_iot_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_pkg_iot_iot_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_iot_iot_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*DeviceDiscovery); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_iot_iot_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Action); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_pkg_iot_iot_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*UpdateRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_pkg_iot_iot_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_pkg_iot_iot_proto_goTypes,
		DependencyIndexes: file_pkg_iot_iot_proto_depIdxs,
		MessageInfos:      file_pkg_iot_iot_proto_msgTypes,
	}.Build()
	File_pkg_iot_iot_proto = out.File
	file_pkg_iot_iot_proto_rawDesc = nil
	file_pkg_iot_iot_proto_goTypes = nil
	file_pkg_iot_iot_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// IOTClient is the client API for IOT service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type IOTClient interface {
	UpdateDevice(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*Empty, error)
}

type iOTClient struct {
	cc grpc.ClientConnInterface
}

func NewIOTClient(cc grpc.ClientConnInterface) IOTClient {
	return &iOTClient{cc}
}

func (c *iOTClient) UpdateDevice(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/iot.IOT/UpdateDevice", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// IOTServer is the server API for IOT service.
type IOTServer interface {
	UpdateDevice(context.Context, *UpdateRequest) (*Empty, error)
}

// UnimplementedIOTServer can be embedded to have forward compatible implementations.
type UnimplementedIOTServer struct {
}

func (*UnimplementedIOTServer) UpdateDevice(context.Context, *UpdateRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateDevice not implemented")
}

func RegisterIOTServer(s *grpc.Server, srv IOTServer) {
	s.RegisterService(&_IOT_serviceDesc, srv)
}

func _IOT_UpdateDevice_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(IOTServer).UpdateDevice(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/iot.IOT/UpdateDevice",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(IOTServer).UpdateDevice(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _IOT_serviceDesc = grpc.ServiceDesc{
	ServiceName: "iot.IOT",
	HandlerType: (*IOTServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "UpdateDevice",
			Handler:    _IOT_UpdateDevice_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "pkg/iot/iot.proto",
}
