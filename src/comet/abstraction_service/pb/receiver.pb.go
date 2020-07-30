// Code generated by protoc-gen-go. DO NOT EDIT.
// source: abstraction_service/pb/receiver.proto

package pb

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type PredictRequest struct {
	ImageVector          []int32  `protobuf:"varint,1,rep,packed,name=image_vector,json=imageVector,proto3" json:"image_vector,omitempty"`
	ModelId              int32    `protobuf:"varint,2,opt,name=model_id,json=modelId,proto3" json:"model_id,omitempty"`
	ContextUuid          string   `protobuf:"bytes,3,opt,name=context_uuid,json=contextUuid,proto3" json:"context_uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PredictRequest) Reset()         { *m = PredictRequest{} }
func (m *PredictRequest) String() string { return proto.CompactTextString(m) }
func (*PredictRequest) ProtoMessage()    {}
func (*PredictRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_307994286c2b9491, []int{0}
}

func (m *PredictRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PredictRequest.Unmarshal(m, b)
}
func (m *PredictRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PredictRequest.Marshal(b, m, deterministic)
}
func (m *PredictRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PredictRequest.Merge(m, src)
}
func (m *PredictRequest) XXX_Size() int {
	return xxx_messageInfo_PredictRequest.Size(m)
}
func (m *PredictRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_PredictRequest.DiscardUnknown(m)
}

var xxx_messageInfo_PredictRequest proto.InternalMessageInfo

func (m *PredictRequest) GetImageVector() []int32 {
	if m != nil {
		return m.ImageVector
	}
	return nil
}

func (m *PredictRequest) GetModelId() int32 {
	if m != nil {
		return m.ModelId
	}
	return 0
}

func (m *PredictRequest) GetContextUuid() string {
	if m != nil {
		return m.ContextUuid
	}
	return ""
}

// in the future we may want to return some data
type PredictReply struct {
	Label                string   `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *PredictReply) Reset()         { *m = PredictReply{} }
func (m *PredictReply) String() string { return proto.CompactTextString(m) }
func (*PredictReply) ProtoMessage()    {}
func (*PredictReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_307994286c2b9491, []int{1}
}

func (m *PredictReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_PredictReply.Unmarshal(m, b)
}
func (m *PredictReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_PredictReply.Marshal(b, m, deterministic)
}
func (m *PredictReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_PredictReply.Merge(m, src)
}
func (m *PredictReply) XXX_Size() int {
	return xxx_messageInfo_PredictReply.Size(m)
}
func (m *PredictReply) XXX_DiscardUnknown() {
	xxx_messageInfo_PredictReply.DiscardUnknown(m)
}

var xxx_messageInfo_PredictReply proto.InternalMessageInfo

func (m *PredictReply) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func init() {
	proto.RegisterType((*PredictRequest)(nil), "pb.PredictRequest")
	proto.RegisterType((*PredictReply)(nil), "pb.PredictReply")
}

func init() {
	proto.RegisterFile("abstraction_service/pb/receiver.proto", fileDescriptor_307994286c2b9491)
}

var fileDescriptor_307994286c2b9491 = []byte{
	// 226 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x54, 0x8f, 0x4d, 0x4b, 0x03, 0x31,
	0x10, 0x86, 0x6d, 0x6b, 0xad, 0x4e, 0x8b, 0xc8, 0xe0, 0xa1, 0x7a, 0xaa, 0x8b, 0x82, 0xa7, 0x2d,
	0xea, 0x2f, 0xd0, 0x8b, 0x78, 0x93, 0x15, 0x3d, 0x78, 0x09, 0xf9, 0x18, 0x24, 0x90, 0x36, 0x69,
	0x76, 0x76, 0xd1, 0x7f, 0x6f, 0x9a, 0x5d, 0x94, 0x3d, 0xce, 0xf3, 0x66, 0xf2, 0xce, 0x03, 0x37,
	0x52, 0xd5, 0x1c, 0xa5, 0x66, 0xeb, 0xb7, 0xa2, 0xa6, 0xd8, 0x5a, 0x4d, 0xeb, 0xa0, 0xd6, 0x91,
	0x34, 0xd9, 0x96, 0x62, 0x19, 0xa2, 0x67, 0x8f, 0xe3, 0xa0, 0x8a, 0x1d, 0x9c, 0xbe, 0x46, 0x32,
	0x56, 0x73, 0x45, 0xbb, 0x86, 0x6a, 0xc6, 0x2b, 0x58, 0xd8, 0x8d, 0xfc, 0x22, 0xd1, 0x92, 0x66,
	0x1f, 0x97, 0xa3, 0xd5, 0xe4, 0x76, 0x5a, 0xcd, 0x33, 0xfb, 0xc8, 0x08, 0x2f, 0xe0, 0x78, 0xe3,
	0x0d, 0x39, 0x61, 0xcd, 0x72, 0xbc, 0x1a, 0xa5, 0x78, 0x96, 0xe7, 0x17, 0xb3, 0xdf, 0xd6, 0x7e,
	0xcb, 0xf4, 0xcd, 0xa2, 0x69, 0x52, 0x3c, 0x49, 0xf1, 0x49, 0x35, 0xef, 0xd9, 0x7b, 0x42, 0xc5,
	0x35, 0x2c, 0xfe, 0x2a, 0x83, 0xfb, 0xc1, 0x73, 0x98, 0x3a, 0xa9, 0xc8, 0xa5, 0xa6, 0xfd, 0xdb,
	0x6e, 0xb8, 0x7f, 0x06, 0x7c, 0xfc, 0xb7, 0x78, 0xeb, 0x24, 0xf0, 0x0e, 0x66, 0xfd, 0x2e, 0x62,
	0x19, 0x54, 0x39, 0xbc, 0xfd, 0xf2, 0x6c, 0xc0, 0xd2, 0xe7, 0xc5, 0xc1, 0xd3, 0xe1, 0x67, 0xf2,
	0x54, 0x47, 0x59, 0xf9, 0xe1, 0x37, 0x00, 0x00, 0xff, 0xff, 0x61, 0x85, 0x9b, 0xed, 0x1b, 0x01,
	0x00, 0x00,
}
