// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.17.3
// source: modules/telemetry/telemetry.proto

package telemetry

import (
	reflect "reflect"

	proto "github.com/golang/protobuf/proto"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"

	inventory "github.com/zachfi/znet/modules/inventory"
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

var File_modules_telemetry_telemetry_proto protoreflect.FileDescriptor

var file_modules_telemetry_telemetry_proto_rawDesc = []byte{
	0x0a, 0x21, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65,
	0x74, 0x72, 0x79, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x12, 0x09, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x1a, 0x21,
	0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72,
	0x79, 0x2f, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x32, 0x81, 0x01, 0x0a, 0x09, 0x54, 0x65, 0x6c, 0x65, 0x6d, 0x65, 0x74, 0x72, 0x79, 0x12,
	0x39, 0x0a, 0x0f, 0x52, 0x65, 0x70, 0x6f, 0x72, 0x74, 0x49, 0x4f, 0x54, 0x44, 0x65, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x14, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x49,
	0x4f, 0x54, 0x44, 0x65, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x10, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e,
	0x74, 0x6f, 0x72, 0x79, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x39, 0x0a, 0x0f, 0x52, 0x65,
	0x70, 0x6f, 0x72, 0x74, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x49, 0x44, 0x12, 0x14, 0x2e,
	0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72,
	0x6b, 0x49, 0x44, 0x1a, 0x10, 0x2e, 0x69, 0x6e, 0x76, 0x65, 0x6e, 0x74, 0x6f, 0x72, 0x79, 0x2e,
	0x45, 0x6d, 0x70, 0x74, 0x79, 0x42, 0x2c, 0x5a, 0x2a, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e,
	0x63, 0x6f, 0x6d, 0x2f, 0x78, 0x61, 0x71, 0x75, 0x65, 0x32, 0x30, 0x38, 0x2f, 0x7a, 0x6e, 0x65,
	0x74, 0x2f, 0x6d, 0x6f, 0x64, 0x75, 0x6c, 0x65, 0x73, 0x2f, 0x74, 0x65, 0x6c, 0x65, 0x6d, 0x65,
	0x74, 0x72, 0x79, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var file_modules_telemetry_telemetry_proto_goTypes = []interface{}{
	(*inventory.IOTDevice)(nil), // 0: inventory.IOTDevice
	(*inventory.NetworkID)(nil), // 1: inventory.NetworkID
	(*inventory.Empty)(nil),     // 2: inventory.Empty
}
var file_modules_telemetry_telemetry_proto_depIdxs = []int32{
	0, // 0: telemetry.Telemetry.ReportIOTDevice:input_type -> inventory.IOTDevice
	1, // 1: telemetry.Telemetry.ReportNetworkID:input_type -> inventory.NetworkID
	2, // 2: telemetry.Telemetry.ReportIOTDevice:output_type -> inventory.Empty
	2, // 3: telemetry.Telemetry.ReportNetworkID:output_type -> inventory.Empty
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_modules_telemetry_telemetry_proto_init() }
func file_modules_telemetry_telemetry_proto_init() {
	if File_modules_telemetry_telemetry_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_modules_telemetry_telemetry_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_modules_telemetry_telemetry_proto_goTypes,
		DependencyIndexes: file_modules_telemetry_telemetry_proto_depIdxs,
	}.Build()
	File_modules_telemetry_telemetry_proto = out.File
	file_modules_telemetry_telemetry_proto_rawDesc = nil
	file_modules_telemetry_telemetry_proto_goTypes = nil
	file_modules_telemetry_telemetry_proto_depIdxs = nil
}
