// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: exprml/v1/value.proto

package exprmlv1

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

// Type of a JSON value.
type Value_Type int32

const (
	// Unspecified.
	Value_UNSPECIFIED Value_Type = 0
	// Null type.
	Value_NULL Value_Type = 1
	// Boolean type.
	Value_BOOL Value_Type = 2
	// Number type.
	Value_NUM Value_Type = 3
	// String type.
	Value_STR Value_Type = 4
	// Array type.
	Value_ARR Value_Type = 5
	// Object type.
	Value_OBJ Value_Type = 6
)

// Enum value maps for Value_Type.
var (
	Value_Type_name = map[int32]string{
		0: "UNSPECIFIED",
		1: "NULL",
		2: "BOOL",
		3: "NUM",
		4: "STR",
		5: "ARR",
		6: "OBJ",
	}
	Value_Type_value = map[string]int32{
		"UNSPECIFIED": 0,
		"NULL":        1,
		"BOOL":        2,
		"NUM":         3,
		"STR":         4,
		"ARR":         5,
		"OBJ":         6,
	}
)

func (x Value_Type) Enum() *Value_Type {
	p := new(Value_Type)
	*p = x
	return p
}

func (x Value_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Value_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_exprml_v1_value_proto_enumTypes[0].Descriptor()
}

func (Value_Type) Type() protoreflect.EnumType {
	return &file_exprml_v1_value_proto_enumTypes[0]
}

func (x Value_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Value_Type.Descriptor instead.
func (Value_Type) EnumDescriptor() ([]byte, []int) {
	return file_exprml_v1_value_proto_rawDescGZIP(), []int{0, 0}
}

// JSON value.
type Value struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Type of the value.
	Type Value_Type `protobuf:"varint,1,opt,name=type,proto3,enum=exprml.v1.Value_Type" json:"type,omitempty"`
	// bool has a boolean value if the type is TYPE_BOOL.
	Bool bool `protobuf:"varint,2,opt,name=bool,proto3" json:"bool,omitempty"`
	// num has a number value if the type is TYPE_NUM.
	Num float64 `protobuf:"fixed64,3,opt,name=num,proto3" json:"num,omitempty"`
	// str has a string value if the type is TYPE_STR.
	Str string `protobuf:"bytes,4,opt,name=str,proto3" json:"str,omitempty"`
	// arr has an array value if the type is TYPE_ARR.
	Arr []*Value `protobuf:"bytes,5,rep,name=arr,proto3" json:"arr,omitempty"`
	// obj has an object value if the type is TYPE_OBJ.
	Obj map[string]*Value `protobuf:"bytes,6,rep,name=obj,proto3" json:"obj,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *Value) Reset() {
	*x = Value{}
	mi := &file_exprml_v1_value_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Value) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Value) ProtoMessage() {}

func (x *Value) ProtoReflect() protoreflect.Message {
	mi := &file_exprml_v1_value_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Value.ProtoReflect.Descriptor instead.
func (*Value) Descriptor() ([]byte, []int) {
	return file_exprml_v1_value_proto_rawDescGZIP(), []int{0}
}

func (x *Value) GetType() Value_Type {
	if x != nil {
		return x.Type
	}
	return Value_UNSPECIFIED
}

func (x *Value) GetBool() bool {
	if x != nil {
		return x.Bool
	}
	return false
}

func (x *Value) GetNum() float64 {
	if x != nil {
		return x.Num
	}
	return 0
}

func (x *Value) GetStr() string {
	if x != nil {
		return x.Str
	}
	return ""
}

func (x *Value) GetArr() []*Value {
	if x != nil {
		return x.Arr
	}
	return nil
}

func (x *Value) GetObj() map[string]*Value {
	if x != nil {
		return x.Obj
	}
	return nil
}

var File_exprml_v1_value_proto protoreflect.FileDescriptor

var file_exprml_v1_value_proto_rawDesc = []byte{
	0x0a, 0x15, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x61, 0x6c, 0x75,
	0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e,
	0x76, 0x31, 0x22, 0xd6, 0x02, 0x0a, 0x05, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x12, 0x29, 0x0a, 0x04,
	0x74, 0x79, 0x70, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x15, 0x2e, 0x65, 0x78, 0x70,
	0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x54, 0x79, 0x70,
	0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x62, 0x6f, 0x6f, 0x6c, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x08, 0x52, 0x04, 0x62, 0x6f, 0x6f, 0x6c, 0x12, 0x10, 0x0a, 0x03, 0x6e,
	0x75, 0x6d, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6e, 0x75, 0x6d, 0x12, 0x10, 0x0a,
	0x03, 0x73, 0x74, 0x72, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x73, 0x74, 0x72, 0x12,
	0x22, 0x0a, 0x03, 0x61, 0x72, 0x72, 0x18, 0x05, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65,
	0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x03,
	0x61, 0x72, 0x72, 0x12, 0x2b, 0x0a, 0x03, 0x6f, 0x62, 0x6a, 0x18, 0x06, 0x20, 0x03, 0x28, 0x0b,
	0x32, 0x19, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c,
	0x75, 0x65, 0x2e, 0x4f, 0x62, 0x6a, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x52, 0x03, 0x6f, 0x62, 0x6a,
	0x1a, 0x48, 0x0a, 0x08, 0x4f, 0x62, 0x6a, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03,
	0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x26,
	0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e,
	0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x22, 0x4f, 0x0a, 0x04, 0x54, 0x79,
	0x70, 0x65, 0x12, 0x0f, 0x0a, 0x0b, 0x55, 0x4e, 0x53, 0x50, 0x45, 0x43, 0x49, 0x46, 0x49, 0x45,
	0x44, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x4e, 0x55, 0x4c, 0x4c, 0x10, 0x01, 0x12, 0x08, 0x0a,
	0x04, 0x42, 0x4f, 0x4f, 0x4c, 0x10, 0x02, 0x12, 0x07, 0x0a, 0x03, 0x4e, 0x55, 0x4d, 0x10, 0x03,
	0x12, 0x07, 0x0a, 0x03, 0x53, 0x54, 0x52, 0x10, 0x04, 0x12, 0x07, 0x0a, 0x03, 0x41, 0x52, 0x52,
	0x10, 0x05, 0x12, 0x07, 0x0a, 0x03, 0x4f, 0x42, 0x4a, 0x10, 0x06, 0x42, 0x93, 0x01, 0x0a, 0x0d,
	0x63, 0x6f, 0x6d, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x42, 0x0a, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x31, 0x67, 0x69, 0x74,
	0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2f, 0x65,
	0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x62, 0x2f, 0x65, 0x78, 0x70, 0x72,
	0x6d, 0x6c, 0x2f, 0x76, 0x31, 0x3b, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x76, 0x31, 0xa2, 0x02,
	0x03, 0x45, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x45, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x56, 0x31,
	0xca, 0x02, 0x09, 0x45, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x45,
	0x78, 0x70, 0x72, 0x6d, 0x6c, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x45, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x3a, 0x3a, 0x56,
	0x31, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_exprml_v1_value_proto_rawDescOnce sync.Once
	file_exprml_v1_value_proto_rawDescData = file_exprml_v1_value_proto_rawDesc
)

func file_exprml_v1_value_proto_rawDescGZIP() []byte {
	file_exprml_v1_value_proto_rawDescOnce.Do(func() {
		file_exprml_v1_value_proto_rawDescData = protoimpl.X.CompressGZIP(file_exprml_v1_value_proto_rawDescData)
	})
	return file_exprml_v1_value_proto_rawDescData
}

var file_exprml_v1_value_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_exprml_v1_value_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_exprml_v1_value_proto_goTypes = []any{
	(Value_Type)(0), // 0: exprml.v1.Value.Type
	(*Value)(nil),   // 1: exprml.v1.Value
	nil,             // 2: exprml.v1.Value.ObjEntry
}
var file_exprml_v1_value_proto_depIdxs = []int32{
	0, // 0: exprml.v1.Value.type:type_name -> exprml.v1.Value.Type
	1, // 1: exprml.v1.Value.arr:type_name -> exprml.v1.Value
	2, // 2: exprml.v1.Value.obj:type_name -> exprml.v1.Value.ObjEntry
	1, // 3: exprml.v1.Value.ObjEntry.value:type_name -> exprml.v1.Value
	4, // [4:4] is the sub-list for method output_type
	4, // [4:4] is the sub-list for method input_type
	4, // [4:4] is the sub-list for extension type_name
	4, // [4:4] is the sub-list for extension extendee
	0, // [0:4] is the sub-list for field type_name
}

func init() { file_exprml_v1_value_proto_init() }
func file_exprml_v1_value_proto_init() {
	if File_exprml_v1_value_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_exprml_v1_value_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_exprml_v1_value_proto_goTypes,
		DependencyIndexes: file_exprml_v1_value_proto_depIdxs,
		EnumInfos:         file_exprml_v1_value_proto_enumTypes,
		MessageInfos:      file_exprml_v1_value_proto_msgTypes,
	}.Build()
	File_exprml_v1_value_proto = out.File
	file_exprml_v1_value_proto_rawDesc = nil
	file_exprml_v1_value_proto_goTypes = nil
	file_exprml_v1_value_proto_depIdxs = nil
}
