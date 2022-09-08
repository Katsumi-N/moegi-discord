// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.21.3
// source: conoha.proto

package grpc

import (
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

type MinecraftRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Command string   `protobuf:"bytes,1,opt,name=command,proto3" json:"command,omitempty"`
	Args    []string `protobuf:"bytes,2,rep,name=args,proto3" json:"args,omitempty"`
}

func (x *MinecraftRequest) Reset() {
	*x = MinecraftRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conoha_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinecraftRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinecraftRequest) ProtoMessage() {}

func (x *MinecraftRequest) ProtoReflect() protoreflect.Message {
	mi := &file_conoha_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinecraftRequest.ProtoReflect.Descriptor instead.
func (*MinecraftRequest) Descriptor() ([]byte, []int) {
	return file_conoha_proto_rawDescGZIP(), []int{0}
}

func (x *MinecraftRequest) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

func (x *MinecraftRequest) GetArgs() []string {
	if x != nil {
		return x.Args
	}
	return nil
}

type MinecraftResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message  string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
	IsNormal bool   `protobuf:"varint,2,opt,name=is_normal,json=isNormal,proto3" json:"is_normal,omitempty"`
}

func (x *MinecraftResponse) Reset() {
	*x = MinecraftResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conoha_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *MinecraftResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*MinecraftResponse) ProtoMessage() {}

func (x *MinecraftResponse) ProtoReflect() protoreflect.Message {
	mi := &file_conoha_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use MinecraftResponse.ProtoReflect.Descriptor instead.
func (*MinecraftResponse) Descriptor() ([]byte, []int) {
	return file_conoha_proto_rawDescGZIP(), []int{1}
}

func (x *MinecraftResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

func (x *MinecraftResponse) GetIsNormal() bool {
	if x != nil {
		return x.IsNormal
	}
	return false
}

type VoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Command string   `protobuf:"bytes,1,opt,name=command,proto3" json:"command,omitempty"`
	Title   string   `protobuf:"bytes,2,opt,name=title,proto3" json:"title,omitempty"`
	Options []string `protobuf:"bytes,3,rep,name=options,proto3" json:"options,omitempty"`
}

func (x *VoteRequest) Reset() {
	*x = VoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conoha_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteRequest) ProtoMessage() {}

func (x *VoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_conoha_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteRequest.ProtoReflect.Descriptor instead.
func (*VoteRequest) Descriptor() ([]byte, []int) {
	return file_conoha_proto_rawDescGZIP(), []int{2}
}

func (x *VoteRequest) GetCommand() string {
	if x != nil {
		return x.Command
	}
	return ""
}

func (x *VoteRequest) GetTitle() string {
	if x != nil {
		return x.Title
	}
	return ""
}

func (x *VoteRequest) GetOptions() []string {
	if x != nil {
		return x.Options
	}
	return nil
}

type VoteResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Message string `protobuf:"bytes,1,opt,name=message,proto3" json:"message,omitempty"`
}

func (x *VoteResponse) Reset() {
	*x = VoteResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_conoha_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *VoteResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*VoteResponse) ProtoMessage() {}

func (x *VoteResponse) ProtoReflect() protoreflect.Message {
	mi := &file_conoha_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use VoteResponse.ProtoReflect.Descriptor instead.
func (*VoteResponse) Descriptor() ([]byte, []int) {
	return file_conoha_proto_rawDescGZIP(), []int{3}
}

func (x *VoteResponse) GetMessage() string {
	if x != nil {
		return x.Message
	}
	return ""
}

var File_conoha_proto protoreflect.FileDescriptor

var file_conoha_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x63, 0x6f, 0x6e, 0x6f, 0x68, 0x61, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x40,
	0x0a, 0x10, 0x4d, 0x69, 0x6e, 0x65, 0x63, 0x72, 0x61, 0x66, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65,
	0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x12, 0x0a, 0x04,
	0x61, 0x72, 0x67, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x09, 0x52, 0x04, 0x61, 0x72, 0x67, 0x73,
	0x22, 0x4a, 0x0a, 0x11, 0x4d, 0x69, 0x6e, 0x65, 0x63, 0x72, 0x61, 0x66, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12,
	0x1b, 0x0a, 0x09, 0x69, 0x73, 0x5f, 0x6e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x4e, 0x6f, 0x72, 0x6d, 0x61, 0x6c, 0x22, 0x57, 0x0a, 0x0b,
	0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f,
	0x6d, 0x6d, 0x61, 0x6e, 0x64, 0x12, 0x14, 0x0a, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x18, 0x02,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x74, 0x69, 0x74, 0x6c, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6f,
	0x70, 0x74, 0x69, 0x6f, 0x6e, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x09, 0x52, 0x07, 0x6f, 0x70,
	0x74, 0x69, 0x6f, 0x6e, 0x73, 0x22, 0x28, 0x0a, 0x0c, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x6d, 0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x32,
	0x6a, 0x0a, 0x0d, 0x43, 0x6f, 0x6e, 0x6f, 0x68, 0x61, 0x53, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65,
	0x12, 0x34, 0x0a, 0x09, 0x4d, 0x69, 0x6e, 0x65, 0x63, 0x72, 0x61, 0x66, 0x74, 0x12, 0x11, 0x2e,
	0x4d, 0x69, 0x6e, 0x65, 0x63, 0x72, 0x61, 0x66, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x1a, 0x12, 0x2e, 0x4d, 0x69, 0x6e, 0x65, 0x63, 0x72, 0x61, 0x66, 0x74, 0x52, 0x65, 0x73, 0x70,
	0x6f, 0x6e, 0x73, 0x65, 0x30, 0x01, 0x12, 0x23, 0x0a, 0x04, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x0c,
	0x2e, 0x56, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0d, 0x2e, 0x56,
	0x6f, 0x74, 0x65, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x70,
	0x6b, 0x67, 0x2f, 0x67, 0x72, 0x70, 0x63, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_conoha_proto_rawDescOnce sync.Once
	file_conoha_proto_rawDescData = file_conoha_proto_rawDesc
)

func file_conoha_proto_rawDescGZIP() []byte {
	file_conoha_proto_rawDescOnce.Do(func() {
		file_conoha_proto_rawDescData = protoimpl.X.CompressGZIP(file_conoha_proto_rawDescData)
	})
	return file_conoha_proto_rawDescData
}

var file_conoha_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_conoha_proto_goTypes = []interface{}{
	(*MinecraftRequest)(nil),  // 0: MinecraftRequest
	(*MinecraftResponse)(nil), // 1: MinecraftResponse
	(*VoteRequest)(nil),       // 2: VoteRequest
	(*VoteResponse)(nil),      // 3: VoteResponse
}
var file_conoha_proto_depIdxs = []int32{
	0, // 0: ConohaService.Minecraft:input_type -> MinecraftRequest
	2, // 1: ConohaService.Vote:input_type -> VoteRequest
	1, // 2: ConohaService.Minecraft:output_type -> MinecraftResponse
	3, // 3: ConohaService.Vote:output_type -> VoteResponse
	2, // [2:4] is the sub-list for method output_type
	0, // [0:2] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_conoha_proto_init() }
func file_conoha_proto_init() {
	if File_conoha_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_conoha_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinecraftRequest); i {
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
		file_conoha_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*MinecraftResponse); i {
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
		file_conoha_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteRequest); i {
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
		file_conoha_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*VoteResponse); i {
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
			RawDescriptor: file_conoha_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_conoha_proto_goTypes,
		DependencyIndexes: file_conoha_proto_depIdxs,
		MessageInfos:      file_conoha_proto_msgTypes,
	}.Build()
	File_conoha_proto = out.File
	file_conoha_proto_rawDesc = nil
	file_conoha_proto_goTypes = nil
	file_conoha_proto_depIdxs = nil
}
