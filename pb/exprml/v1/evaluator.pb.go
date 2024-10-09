// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        (unknown)
// source: exprml/v1/evaluator.proto

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

// Status of the evaluation.
type EvaluateOutput_Status int32

const (
	// Evaluation was successful.
	EvaluateOutput_OK EvaluateOutput_Status = 0
	// Index is invalid.
	EvaluateOutput_INVALID_INDEX EvaluateOutput_Status = 100
	// Key is invalid.
	EvaluateOutput_INVALID_KEY EvaluateOutput_Status = 101
	// Type is unexpected.
	EvaluateOutput_UNEXPECTED_TYPE EvaluateOutput_Status = 102
	// Argument mismatch.
	EvaluateOutput_ARGUMENT_MISMATCH EvaluateOutput_Status = 103
	// Cases are not exhaustive.
	EvaluateOutput_CASES_NOT_EXHAUSTIVE EvaluateOutput_Status = 104
	// Reference not found.
	EvaluateOutput_REFERENCE_NOT_FOUND EvaluateOutput_Status = 105
	// Values are not comparable.
	EvaluateOutput_NOT_COMPARABLE EvaluateOutput_Status = 106
	// Not a finite number.
	EvaluateOutput_NOT_FINITE_NUMBER EvaluateOutput_Status = 107
	// Evaluation was aborted.
	EvaluateOutput_ABORTED EvaluateOutput_Status = 108
	// Unknown error.
	EvaluateOutput_UNKNOWN_ERROR EvaluateOutput_Status = 109
)

// Enum value maps for EvaluateOutput_Status.
var (
	EvaluateOutput_Status_name = map[int32]string{
		0:   "OK",
		100: "INVALID_INDEX",
		101: "INVALID_KEY",
		102: "UNEXPECTED_TYPE",
		103: "ARGUMENT_MISMATCH",
		104: "CASES_NOT_EXHAUSTIVE",
		105: "REFERENCE_NOT_FOUND",
		106: "NOT_COMPARABLE",
		107: "NOT_FINITE_NUMBER",
		108: "ABORTED",
		109: "UNKNOWN_ERROR",
	}
	EvaluateOutput_Status_value = map[string]int32{
		"OK":                   0,
		"INVALID_INDEX":        100,
		"INVALID_KEY":          101,
		"UNEXPECTED_TYPE":      102,
		"ARGUMENT_MISMATCH":    103,
		"CASES_NOT_EXHAUSTIVE": 104,
		"REFERENCE_NOT_FOUND":  105,
		"NOT_COMPARABLE":       106,
		"NOT_FINITE_NUMBER":    107,
		"ABORTED":              108,
		"UNKNOWN_ERROR":        109,
	}
)

func (x EvaluateOutput_Status) Enum() *EvaluateOutput_Status {
	p := new(EvaluateOutput_Status)
	*p = x
	return p
}

func (x EvaluateOutput_Status) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (EvaluateOutput_Status) Descriptor() protoreflect.EnumDescriptor {
	return file_exprml_v1_evaluator_proto_enumTypes[0].Descriptor()
}

func (EvaluateOutput_Status) Type() protoreflect.EnumType {
	return &file_exprml_v1_evaluator_proto_enumTypes[0]
}

func (x EvaluateOutput_Status) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use EvaluateOutput_Status.Descriptor instead.
func (EvaluateOutput_Status) EnumDescriptor() ([]byte, []int) {
	return file_exprml_v1_evaluator_proto_rawDescGZIP(), []int{2, 0}
}

// FunDefList is a list of function definitions.
type FunDefList struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Parent function definition list.
	Parent *FunDefList `protobuf:"bytes,1,opt,name=parent,proto3" json:"parent,omitempty"`
	// Function definitions.
	FunDef *Node `protobuf:"bytes,3,opt,name=fun_def,json=funDef,proto3" json:"fun_def,omitempty"`
}

func (x *FunDefList) Reset() {
	*x = FunDefList{}
	mi := &file_exprml_v1_evaluator_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *FunDefList) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*FunDefList) ProtoMessage() {}

func (x *FunDefList) ProtoReflect() protoreflect.Message {
	mi := &file_exprml_v1_evaluator_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use FunDefList.ProtoReflect.Descriptor instead.
func (*FunDefList) Descriptor() ([]byte, []int) {
	return file_exprml_v1_evaluator_proto_rawDescGZIP(), []int{0}
}

func (x *FunDefList) GetParent() *FunDefList {
	if x != nil {
		return x.Parent
	}
	return nil
}

func (x *FunDefList) GetFunDef() *Node {
	if x != nil {
		return x.FunDef
	}
	return nil
}

// EvaluateInput is the input message for the EvaluateExpr method.
type EvaluateInput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Function definition stack.
	DefStack *FunDefList `protobuf:"bytes,1,opt,name=def_stack,json=defStack,proto3" json:"def_stack,omitempty"`
	// Expression to evaluate.
	Expr *Node `protobuf:"bytes,2,opt,name=expr,proto3" json:"expr,omitempty"`
}

func (x *EvaluateInput) Reset() {
	*x = EvaluateInput{}
	mi := &file_exprml_v1_evaluator_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EvaluateInput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EvaluateInput) ProtoMessage() {}

func (x *EvaluateInput) ProtoReflect() protoreflect.Message {
	mi := &file_exprml_v1_evaluator_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EvaluateInput.ProtoReflect.Descriptor instead.
func (*EvaluateInput) Descriptor() ([]byte, []int) {
	return file_exprml_v1_evaluator_proto_rawDescGZIP(), []int{1}
}

func (x *EvaluateInput) GetDefStack() *FunDefList {
	if x != nil {
		return x.DefStack
	}
	return nil
}

func (x *EvaluateInput) GetExpr() *Node {
	if x != nil {
		return x.Expr
	}
	return nil
}

// EvaluateOutput is the output message for the EvaluateExpr method.
type EvaluateOutput struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// Status of the evaluation.
	Status EvaluateOutput_Status `protobuf:"varint,1,opt,name=status,proto3,enum=exprml.v1.EvaluateOutput_Status" json:"status,omitempty"`
	// Error message if status is not OK.
	ErrorMessage string `protobuf:"bytes,2,opt,name=error_message,json=errorMessage,proto3" json:"error_message,omitempty"`
	// Error path if status is not OK.
	ErrorPath *Node_Path `protobuf:"bytes,3,opt,name=error_path,json=errorPath,proto3" json:"error_path,omitempty"`
	// Result of the evaluation.
	Value *Value `protobuf:"bytes,4,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *EvaluateOutput) Reset() {
	*x = EvaluateOutput{}
	mi := &file_exprml_v1_evaluator_proto_msgTypes[2]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *EvaluateOutput) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EvaluateOutput) ProtoMessage() {}

func (x *EvaluateOutput) ProtoReflect() protoreflect.Message {
	mi := &file_exprml_v1_evaluator_proto_msgTypes[2]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EvaluateOutput.ProtoReflect.Descriptor instead.
func (*EvaluateOutput) Descriptor() ([]byte, []int) {
	return file_exprml_v1_evaluator_proto_rawDescGZIP(), []int{2}
}

func (x *EvaluateOutput) GetStatus() EvaluateOutput_Status {
	if x != nil {
		return x.Status
	}
	return EvaluateOutput_OK
}

func (x *EvaluateOutput) GetErrorMessage() string {
	if x != nil {
		return x.ErrorMessage
	}
	return ""
}

func (x *EvaluateOutput) GetErrorPath() *Node_Path {
	if x != nil {
		return x.ErrorPath
	}
	return nil
}

func (x *EvaluateOutput) GetValue() *Value {
	if x != nil {
		return x.Value
	}
	return nil
}

var File_exprml_v1_evaluator_proto protoreflect.FileDescriptor

var file_exprml_v1_evaluator_proto_rawDesc = []byte{
	0x0a, 0x19, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x65, 0x76, 0x61, 0x6c,
	0x75, 0x61, 0x74, 0x6f, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x65, 0x78, 0x70,
	0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x1a, 0x14, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2f, 0x76,
	0x31, 0x2f, 0x6e, 0x6f, 0x64, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x65, 0x78,
	0x70, 0x72, 0x6d, 0x6c, 0x2f, 0x76, 0x31, 0x2f, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x22, 0x65, 0x0a, 0x0a, 0x46, 0x75, 0x6e, 0x44, 0x65, 0x66, 0x4c, 0x69, 0x73,
	0x74, 0x12, 0x2d, 0x0a, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x15, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x75,
	0x6e, 0x44, 0x65, 0x66, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x06, 0x70, 0x61, 0x72, 0x65, 0x6e, 0x74,
	0x12, 0x28, 0x0a, 0x07, 0x66, 0x75, 0x6e, 0x5f, 0x64, 0x65, 0x66, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0f, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f,
	0x64, 0x65, 0x52, 0x06, 0x66, 0x75, 0x6e, 0x44, 0x65, 0x66, 0x22, 0x68, 0x0a, 0x0d, 0x45, 0x76,
	0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x12, 0x32, 0x0a, 0x09, 0x64,
	0x65, 0x66, 0x5f, 0x73, 0x74, 0x61, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15,
	0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x46, 0x75, 0x6e, 0x44, 0x65,
	0x66, 0x4c, 0x69, 0x73, 0x74, 0x52, 0x08, 0x64, 0x65, 0x66, 0x53, 0x74, 0x61, 0x63, 0x6b, 0x12,
	0x23, 0x0a, 0x04, 0x65, 0x78, 0x70, 0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0f, 0x2e,
	0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x52, 0x04,
	0x65, 0x78, 0x70, 0x72, 0x22, 0xad, 0x03, 0x0a, 0x0e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74,
	0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x12, 0x38, 0x0a, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x20, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x2e, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x52, 0x06, 0x73, 0x74, 0x61, 0x74, 0x75,
	0x73, 0x12, 0x23, 0x0a, 0x0d, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f, 0x6d, 0x65, 0x73, 0x73, 0x61,
	0x67, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0c, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x4d,
	0x65, 0x73, 0x73, 0x61, 0x67, 0x65, 0x12, 0x33, 0x0a, 0x0a, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x5f,
	0x70, 0x61, 0x74, 0x68, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x14, 0x2e, 0x65, 0x78, 0x70,
	0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x4e, 0x6f, 0x64, 0x65, 0x2e, 0x50, 0x61, 0x74, 0x68,
	0x52, 0x09, 0x65, 0x72, 0x72, 0x6f, 0x72, 0x50, 0x61, 0x74, 0x68, 0x12, 0x26, 0x0a, 0x05, 0x76,
	0x61, 0x6c, 0x75, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x65, 0x78, 0x70,
	0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x22, 0xde, 0x01, 0x0a, 0x06, 0x53, 0x74, 0x61, 0x74, 0x75, 0x73, 0x12, 0x06,
	0x0a, 0x02, 0x4f, 0x4b, 0x10, 0x00, 0x12, 0x11, 0x0a, 0x0d, 0x49, 0x4e, 0x56, 0x41, 0x4c, 0x49,
	0x44, 0x5f, 0x49, 0x4e, 0x44, 0x45, 0x58, 0x10, 0x64, 0x12, 0x0f, 0x0a, 0x0b, 0x49, 0x4e, 0x56,
	0x41, 0x4c, 0x49, 0x44, 0x5f, 0x4b, 0x45, 0x59, 0x10, 0x65, 0x12, 0x13, 0x0a, 0x0f, 0x55, 0x4e,
	0x45, 0x58, 0x50, 0x45, 0x43, 0x54, 0x45, 0x44, 0x5f, 0x54, 0x59, 0x50, 0x45, 0x10, 0x66, 0x12,
	0x15, 0x0a, 0x11, 0x41, 0x52, 0x47, 0x55, 0x4d, 0x45, 0x4e, 0x54, 0x5f, 0x4d, 0x49, 0x53, 0x4d,
	0x41, 0x54, 0x43, 0x48, 0x10, 0x67, 0x12, 0x18, 0x0a, 0x14, 0x43, 0x41, 0x53, 0x45, 0x53, 0x5f,
	0x4e, 0x4f, 0x54, 0x5f, 0x45, 0x58, 0x48, 0x41, 0x55, 0x53, 0x54, 0x49, 0x56, 0x45, 0x10, 0x68,
	0x12, 0x17, 0x0a, 0x13, 0x52, 0x45, 0x46, 0x45, 0x52, 0x45, 0x4e, 0x43, 0x45, 0x5f, 0x4e, 0x4f,
	0x54, 0x5f, 0x46, 0x4f, 0x55, 0x4e, 0x44, 0x10, 0x69, 0x12, 0x12, 0x0a, 0x0e, 0x4e, 0x4f, 0x54,
	0x5f, 0x43, 0x4f, 0x4d, 0x50, 0x41, 0x52, 0x41, 0x42, 0x4c, 0x45, 0x10, 0x6a, 0x12, 0x15, 0x0a,
	0x11, 0x4e, 0x4f, 0x54, 0x5f, 0x46, 0x49, 0x4e, 0x49, 0x54, 0x45, 0x5f, 0x4e, 0x55, 0x4d, 0x42,
	0x45, 0x52, 0x10, 0x6b, 0x12, 0x0b, 0x0a, 0x07, 0x41, 0x42, 0x4f, 0x52, 0x54, 0x45, 0x44, 0x10,
	0x6c, 0x12, 0x11, 0x0a, 0x0d, 0x55, 0x4e, 0x4b, 0x4e, 0x4f, 0x57, 0x4e, 0x5f, 0x45, 0x52, 0x52,
	0x4f, 0x52, 0x10, 0x6d, 0x32, 0xbf, 0x07, 0x0a, 0x09, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74,
	0x6f, 0x72, 0x12, 0x45, 0x0a, 0x0c, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x45, 0x78,
	0x70, 0x72, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x19, 0x2e, 0x65,
	0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74,
	0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x0c, 0x45, 0x76, 0x61,
	0x6c, 0x75, 0x61, 0x74, 0x65, 0x45, 0x76, 0x61, 0x6c, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x70, 0x72,
	0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x49, 0x6e,
	0x70, 0x75, 0x74, 0x1a, 0x19, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e,
	0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00,
	0x12, 0x47, 0x0a, 0x0e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x53, 0x63, 0x61, 0x6c,
	0x61, 0x72, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x19, 0x2e, 0x65,
	0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74,
	0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12, 0x44, 0x0a, 0x0b, 0x45, 0x76, 0x61,
	0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x62, 0x6a, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d,
	0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x70,
	0x75, 0x74, 0x1a, 0x19, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12,
	0x44, 0x0a, 0x0b, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x41, 0x72, 0x72, 0x12, 0x18,
	0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75,
	0x61, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x19, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d,
	0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x75, 0x74,
	0x70, 0x75, 0x74, 0x22, 0x00, 0x12, 0x45, 0x0a, 0x0c, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74,
	0x65, 0x4a, 0x73, 0x6f, 0x6e, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76,
	0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a,
	0x19, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c,
	0x75, 0x61, 0x74, 0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12, 0x4a, 0x0a, 0x11,
	0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x52, 0x61, 0x6e, 0x67, 0x65, 0x49, 0x74, 0x65,
	0x72, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76,
	0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x19, 0x2e, 0x65, 0x78,
	0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65,
	0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x0f, 0x45, 0x76, 0x61, 0x6c,
	0x75, 0x61, 0x74, 0x65, 0x47, 0x65, 0x74, 0x45, 0x6c, 0x65, 0x6d, 0x12, 0x18, 0x2e, 0x65, 0x78,
	0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65,
	0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x19, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76,
	0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74,
	0x22, 0x00, 0x12, 0x48, 0x0a, 0x0f, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x46, 0x75,
	0x6e, 0x43, 0x61, 0x6c, 0x6c, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76,
	0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a,
	0x19, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c,
	0x75, 0x61, 0x74, 0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12, 0x46, 0x0a, 0x0d,
	0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x43, 0x61, 0x73, 0x65, 0x73, 0x12, 0x18, 0x2e,
	0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61,
	0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x19, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x75, 0x74, 0x70,
	0x75, 0x74, 0x22, 0x00, 0x12, 0x48, 0x0a, 0x0f, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65,
	0x4f, 0x70, 0x55, 0x6e, 0x61, 0x72, 0x79, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c,
	0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75,
	0x74, 0x1a, 0x19, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76,
	0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12, 0x49,
	0x0a, 0x10, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x70, 0x42, 0x69, 0x6e, 0x61,
	0x72, 0x79, 0x12, 0x18, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45,
	0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x19, 0x2e, 0x65,
	0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74,
	0x65, 0x4f, 0x75, 0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x12, 0x4b, 0x0a, 0x12, 0x45, 0x76, 0x61,
	0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x70, 0x56, 0x61, 0x72, 0x69, 0x61, 0x64, 0x69, 0x63, 0x12,
	0x18, 0x2e, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c,
	0x75, 0x61, 0x74, 0x65, 0x49, 0x6e, 0x70, 0x75, 0x74, 0x1a, 0x19, 0x2e, 0x65, 0x78, 0x70, 0x72,
	0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x2e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61, 0x74, 0x65, 0x4f, 0x75,
	0x74, 0x70, 0x75, 0x74, 0x22, 0x00, 0x42, 0x97, 0x01, 0x0a, 0x0d, 0x63, 0x6f, 0x6d, 0x2e, 0x65,
	0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x76, 0x31, 0x42, 0x0e, 0x45, 0x76, 0x61, 0x6c, 0x75, 0x61,
	0x74, 0x6f, 0x72, 0x50, 0x72, 0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x31, 0x67, 0x69, 0x74, 0x68,
	0x75, 0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2f, 0x65, 0x78,
	0x70, 0x72, 0x6d, 0x6c, 0x2d, 0x67, 0x6f, 0x2f, 0x70, 0x62, 0x2f, 0x65, 0x78, 0x70, 0x72, 0x6d,
	0x6c, 0x2f, 0x76, 0x31, 0x3b, 0x65, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x76, 0x31, 0xa2, 0x02, 0x03,
	0x45, 0x58, 0x58, 0xaa, 0x02, 0x09, 0x45, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x2e, 0x56, 0x31, 0xca,
	0x02, 0x09, 0x45, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x5c, 0x56, 0x31, 0xe2, 0x02, 0x15, 0x45, 0x78,
	0x70, 0x72, 0x6d, 0x6c, 0x5c, 0x56, 0x31, 0x5c, 0x47, 0x50, 0x42, 0x4d, 0x65, 0x74, 0x61, 0x64,
	0x61, 0x74, 0x61, 0xea, 0x02, 0x0a, 0x45, 0x78, 0x70, 0x72, 0x6d, 0x6c, 0x3a, 0x3a, 0x56, 0x31,
	0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_exprml_v1_evaluator_proto_rawDescOnce sync.Once
	file_exprml_v1_evaluator_proto_rawDescData = file_exprml_v1_evaluator_proto_rawDesc
)

func file_exprml_v1_evaluator_proto_rawDescGZIP() []byte {
	file_exprml_v1_evaluator_proto_rawDescOnce.Do(func() {
		file_exprml_v1_evaluator_proto_rawDescData = protoimpl.X.CompressGZIP(file_exprml_v1_evaluator_proto_rawDescData)
	})
	return file_exprml_v1_evaluator_proto_rawDescData
}

var file_exprml_v1_evaluator_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_exprml_v1_evaluator_proto_msgTypes = make([]protoimpl.MessageInfo, 3)
var file_exprml_v1_evaluator_proto_goTypes = []any{
	(EvaluateOutput_Status)(0), // 0: exprml.v1.EvaluateOutput.Status
	(*FunDefList)(nil),         // 1: exprml.v1.FunDefList
	(*EvaluateInput)(nil),      // 2: exprml.v1.EvaluateInput
	(*EvaluateOutput)(nil),     // 3: exprml.v1.EvaluateOutput
	(*Node)(nil),               // 4: exprml.v1.Node
	(*Node_Path)(nil),          // 5: exprml.v1.Node.Path
	(*Value)(nil),              // 6: exprml.v1.Value
}
var file_exprml_v1_evaluator_proto_depIdxs = []int32{
	1,  // 0: exprml.v1.FunDefList.parent:type_name -> exprml.v1.FunDefList
	4,  // 1: exprml.v1.FunDefList.fun_def:type_name -> exprml.v1.Node
	1,  // 2: exprml.v1.EvaluateInput.def_stack:type_name -> exprml.v1.FunDefList
	4,  // 3: exprml.v1.EvaluateInput.expr:type_name -> exprml.v1.Node
	0,  // 4: exprml.v1.EvaluateOutput.status:type_name -> exprml.v1.EvaluateOutput.Status
	5,  // 5: exprml.v1.EvaluateOutput.error_path:type_name -> exprml.v1.Node.Path
	6,  // 6: exprml.v1.EvaluateOutput.value:type_name -> exprml.v1.Value
	2,  // 7: exprml.v1.Evaluator.EvaluateExpr:input_type -> exprml.v1.EvaluateInput
	2,  // 8: exprml.v1.Evaluator.EvaluateEval:input_type -> exprml.v1.EvaluateInput
	2,  // 9: exprml.v1.Evaluator.EvaluateScalar:input_type -> exprml.v1.EvaluateInput
	2,  // 10: exprml.v1.Evaluator.EvaluateObj:input_type -> exprml.v1.EvaluateInput
	2,  // 11: exprml.v1.Evaluator.EvaluateArr:input_type -> exprml.v1.EvaluateInput
	2,  // 12: exprml.v1.Evaluator.EvaluateJson:input_type -> exprml.v1.EvaluateInput
	2,  // 13: exprml.v1.Evaluator.EvaluateRangeIter:input_type -> exprml.v1.EvaluateInput
	2,  // 14: exprml.v1.Evaluator.EvaluateGetElem:input_type -> exprml.v1.EvaluateInput
	2,  // 15: exprml.v1.Evaluator.EvaluateFunCall:input_type -> exprml.v1.EvaluateInput
	2,  // 16: exprml.v1.Evaluator.EvaluateCases:input_type -> exprml.v1.EvaluateInput
	2,  // 17: exprml.v1.Evaluator.EvaluateOpUnary:input_type -> exprml.v1.EvaluateInput
	2,  // 18: exprml.v1.Evaluator.EvaluateOpBinary:input_type -> exprml.v1.EvaluateInput
	2,  // 19: exprml.v1.Evaluator.EvaluateOpVariadic:input_type -> exprml.v1.EvaluateInput
	3,  // 20: exprml.v1.Evaluator.EvaluateExpr:output_type -> exprml.v1.EvaluateOutput
	3,  // 21: exprml.v1.Evaluator.EvaluateEval:output_type -> exprml.v1.EvaluateOutput
	3,  // 22: exprml.v1.Evaluator.EvaluateScalar:output_type -> exprml.v1.EvaluateOutput
	3,  // 23: exprml.v1.Evaluator.EvaluateObj:output_type -> exprml.v1.EvaluateOutput
	3,  // 24: exprml.v1.Evaluator.EvaluateArr:output_type -> exprml.v1.EvaluateOutput
	3,  // 25: exprml.v1.Evaluator.EvaluateJson:output_type -> exprml.v1.EvaluateOutput
	3,  // 26: exprml.v1.Evaluator.EvaluateRangeIter:output_type -> exprml.v1.EvaluateOutput
	3,  // 27: exprml.v1.Evaluator.EvaluateGetElem:output_type -> exprml.v1.EvaluateOutput
	3,  // 28: exprml.v1.Evaluator.EvaluateFunCall:output_type -> exprml.v1.EvaluateOutput
	3,  // 29: exprml.v1.Evaluator.EvaluateCases:output_type -> exprml.v1.EvaluateOutput
	3,  // 30: exprml.v1.Evaluator.EvaluateOpUnary:output_type -> exprml.v1.EvaluateOutput
	3,  // 31: exprml.v1.Evaluator.EvaluateOpBinary:output_type -> exprml.v1.EvaluateOutput
	3,  // 32: exprml.v1.Evaluator.EvaluateOpVariadic:output_type -> exprml.v1.EvaluateOutput
	20, // [20:33] is the sub-list for method output_type
	7,  // [7:20] is the sub-list for method input_type
	7,  // [7:7] is the sub-list for extension type_name
	7,  // [7:7] is the sub-list for extension extendee
	0,  // [0:7] is the sub-list for field type_name
}

func init() { file_exprml_v1_evaluator_proto_init() }
func file_exprml_v1_evaluator_proto_init() {
	if File_exprml_v1_evaluator_proto != nil {
		return
	}
	file_exprml_v1_node_proto_init()
	file_exprml_v1_value_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_exprml_v1_evaluator_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   3,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_exprml_v1_evaluator_proto_goTypes,
		DependencyIndexes: file_exprml_v1_evaluator_proto_depIdxs,
		EnumInfos:         file_exprml_v1_evaluator_proto_enumTypes,
		MessageInfos:      file_exprml_v1_evaluator_proto_msgTypes,
	}.Build()
	File_exprml_v1_evaluator_proto = out.File
	file_exprml_v1_evaluator_proto_rawDesc = nil
	file_exprml_v1_evaluator_proto_goTypes = nil
	file_exprml_v1_evaluator_proto_depIdxs = nil
}
