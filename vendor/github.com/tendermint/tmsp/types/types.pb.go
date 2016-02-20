// Code generated by protoc-gen-go.
// source: types/types.proto
// DO NOT EDIT!

/*
Package types is a generated protocol buffer package.

It is generated from these files:
	types/types.proto

It has these top-level messages:
	Request
	Response
*/
package types

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type MessageType int32

const (
	MessageType_NullMessage MessageType = 0
	MessageType_Echo        MessageType = 1
	MessageType_Flush       MessageType = 2
	MessageType_Info        MessageType = 3
	MessageType_SetOption   MessageType = 4
	MessageType_Exception   MessageType = 5
	MessageType_AppendTx    MessageType = 17
	MessageType_CheckTx     MessageType = 18
	MessageType_Commit      MessageType = 19
	MessageType_Query       MessageType = 20
)

var MessageType_name = map[int32]string{
	0:  "NullMessage",
	1:  "Echo",
	2:  "Flush",
	3:  "Info",
	4:  "SetOption",
	5:  "Exception",
	17: "AppendTx",
	18: "CheckTx",
	19: "Commit",
	20: "Query",
}
var MessageType_value = map[string]int32{
	"NullMessage": 0,
	"Echo":        1,
	"Flush":       2,
	"Info":        3,
	"SetOption":   4,
	"Exception":   5,
	"AppendTx":    17,
	"CheckTx":     18,
	"Commit":      19,
	"Query":       20,
}

func (x MessageType) String() string {
	return proto.EnumName(MessageType_name, int32(x))
}
func (MessageType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type CodeType int32

const (
	CodeType_OK                CodeType = 0
	CodeType_InternalError     CodeType = 1
	CodeType_Unauthorized      CodeType = 2
	CodeType_InsufficientFees  CodeType = 3
	CodeType_UnknownRequest    CodeType = 4
	CodeType_EncodingError     CodeType = 5
	CodeType_BadNonce          CodeType = 6
	CodeType_UnknownAccount    CodeType = 7
	CodeType_InsufficientFunds CodeType = 8
)

var CodeType_name = map[int32]string{
	0: "OK",
	1: "InternalError",
	2: "Unauthorized",
	3: "InsufficientFees",
	4: "UnknownRequest",
	5: "EncodingError",
	6: "BadNonce",
	7: "UnknownAccount",
	8: "InsufficientFunds",
}
var CodeType_value = map[string]int32{
	"OK":                0,
	"InternalError":     1,
	"Unauthorized":      2,
	"InsufficientFees":  3,
	"UnknownRequest":    4,
	"EncodingError":     5,
	"BadNonce":          6,
	"UnknownAccount":    7,
	"InsufficientFunds": 8,
}

func (x CodeType) String() string {
	return proto.EnumName(CodeType_name, int32(x))
}
func (CodeType) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type Request struct {
	Type  MessageType `protobuf:"varint,1,opt,name=type,enum=types.MessageType" json:"type,omitempty"`
	Data  []byte      `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Key   string      `protobuf:"bytes,3,opt,name=key" json:"key,omitempty"`
	Value string      `protobuf:"bytes,4,opt,name=value" json:"value,omitempty"`
}

func (m *Request) Reset()                    { *m = Request{} }
func (m *Request) String() string            { return proto.CompactTextString(m) }
func (*Request) ProtoMessage()               {}
func (*Request) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type Response struct {
	Type  MessageType `protobuf:"varint,1,opt,name=type,enum=types.MessageType" json:"type,omitempty"`
	Data  []byte      `protobuf:"bytes,2,opt,name=data,proto3" json:"data,omitempty"`
	Code  CodeType    `protobuf:"varint,3,opt,name=code,enum=types.CodeType" json:"code,omitempty"`
	Error string      `protobuf:"bytes,4,opt,name=error" json:"error,omitempty"`
	Log   string      `protobuf:"bytes,5,opt,name=log" json:"log,omitempty"`
}

func (m *Response) Reset()                    { *m = Response{} }
func (m *Response) String() string            { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()               {}
func (*Response) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func init() {
	proto.RegisterType((*Request)(nil), "types.Request")
	proto.RegisterType((*Response)(nil), "types.Response")
	proto.RegisterEnum("types.MessageType", MessageType_name, MessageType_value)
	proto.RegisterEnum("types.CodeType", CodeType_name, CodeType_value)
}

var fileDescriptor0 = []byte{
	// 406 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0xa4, 0x92, 0x5f, 0x6e, 0xd3, 0x40,
	0x10, 0xc6, 0x71, 0x62, 0xe7, 0xcf, 0x24, 0x4d, 0x37, 0x43, 0x90, 0xfc, 0x88, 0x8a, 0x84, 0x50,
	0x1f, 0x8a, 0x54, 0x4e, 0x50, 0xa2, 0x54, 0x8a, 0x10, 0xad, 0x30, 0xed, 0x01, 0xcc, 0x7a, 0x12,
	0x5b, 0x71, 0x66, 0x8d, 0x77, 0x17, 0x1a, 0xee, 0xc0, 0x13, 0xe7, 0xe0, 0x8e, 0xec, 0xae, 0x53,
	0xa9, 0x3c, 0xf7, 0xc5, 0x9a, 0xef, 0xdb, 0x9d, 0x99, 0xdf, 0x8c, 0x17, 0xe6, 0xe6, 0xd0, 0x90,
	0x7e, 0x1f, 0xbe, 0x17, 0x4d, 0xab, 0x8c, 0xc2, 0x24, 0x88, 0xb3, 0x3d, 0x0c, 0x33, 0xfa, 0x6e,
	0x49, 0x1b, 0x7c, 0x0b, 0xb1, 0xf7, 0xd2, 0xe8, 0x75, 0xf4, 0x6e, 0x76, 0x89, 0x17, 0xdd, 0xed,
	0xcf, 0xa4, 0x75, 0xbe, 0xa5, 0x3b, 0x27, 0xb2, 0x70, 0x8e, 0x08, 0x71, 0x91, 0x9b, 0x3c, 0xed,
	0xb9, 0x7b, 0xd3, 0x2c, 0xc4, 0x28, 0xa0, 0xbf, 0xa3, 0x43, 0xda, 0x77, 0xd6, 0x38, 0xf3, 0x21,
	0x2e, 0x20, 0xf9, 0x91, 0xd7, 0x96, 0xd2, 0x38, 0x78, 0x9d, 0x38, 0xfb, 0x13, 0xc1, 0x28, 0x23,
	0xdd, 0x28, 0xd6, 0xf4, 0xac, 0x86, 0x6f, 0x20, 0x96, 0xaa, 0xa0, 0xd0, 0x71, 0x76, 0x79, 0x7a,
	0xcc, 0x5d, 0x3a, 0xab, 0x4b, 0xf4, 0x87, 0x9e, 0x81, 0xda, 0x56, 0xb5, 0x8f, 0x0c, 0x41, 0x78,
	0xd6, 0x5a, 0x6d, 0xd3, 0xa4, 0x63, 0x75, 0xe1, 0xf9, 0xef, 0x08, 0x26, 0x4f, 0xda, 0xe2, 0x29,
	0x4c, 0x6e, 0x6c, 0x5d, 0x1f, 0x2d, 0xf1, 0x02, 0x47, 0x10, 0xaf, 0x64, 0xa9, 0x44, 0x84, 0x63,
	0x48, 0xae, 0x6b, 0xab, 0x4b, 0xd1, 0xf3, 0xe6, 0x9a, 0x37, 0x4a, 0xf4, 0xf1, 0x04, 0xc6, 0x5f,
	0xc9, 0xdc, 0x36, 0xa6, 0x52, 0x2c, 0x62, 0x2f, 0x57, 0x0f, 0x92, 0x3a, 0x99, 0xe0, 0x14, 0x46,
	0x57, 0x4d, 0x43, 0x5c, 0xdc, 0x3d, 0x88, 0x39, 0x4e, 0x60, 0xb8, 0x2c, 0x49, 0xee, 0x9c, 0x70,
	0x83, 0xc1, 0x60, 0xa9, 0xf6, 0xfb, 0xca, 0x88, 0x97, 0xbe, 0xf2, 0x17, 0x4b, 0xed, 0x41, 0x2c,
	0xce, 0xff, 0xba, 0x2d, 0x3d, 0x8e, 0x82, 0x03, 0xe8, 0xdd, 0x7e, 0x72, 0x0c, 0x73, 0x38, 0x59,
	0xb3, 0xa1, 0x96, 0xf3, 0x7a, 0xe5, 0xe7, 0x70, 0x30, 0x02, 0xa6, 0xf7, 0x9c, 0x5b, 0x53, 0xaa,
	0xb6, 0xfa, 0x45, 0x85, 0x63, 0x5a, 0x80, 0x58, 0xb3, 0xb6, 0x9b, 0x4d, 0x25, 0x2b, 0x62, 0x73,
	0x4d, 0xa4, 0x1d, 0x1f, 0xc2, 0xec, 0x9e, 0x77, 0xac, 0x7e, 0xf2, 0xf1, 0x5f, 0x3b, 0x48, 0x57,
	0x6e, 0xc5, 0x6e, 0x4b, 0x15, 0x6f, 0xbb, 0x72, 0x01, 0xf4, 0x63, 0x5e, 0xdc, 0x28, 0x96, 0x24,
	0x06, 0x4f, 0x92, 0xae, 0xa4, 0x54, 0x96, 0x8d, 0x18, 0xe2, 0x2b, 0x98, 0xff, 0x57, 0xde, 0x72,
	0xa1, 0xc5, 0xe8, 0xdb, 0x20, 0x3c, 0xa9, 0x0f, 0xff, 0x02, 0x00, 0x00, 0xff, 0xff, 0x41, 0x3b,
	0x4f, 0x97, 0x67, 0x02, 0x00, 0x00,
}