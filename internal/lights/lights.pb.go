// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.17.3
// source: internal/lights/lights.proto

package lights

import (
	proto "github.com/golang/protobuf/proto"
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

type State struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	On         bool  `protobuf:"varint,1,opt,name=on,proto3" json:"on,omitempty"`
	Brightness int32 `protobuf:"varint,2,opt,name=brightness,proto3" json:"brightness,omitempty"`
}

func (x *State) Reset() {
	*x = State{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_lights_lights_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *State) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*State) ProtoMessage() {}

func (x *State) ProtoReflect() protoreflect.Message {
	mi := &file_internal_lights_lights_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use State.ProtoReflect.Descriptor instead.
func (*State) Descriptor() ([]byte, []int) {
	return file_internal_lights_lights_proto_rawDescGZIP(), []int{0}
}

func (x *State) GetOn() bool {
	if x != nil {
		return x.On
	}
	return false
}

func (x *State) GetBrightness() int32 {
	if x != nil {
		return x.Brightness
	}
	return 0
}

type LightGroup struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type   string  `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Id     int32   `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Name   string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	State  *State  `protobuf:"bytes,4,opt,name=state,proto3" json:"state,omitempty"`
	Lights []int32 `protobuf:"varint,5,rep,packed,name=lights,proto3" json:"lights,omitempty"`
}

func (x *LightGroup) Reset() {
	*x = LightGroup{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_lights_lights_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LightGroup) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LightGroup) ProtoMessage() {}

func (x *LightGroup) ProtoReflect() protoreflect.Message {
	mi := &file_internal_lights_lights_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LightGroup.ProtoReflect.Descriptor instead.
func (*LightGroup) Descriptor() ([]byte, []int) {
	return file_internal_lights_lights_proto_rawDescGZIP(), []int{1}
}

func (x *LightGroup) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *LightGroup) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *LightGroup) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *LightGroup) GetState() *State {
	if x != nil {
		return x.State
	}
	return nil
}

func (x *LightGroup) GetLights() []int32 {
	if x != nil {
		return x.Lights
	}
	return nil
}

type Light struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Type  string `protobuf:"bytes,1,opt,name=type,proto3" json:"type,omitempty"`
	Id    int32  `protobuf:"varint,2,opt,name=id,proto3" json:"id,omitempty"`
	Name  string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	State *State `protobuf:"bytes,4,opt,name=state,proto3" json:"state,omitempty"`
}

func (x *Light) Reset() {
	*x = Light{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_lights_lights_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Light) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Light) ProtoMessage() {}

func (x *Light) ProtoReflect() protoreflect.Message {
	mi := &file_internal_lights_lights_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Light.ProtoReflect.Descriptor instead.
func (*Light) Descriptor() ([]byte, []int) {
	return file_internal_lights_lights_proto_rawDescGZIP(), []int{2}
}

func (x *Light) GetType() string {
	if x != nil {
		return x.Type
	}
	return ""
}

func (x *Light) GetId() int32 {
	if x != nil {
		return x.Id
	}
	return 0
}

func (x *Light) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Light) GetState() *State {
	if x != nil {
		return x.State
	}
	return nil
}

type LightGroupRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name       string   `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Color      string   `protobuf:"bytes,2,opt,name=color,proto3" json:"color,omitempty"`
	Colors     []string `protobuf:"bytes,3,rep,name=colors,proto3" json:"colors,omitempty"`
	Brightness int32    `protobuf:"varint,4,opt,name=brightness,proto3" json:"brightness,omitempty"`
}

func (x *LightGroupRequest) Reset() {
	*x = LightGroupRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_lights_lights_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LightGroupRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LightGroupRequest) ProtoMessage() {}

func (x *LightGroupRequest) ProtoReflect() protoreflect.Message {
	mi := &file_internal_lights_lights_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LightGroupRequest.ProtoReflect.Descriptor instead.
func (*LightGroupRequest) Descriptor() ([]byte, []int) {
	return file_internal_lights_lights_proto_rawDescGZIP(), []int{3}
}

func (x *LightGroupRequest) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *LightGroupRequest) GetColor() string {
	if x != nil {
		return x.Color
	}
	return ""
}

func (x *LightGroupRequest) GetColors() []string {
	if x != nil {
		return x.Colors
	}
	return nil
}

func (x *LightGroupRequest) GetBrightness() int32 {
	if x != nil {
		return x.Brightness
	}
	return 0
}

type LightResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Lights []*Light      `protobuf:"bytes,1,rep,name=lights,proto3" json:"lights,omitempty"`
	Groups []*LightGroup `protobuf:"bytes,2,rep,name=groups,proto3" json:"groups,omitempty"`
}

func (x *LightResponse) Reset() {
	*x = LightResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_lights_lights_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *LightResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*LightResponse) ProtoMessage() {}

func (x *LightResponse) ProtoReflect() protoreflect.Message {
	mi := &file_internal_lights_lights_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use LightResponse.ProtoReflect.Descriptor instead.
func (*LightResponse) Descriptor() ([]byte, []int) {
	return file_internal_lights_lights_proto_rawDescGZIP(), []int{4}
}

func (x *LightResponse) GetLights() []*Light {
	if x != nil {
		return x.Lights
	}
	return nil
}

func (x *LightResponse) GetGroups() []*LightGroup {
	if x != nil {
		return x.Groups
	}
	return nil
}

var File_internal_lights_lights_proto protoreflect.FileDescriptor

var file_internal_lights_lights_proto_rawDesc = []byte{
	0x0a, 0x1c, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6c, 0x69, 0x67, 0x68, 0x74,
	0x73, 0x2f, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06,
	0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x22, 0x37, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x65, 0x12,
	0x0e, 0x0a, 0x02, 0x6f, 0x6e, 0x18, 0x01, 0x20, 0x01, 0x28, 0x08, 0x52, 0x02, 0x6f, 0x6e, 0x12,
	0x1e, 0x0a, 0x0a, 0x62, 0x72, 0x69, 0x67, 0x68, 0x74, 0x6e, 0x65, 0x73, 0x73, 0x18, 0x02, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x0a, 0x62, 0x72, 0x69, 0x67, 0x68, 0x74, 0x6e, 0x65, 0x73, 0x73, 0x22,
	0x81, 0x01, 0x0a, 0x0a, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x12, 0x12,
	0x0a, 0x04, 0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79,
	0x70, 0x65, 0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02,
	0x69, 0x64, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18,
	0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x53,
	0x74, 0x61, 0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x12, 0x16, 0x0a, 0x06, 0x6c,
	0x69, 0x67, 0x68, 0x74, 0x73, 0x18, 0x05, 0x20, 0x03, 0x28, 0x05, 0x52, 0x06, 0x6c, 0x69, 0x67,
	0x68, 0x74, 0x73, 0x22, 0x64, 0x0a, 0x05, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x0e, 0x0a, 0x02, 0x69, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x02, 0x69, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x53, 0x74, 0x61,
	0x74, 0x65, 0x52, 0x05, 0x73, 0x74, 0x61, 0x74, 0x65, 0x22, 0x75, 0x0a, 0x11, 0x4c, 0x69, 0x67,
	0x68, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12,
	0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61,
	0x6d, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x05, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6c, 0x6f,
	0x72, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x06, 0x63, 0x6f, 0x6c, 0x6f, 0x72, 0x73,
	0x12, 0x1e, 0x0a, 0x0a, 0x62, 0x72, 0x69, 0x67, 0x68, 0x74, 0x6e, 0x65, 0x73, 0x73, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x05, 0x52, 0x0a, 0x62, 0x72, 0x69, 0x67, 0x68, 0x74, 0x6e, 0x65, 0x73, 0x73,
	0x22, 0x62, 0x0a, 0x0d, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73,
	0x65, 0x12, 0x25, 0x0a, 0x06, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x18, 0x01, 0x20, 0x03, 0x28,
	0x0b, 0x32, 0x0d, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74,
	0x52, 0x06, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x12, 0x2a, 0x0a, 0x06, 0x67, 0x72, 0x6f, 0x75,
	0x70, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74,
	0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x06, 0x67, 0x72,
	0x6f, 0x75, 0x70, 0x73, 0x32, 0xa8, 0x03, 0x0a, 0x06, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x12,
	0x39, 0x0a, 0x05, 0x41, 0x6c, 0x65, 0x72, 0x74, 0x12, 0x19, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74,
	0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67,
	0x68, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x03, 0x44, 0x69,
	0x6d, 0x12, 0x19, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6c,
	0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x37, 0x0a, 0x03, 0x4f, 0x66, 0x66, 0x12, 0x19, 0x2e, 0x6c, 0x69, 0x67,
	0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c,
	0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x36, 0x0a, 0x02,
	0x4f, 0x6e, 0x12, 0x19, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68,
	0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e,
	0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3f, 0x0a, 0x0b, 0x52, 0x61, 0x6e, 0x64, 0x6f, 0x6d, 0x43, 0x6f,
	0x6c, 0x6f, 0x72, 0x12, 0x19, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67,
	0x68, 0x74, 0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15,
	0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x3c, 0x0a, 0x08, 0x53, 0x65, 0x74, 0x43, 0x6f, 0x6c, 0x6f,
	0x72, 0x12, 0x19, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74,
	0x47, 0x72, 0x6f, 0x75, 0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6c,
	0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f,
	0x6e, 0x73, 0x65, 0x12, 0x3a, 0x0a, 0x06, 0x54, 0x6f, 0x67, 0x67, 0x6c, 0x65, 0x12, 0x19, 0x2e,
	0x6c, 0x69, 0x67, 0x68, 0x74, 0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x47, 0x72, 0x6f, 0x75,
	0x70, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x15, 0x2e, 0x6c, 0x69, 0x67, 0x68, 0x74,
	0x73, 0x2e, 0x4c, 0x69, 0x67, 0x68, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42,
	0x11, 0x5a, 0x0f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x6c, 0x69, 0x67, 0x68,
	0x74, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_lights_lights_proto_rawDescOnce sync.Once
	file_internal_lights_lights_proto_rawDescData = file_internal_lights_lights_proto_rawDesc
)

func file_internal_lights_lights_proto_rawDescGZIP() []byte {
	file_internal_lights_lights_proto_rawDescOnce.Do(func() {
		file_internal_lights_lights_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_lights_lights_proto_rawDescData)
	})
	return file_internal_lights_lights_proto_rawDescData
}

var file_internal_lights_lights_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_internal_lights_lights_proto_goTypes = []interface{}{
	(*State)(nil),             // 0: lights.State
	(*LightGroup)(nil),        // 1: lights.LightGroup
	(*Light)(nil),             // 2: lights.Light
	(*LightGroupRequest)(nil), // 3: lights.LightGroupRequest
	(*LightResponse)(nil),     // 4: lights.LightResponse
}
var file_internal_lights_lights_proto_depIdxs = []int32{
	0,  // 0: lights.LightGroup.state:type_name -> lights.State
	0,  // 1: lights.Light.state:type_name -> lights.State
	2,  // 2: lights.LightResponse.lights:type_name -> lights.Light
	1,  // 3: lights.LightResponse.groups:type_name -> lights.LightGroup
	3,  // 4: lights.Lights.Alert:input_type -> lights.LightGroupRequest
	3,  // 5: lights.Lights.Dim:input_type -> lights.LightGroupRequest
	3,  // 6: lights.Lights.Off:input_type -> lights.LightGroupRequest
	3,  // 7: lights.Lights.On:input_type -> lights.LightGroupRequest
	3,  // 8: lights.Lights.RandomColor:input_type -> lights.LightGroupRequest
	3,  // 9: lights.Lights.SetColor:input_type -> lights.LightGroupRequest
	3,  // 10: lights.Lights.Toggle:input_type -> lights.LightGroupRequest
	4,  // 11: lights.Lights.Alert:output_type -> lights.LightResponse
	4,  // 12: lights.Lights.Dim:output_type -> lights.LightResponse
	4,  // 13: lights.Lights.Off:output_type -> lights.LightResponse
	4,  // 14: lights.Lights.On:output_type -> lights.LightResponse
	4,  // 15: lights.Lights.RandomColor:output_type -> lights.LightResponse
	4,  // 16: lights.Lights.SetColor:output_type -> lights.LightResponse
	4,  // 17: lights.Lights.Toggle:output_type -> lights.LightResponse
	11, // [11:18] is the sub-list for method output_type
	4,  // [4:11] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_internal_lights_lights_proto_init() }
func file_internal_lights_lights_proto_init() {
	if File_internal_lights_lights_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_lights_lights_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*State); i {
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
		file_internal_lights_lights_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LightGroup); i {
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
		file_internal_lights_lights_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Light); i {
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
		file_internal_lights_lights_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LightGroupRequest); i {
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
		file_internal_lights_lights_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*LightResponse); i {
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
			RawDescriptor: file_internal_lights_lights_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_internal_lights_lights_proto_goTypes,
		DependencyIndexes: file_internal_lights_lights_proto_depIdxs,
		MessageInfos:      file_internal_lights_lights_proto_msgTypes,
	}.Build()
	File_internal_lights_lights_proto = out.File
	file_internal_lights_lights_proto_rawDesc = nil
	file_internal_lights_lights_proto_goTypes = nil
	file_internal_lights_lights_proto_depIdxs = nil
}
