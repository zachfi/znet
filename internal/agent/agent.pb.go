// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.17.3
// source: internal/agent/agent.proto

package agent

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

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_agent_agent_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_internal_agent_agent_proto_msgTypes[0]
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
	return file_internal_agent_agent_proto_rawDescGZIP(), []int{0}
}

type Jails struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Jails []*Jail `protobuf:"bytes,1,rep,name=jails,proto3" json:"jails,omitempty"`
}

func (x *Jails) Reset() {
	*x = Jails{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_agent_agent_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Jails) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Jails) ProtoMessage() {}

func (x *Jails) ProtoReflect() protoreflect.Message {
	mi := &file_internal_agent_agent_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Jails.ProtoReflect.Descriptor instead.
func (*Jails) Descriptor() ([]byte, []int) {
	return file_internal_agent_agent_proto_rawDescGZIP(), []int{1}
}

func (x *Jails) GetJails() []*Jail {
	if x != nil {
		return x.Jails
	}
	return nil
}

type Jail struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name     string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Hostname string `protobuf:"bytes,2,opt,name=hostname,proto3" json:"hostname,omitempty"`
	Path     string `protobuf:"bytes,3,opt,name=path,proto3" json:"path,omitempty"`
}

func (x *Jail) Reset() {
	*x = Jail{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_agent_agent_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Jail) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Jail) ProtoMessage() {}

func (x *Jail) ProtoReflect() protoreflect.Message {
	mi := &file_internal_agent_agent_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Jail.ProtoReflect.Descriptor instead.
func (*Jail) Descriptor() ([]byte, []int) {
	return file_internal_agent_agent_proto_rawDescGZIP(), []int{2}
}

func (x *Jail) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Jail) GetHostname() string {
	if x != nil {
		return x.Hostname
	}
	return ""
}

func (x *Jail) GetPath() string {
	if x != nil {
		return x.Path
	}
	return ""
}

type CommandResult struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Output   []byte `protobuf:"bytes,1,opt,name=output,proto3" json:"output,omitempty"`
	Error    []byte `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	ExitCode int32  `protobuf:"varint,3,opt,name=exit_code,json=exitCode,proto3" json:"exit_code,omitempty"`
}

func (x *CommandResult) Reset() {
	*x = CommandResult{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_agent_agent_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommandResult) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommandResult) ProtoMessage() {}

func (x *CommandResult) ProtoReflect() protoreflect.Message {
	mi := &file_internal_agent_agent_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommandResult.ProtoReflect.Descriptor instead.
func (*CommandResult) Descriptor() ([]byte, []int) {
	return file_internal_agent_agent_proto_rawDescGZIP(), []int{3}
}

func (x *CommandResult) GetOutput() []byte {
	if x != nil {
		return x.Output
	}
	return nil
}

func (x *CommandResult) GetError() []byte {
	if x != nil {
		return x.Error
	}
	return nil
}

func (x *CommandResult) GetExitCode() int32 {
	if x != nil {
		return x.ExitCode
	}
	return 0
}

type ProjectSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Name string `protobuf:"bytes,1,opt,name=name,proto3" json:"name,omitempty"`
	Url  string `protobuf:"bytes,2,opt,name=url,proto3" json:"url,omitempty"`
}

func (x *ProjectSpec) Reset() {
	*x = ProjectSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_agent_agent_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectSpec) ProtoMessage() {}

func (x *ProjectSpec) ProtoReflect() protoreflect.Message {
	mi := &file_internal_agent_agent_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectSpec.ProtoReflect.Descriptor instead.
func (*ProjectSpec) Descriptor() ([]byte, []int) {
	return file_internal_agent_agent_proto_rawDescGZIP(), []int{4}
}

func (x *ProjectSpec) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *ProjectSpec) GetUrl() string {
	if x != nil {
		return x.Url
	}
	return ""
}

type ProjectStatus struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	HeadCommit string `protobuf:"bytes,1,opt,name=head_commit,json=headCommit,proto3" json:"head_commit,omitempty"`
	LatestTag  string `protobuf:"bytes,2,opt,name=latest_tag,json=latestTag,proto3" json:"latest_tag,omitempty"`
}

func (x *ProjectStatus) Reset() {
	*x = ProjectStatus{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_agent_agent_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProjectStatus) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProjectStatus) ProtoMessage() {}

func (x *ProjectStatus) ProtoReflect() protoreflect.Message {
	mi := &file_internal_agent_agent_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProjectStatus.ProtoReflect.Descriptor instead.
func (*ProjectStatus) Descriptor() ([]byte, []int) {
	return file_internal_agent_agent_proto_rawDescGZIP(), []int{5}
}

func (x *ProjectStatus) GetHeadCommit() string {
	if x != nil {
		return x.HeadCommit
	}
	return ""
}

func (x *ProjectStatus) GetLatestTag() string {
	if x != nil {
		return x.LatestTag
	}
	return ""
}

type BuildSpec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Project *ProjectSpec `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	Commit  string       `protobuf:"bytes,2,opt,name=commit,proto3" json:"commit,omitempty"`
	Branch  string       `protobuf:"bytes,3,opt,name=branch,proto3" json:"branch,omitempty"`
	Tag     string       `protobuf:"bytes,4,opt,name=tag,proto3" json:"tag,omitempty"`
}

func (x *BuildSpec) Reset() {
	*x = BuildSpec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_internal_agent_agent_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *BuildSpec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*BuildSpec) ProtoMessage() {}

func (x *BuildSpec) ProtoReflect() protoreflect.Message {
	mi := &file_internal_agent_agent_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use BuildSpec.ProtoReflect.Descriptor instead.
func (*BuildSpec) Descriptor() ([]byte, []int) {
	return file_internal_agent_agent_proto_rawDescGZIP(), []int{6}
}

func (x *BuildSpec) GetProject() *ProjectSpec {
	if x != nil {
		return x.Project
	}
	return nil
}

func (x *BuildSpec) GetCommit() string {
	if x != nil {
		return x.Commit
	}
	return ""
}

func (x *BuildSpec) GetBranch() string {
	if x != nil {
		return x.Branch
	}
	return ""
}

func (x *BuildSpec) GetTag() string {
	if x != nil {
		return x.Tag
	}
	return ""
}

var File_internal_agent_agent_proto protoreflect.FileDescriptor

var file_internal_agent_agent_proto_rawDesc = []byte{
	0x0a, 0x1a, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2f, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x22, 0x07, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x2a, 0x0a, 0x05,
	0x4a, 0x61, 0x69, 0x6c, 0x73, 0x12, 0x21, 0x0a, 0x05, 0x6a, 0x61, 0x69, 0x6c, 0x73, 0x18, 0x01,
	0x20, 0x03, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x4a, 0x61, 0x69,
	0x6c, 0x52, 0x05, 0x6a, 0x61, 0x69, 0x6c, 0x73, 0x22, 0x4a, 0x0a, 0x04, 0x4a, 0x61, 0x69, 0x6c,
	0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x68, 0x6f, 0x73, 0x74, 0x6e, 0x61, 0x6d, 0x65,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x70, 0x61, 0x74, 0x68, 0x22, 0x5a, 0x0a, 0x0d, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52,
	0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x06, 0x6f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x05, 0x65, 0x72,
	0x72, 0x6f, 0x72, 0x12, 0x1b, 0x0a, 0x09, 0x65, 0x78, 0x69, 0x74, 0x5f, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x65, 0x78, 0x69, 0x74, 0x43, 0x6f, 0x64, 0x65,
	0x22, 0x33, 0x0a, 0x0b, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x53, 0x70, 0x65, 0x63, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x12, 0x10, 0x0a, 0x03, 0x75, 0x72, 0x6c, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x03, 0x75, 0x72, 0x6c, 0x22, 0x4f, 0x0a, 0x0d, 0x50, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74,
	0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x1f, 0x0a, 0x0b, 0x68, 0x65, 0x61, 0x64, 0x5f, 0x63,
	0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x68, 0x65, 0x61,
	0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x1d, 0x0a, 0x0a, 0x6c, 0x61, 0x74, 0x65, 0x73,
	0x74, 0x5f, 0x74, 0x61, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x6c, 0x61, 0x74,
	0x65, 0x73, 0x74, 0x54, 0x61, 0x67, 0x22, 0x7b, 0x0a, 0x09, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x53,
	0x70, 0x65, 0x63, 0x12, 0x2c, 0x0a, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x50, 0x72, 0x6f,
	0x6a, 0x65, 0x63, 0x74, 0x53, 0x70, 0x65, 0x63, 0x52, 0x07, 0x70, 0x72, 0x6f, 0x6a, 0x65, 0x63,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x63, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x16, 0x0a, 0x06, 0x62, 0x72, 0x61,
	0x6e, 0x63, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x62, 0x72, 0x61, 0x6e, 0x63,
	0x68, 0x12, 0x10, 0x0a, 0x03, 0x74, 0x61, 0x67, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03,
	0x74, 0x61, 0x67, 0x32, 0xb3, 0x01, 0x0a, 0x08, 0x4a, 0x61, 0x69, 0x6c, 0x48, 0x6f, 0x73, 0x74,
	0x12, 0x22, 0x0a, 0x04, 0x4c, 0x69, 0x73, 0x74, 0x12, 0x0c, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x4a,
	0x61, 0x69, 0x6c, 0x73, 0x12, 0x2c, 0x0a, 0x07, 0x52, 0x65, 0x73, 0x74, 0x61, 0x72, 0x74, 0x12,
	0x0b, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x4a, 0x61, 0x69, 0x6c, 0x1a, 0x14, 0x2e, 0x61,
	0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x75,
	0x6c, 0x74, 0x12, 0x2a, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x72, 0x74, 0x12, 0x0b, 0x2e, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2e, 0x4a, 0x61, 0x69, 0x6c, 0x1a, 0x14, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74,
	0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x29,
	0x0a, 0x04, 0x53, 0x74, 0x6f, 0x70, 0x12, 0x0b, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x4a,
	0x61, 0x69, 0x6c, 0x1a, 0x14, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6d, 0x6d,
	0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x72, 0x0a, 0x04, 0x4e, 0x6f, 0x64,
	0x65, 0x12, 0x34, 0x0a, 0x0e, 0x52, 0x75, 0x6e, 0x50, 0x75, 0x70, 0x70, 0x65, 0x74, 0x41, 0x67,
	0x65, 0x6e, 0x74, 0x12, 0x0c, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74,
	0x79, 0x1a, 0x14, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e,
	0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x34, 0x0a, 0x0e, 0x50, 0x61, 0x63, 0x6b, 0x61,
	0x67, 0x65, 0x55, 0x70, 0x67, 0x72, 0x61, 0x64, 0x65, 0x12, 0x0c, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x14, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e,
	0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x32, 0x72, 0x0a,
	0x05, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x12, 0x32, 0x0a, 0x08, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x54,
	0x61, 0x67, 0x12, 0x10, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64,
	0x53, 0x70, 0x65, 0x63, 0x1a, 0x14, 0x2e, 0x61, 0x67, 0x65, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6d,
	0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x12, 0x35, 0x0a, 0x0b, 0x42, 0x75,
	0x69, 0x6c, 0x64, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x12, 0x10, 0x2e, 0x61, 0x67, 0x65, 0x6e,
	0x74, 0x2e, 0x42, 0x75, 0x69, 0x6c, 0x64, 0x53, 0x70, 0x65, 0x63, 0x1a, 0x14, 0x2e, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x52, 0x65, 0x73, 0x75, 0x6c,
	0x74, 0x42, 0x10, 0x5a, 0x0e, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x61, 0x67,
	0x65, 0x6e, 0x74, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_internal_agent_agent_proto_rawDescOnce sync.Once
	file_internal_agent_agent_proto_rawDescData = file_internal_agent_agent_proto_rawDesc
)

func file_internal_agent_agent_proto_rawDescGZIP() []byte {
	file_internal_agent_agent_proto_rawDescOnce.Do(func() {
		file_internal_agent_agent_proto_rawDescData = protoimpl.X.CompressGZIP(file_internal_agent_agent_proto_rawDescData)
	})
	return file_internal_agent_agent_proto_rawDescData
}

var file_internal_agent_agent_proto_msgTypes = make([]protoimpl.MessageInfo, 7)
var file_internal_agent_agent_proto_goTypes = []interface{}{
	(*Empty)(nil),         // 0: agent.Empty
	(*Jails)(nil),         // 1: agent.Jails
	(*Jail)(nil),          // 2: agent.Jail
	(*CommandResult)(nil), // 3: agent.CommandResult
	(*ProjectSpec)(nil),   // 4: agent.ProjectSpec
	(*ProjectStatus)(nil), // 5: agent.ProjectStatus
	(*BuildSpec)(nil),     // 6: agent.BuildSpec
}
var file_internal_agent_agent_proto_depIdxs = []int32{
	2,  // 0: agent.Jails.jails:type_name -> agent.Jail
	4,  // 1: agent.BuildSpec.project:type_name -> agent.ProjectSpec
	0,  // 2: agent.JailHost.List:input_type -> agent.Empty
	2,  // 3: agent.JailHost.Restart:input_type -> agent.Jail
	2,  // 4: agent.JailHost.Start:input_type -> agent.Jail
	2,  // 5: agent.JailHost.Stop:input_type -> agent.Jail
	0,  // 6: agent.Node.RunPuppetAgent:input_type -> agent.Empty
	0,  // 7: agent.Node.PackageUpgrade:input_type -> agent.Empty
	6,  // 8: agent.Build.BuildTag:input_type -> agent.BuildSpec
	6,  // 9: agent.Build.BuildCommit:input_type -> agent.BuildSpec
	1,  // 10: agent.JailHost.List:output_type -> agent.Jails
	3,  // 11: agent.JailHost.Restart:output_type -> agent.CommandResult
	3,  // 12: agent.JailHost.Start:output_type -> agent.CommandResult
	3,  // 13: agent.JailHost.Stop:output_type -> agent.CommandResult
	3,  // 14: agent.Node.RunPuppetAgent:output_type -> agent.CommandResult
	3,  // 15: agent.Node.PackageUpgrade:output_type -> agent.CommandResult
	3,  // 16: agent.Build.BuildTag:output_type -> agent.CommandResult
	3,  // 17: agent.Build.BuildCommit:output_type -> agent.CommandResult
	10, // [10:18] is the sub-list for method output_type
	2,  // [2:10] is the sub-list for method input_type
	2,  // [2:2] is the sub-list for extension type_name
	2,  // [2:2] is the sub-list for extension extendee
	0,  // [0:2] is the sub-list for field type_name
}

func init() { file_internal_agent_agent_proto_init() }
func file_internal_agent_agent_proto_init() {
	if File_internal_agent_agent_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_internal_agent_agent_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
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
		file_internal_agent_agent_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Jails); i {
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
		file_internal_agent_agent_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Jail); i {
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
		file_internal_agent_agent_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommandResult); i {
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
		file_internal_agent_agent_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProjectSpec); i {
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
		file_internal_agent_agent_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProjectStatus); i {
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
		file_internal_agent_agent_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*BuildSpec); i {
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
			RawDescriptor: file_internal_agent_agent_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   7,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_internal_agent_agent_proto_goTypes,
		DependencyIndexes: file_internal_agent_agent_proto_depIdxs,
		MessageInfos:      file_internal_agent_agent_proto_msgTypes,
	}.Build()
	File_internal_agent_agent_proto = out.File
	file_internal_agent_agent_proto_rawDesc = nil
	file_internal_agent_agent_proto_goTypes = nil
	file_internal_agent_agent_proto_depIdxs = nil
}
