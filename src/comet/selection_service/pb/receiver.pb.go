// Code generated by protoc-gen-go. DO NOT EDIT.
// source: selection_service/pb/receiver.proto

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

type QueryRequest struct {
	FeatureVector        []int32  `protobuf:"varint,1,rep,packed,name=feature_vector,json=featureVector,proto3" json:"feature_vector,omitempty"`
	FType                string   `protobuf:"bytes,2,opt,name=fType,proto3" json:"fType,omitempty"`
	XDim                 int32    `protobuf:"varint,3,opt,name=xDim,proto3" json:"xDim,omitempty"`
	YDim                 int32    `protobuf:"varint,4,opt,name=yDim,proto3" json:"yDim,omitempty"`
	Layers               int32    `protobuf:"varint,5,opt,name=layers,proto3" json:"layers,omitempty"`
	ContextUuid          string   `protobuf:"bytes,6,opt,name=context_uuid,json=contextUuid,proto3" json:"context_uuid,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryRequest) Reset()         { *m = QueryRequest{} }
func (m *QueryRequest) String() string { return proto.CompactTextString(m) }
func (*QueryRequest) ProtoMessage()    {}
func (*QueryRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8b0d817558cbf0b1, []int{0}
}

func (m *QueryRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryRequest.Unmarshal(m, b)
}
func (m *QueryRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryRequest.Marshal(b, m, deterministic)
}
func (m *QueryRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryRequest.Merge(m, src)
}
func (m *QueryRequest) XXX_Size() int {
	return xxx_messageInfo_QueryRequest.Size(m)
}
func (m *QueryRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryRequest.DiscardUnknown(m)
}

var xxx_messageInfo_QueryRequest proto.InternalMessageInfo

func (m *QueryRequest) GetFeatureVector() []int32 {
	if m != nil {
		return m.FeatureVector
	}
	return nil
}

func (m *QueryRequest) GetFType() string {
	if m != nil {
		return m.FType
	}
	return ""
}

func (m *QueryRequest) GetXDim() int32 {
	if m != nil {
		return m.XDim
	}
	return 0
}

func (m *QueryRequest) GetYDim() int32 {
	if m != nil {
		return m.YDim
	}
	return 0
}

func (m *QueryRequest) GetLayers() int32 {
	if m != nil {
		return m.Layers
	}
	return 0
}

func (m *QueryRequest) GetContextUuid() string {
	if m != nil {
		return m.ContextUuid
	}
	return ""
}

type QueryReply struct {
	Label                string   `protobuf:"bytes,1,opt,name=label,proto3" json:"label,omitempty"`
	QueryID              int32    `protobuf:"varint,2,opt,name=queryID,proto3" json:"queryID,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *QueryReply) Reset()         { *m = QueryReply{} }
func (m *QueryReply) String() string { return proto.CompactTextString(m) }
func (*QueryReply) ProtoMessage()    {}
func (*QueryReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_8b0d817558cbf0b1, []int{1}
}

func (m *QueryReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_QueryReply.Unmarshal(m, b)
}
func (m *QueryReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_QueryReply.Marshal(b, m, deterministic)
}
func (m *QueryReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_QueryReply.Merge(m, src)
}
func (m *QueryReply) XXX_Size() int {
	return xxx_messageInfo_QueryReply.Size(m)
}
func (m *QueryReply) XXX_DiscardUnknown() {
	xxx_messageInfo_QueryReply.DiscardUnknown(m)
}

var xxx_messageInfo_QueryReply proto.InternalMessageInfo

func (m *QueryReply) GetLabel() string {
	if m != nil {
		return m.Label
	}
	return ""
}

func (m *QueryReply) GetQueryID() int32 {
	if m != nil {
		return m.QueryID
	}
	return 0
}

type FeedbackRequest struct {
	ContextUuid          string   `protobuf:"bytes,1,opt,name=context_uuid,json=contextUuid,proto3" json:"context_uuid,omitempty"`
	QueryID              string   `protobuf:"bytes,2,opt,name=queryID,proto3" json:"queryID,omitempty"`
	Evaluation           int32    `protobuf:"varint,3,opt,name=evaluation,proto3" json:"evaluation,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FeedbackRequest) Reset()         { *m = FeedbackRequest{} }
func (m *FeedbackRequest) String() string { return proto.CompactTextString(m) }
func (*FeedbackRequest) ProtoMessage()    {}
func (*FeedbackRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_8b0d817558cbf0b1, []int{2}
}

func (m *FeedbackRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeedbackRequest.Unmarshal(m, b)
}
func (m *FeedbackRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeedbackRequest.Marshal(b, m, deterministic)
}
func (m *FeedbackRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeedbackRequest.Merge(m, src)
}
func (m *FeedbackRequest) XXX_Size() int {
	return xxx_messageInfo_FeedbackRequest.Size(m)
}
func (m *FeedbackRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FeedbackRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FeedbackRequest proto.InternalMessageInfo

func (m *FeedbackRequest) GetContextUuid() string {
	if m != nil {
		return m.ContextUuid
	}
	return ""
}

func (m *FeedbackRequest) GetQueryID() string {
	if m != nil {
		return m.QueryID
	}
	return ""
}

func (m *FeedbackRequest) GetEvaluation() int32 {
	if m != nil {
		return m.Evaluation
	}
	return 0
}

type FeedbackReply struct {
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FeedbackReply) Reset()         { *m = FeedbackReply{} }
func (m *FeedbackReply) String() string { return proto.CompactTextString(m) }
func (*FeedbackReply) ProtoMessage()    {}
func (*FeedbackReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_8b0d817558cbf0b1, []int{3}
}

func (m *FeedbackReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FeedbackReply.Unmarshal(m, b)
}
func (m *FeedbackReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FeedbackReply.Marshal(b, m, deterministic)
}
func (m *FeedbackReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FeedbackReply.Merge(m, src)
}
func (m *FeedbackReply) XXX_Size() int {
	return xxx_messageInfo_FeedbackReply.Size(m)
}
func (m *FeedbackReply) XXX_DiscardUnknown() {
	xxx_messageInfo_FeedbackReply.DiscardUnknown(m)
}

var xxx_messageInfo_FeedbackReply proto.InternalMessageInfo

func init() {
	proto.RegisterType((*QueryRequest)(nil), "pb.QueryRequest")
	proto.RegisterType((*QueryReply)(nil), "pb.QueryReply")
	proto.RegisterType((*FeedbackRequest)(nil), "pb.FeedbackRequest")
	proto.RegisterType((*FeedbackReply)(nil), "pb.FeedbackReply")
}

func init() {
	proto.RegisterFile("selection_service/pb/receiver.proto", fileDescriptor_8b0d817558cbf0b1)
}

var fileDescriptor_8b0d817558cbf0b1 = []byte{
	// 330 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x64, 0x91, 0xcf, 0x4e, 0xf2, 0x40,
	0x14, 0xc5, 0xbf, 0x02, 0xe5, 0x93, 0x2b, 0xff, 0xbc, 0x1a, 0xd3, 0xb0, 0x30, 0x58, 0x63, 0x42,
	0x62, 0x02, 0x89, 0xba, 0x74, 0x65, 0x88, 0x89, 0x4b, 0x8b, 0xba, 0x70, 0x43, 0xda, 0x72, 0x49,
	0x1a, 0x47, 0x3a, 0x4c, 0x67, 0x1a, 0xfa, 0x52, 0x3e, 0xa3, 0x33, 0x43, 0xab, 0x88, 0xbb, 0xb9,
	0xbf, 0x76, 0xce, 0x3d, 0x67, 0x0e, 0x5c, 0x64, 0xc4, 0x28, 0x96, 0x49, 0xba, 0x9a, 0x67, 0x24,
	0xf2, 0x24, 0xa6, 0x09, 0x8f, 0x26, 0x82, 0x62, 0x4a, 0x72, 0x12, 0x63, 0x2e, 0x52, 0x99, 0x62,
	0x8d, 0x47, 0xfe, 0xa7, 0x03, 0xed, 0x27, 0x45, 0xa2, 0x08, 0x68, 0xad, 0x28, 0x93, 0x78, 0x09,
	0xdd, 0x25, 0x85, 0x52, 0x09, 0x9a, 0xe7, 0x5a, 0x21, 0x15, 0x9e, 0x33, 0xac, 0x8f, 0xdc, 0xa0,
	0x53, 0xd2, 0x57, 0x0b, 0xf1, 0x04, 0xdc, 0xe5, 0x73, 0xc1, 0xc9, 0xab, 0x0d, 0x9d, 0x51, 0x2b,
	0xd8, 0x0e, 0x88, 0xd0, 0xd8, 0x4c, 0x93, 0x0f, 0xaf, 0xae, 0xa1, 0x1b, 0xd8, 0xb3, 0x61, 0x85,
	0x61, 0x8d, 0x2d, 0x33, 0x67, 0x3c, 0x85, 0x26, 0x0b, 0x0b, 0x12, 0x99, 0xe7, 0x5a, 0x5a, 0x4e,
	0x78, 0x0e, 0xed, 0x38, 0x5d, 0x49, 0xda, 0xc8, 0xb9, 0x52, 0xc9, 0xc2, 0x6b, 0x5a, 0xf1, 0xc3,
	0x92, 0xbd, 0x68, 0xe4, 0xdf, 0x01, 0x94, 0x7e, 0x39, 0x2b, 0x8c, 0x0d, 0x16, 0x46, 0xc4, 0xb4,
	0x49, 0x6b, 0xc3, 0x0e, 0xe8, 0xc1, 0xff, 0xb5, 0xf9, 0xe7, 0x71, 0x6a, 0xed, 0xb9, 0x41, 0x35,
	0xfa, 0x2b, 0xe8, 0x3d, 0x10, 0x2d, 0xa2, 0x30, 0x7e, 0xaf, 0x02, 0xef, 0xef, 0x74, 0xfe, 0xec,
	0xdc, 0xd7, 0x6b, 0x7d, 0xeb, 0xe1, 0x19, 0x00, 0xe5, 0x21, 0x53, 0xa1, 0x79, 0xea, 0x32, 0xf6,
	0x0e, 0xf1, 0x7b, 0xd0, 0xf9, 0xd9, 0xa7, 0x0d, 0x5f, 0x2b, 0xe8, 0xcf, 0xaa, 0x6a, 0x66, 0xdb,
	0x66, 0xf0, 0x0a, 0x5c, 0x1b, 0x09, 0xfb, 0x63, 0x1e, 0x8d, 0x77, 0xdb, 0x18, 0x74, 0x77, 0x88,
	0xbe, 0xee, 0xff, 0xc3, 0x5b, 0x38, 0xa8, 0x14, 0xf1, 0xd8, 0x7c, 0xdd, 0xcb, 0x33, 0x38, 0xfa,
	0x0d, 0xed, 0xad, 0xfb, 0xc6, 0x9b, 0x2e, 0x3b, 0x6a, 0xda, 0xde, 0x6f, 0xbe, 0x02, 0x00, 0x00,
	0xff, 0xff, 0x33, 0x0a, 0xcd, 0x74, 0x1e, 0x02, 0x00, 0x00,
}
