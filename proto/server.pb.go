// Code generated by protoc-gen-go. DO NOT EDIT.
// source: server.proto

package proto

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
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

type CreateRequest struct {
	Document             []byte   `protobuf:"bytes,1,opt,name=document,proto3" json:"document,omitempty"`
	Operation            string   `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
	Meta                 *Meta    `protobuf:"bytes,3,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{0}
}

func (m *CreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRequest.Unmarshal(m, b)
}
func (m *CreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRequest.Marshal(b, m, deterministic)
}
func (m *CreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRequest.Merge(m, src)
}
func (m *CreateRequest) XXX_Size() int {
	return xxx_messageInfo_CreateRequest.Size(m)
}
func (m *CreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRequest proto.InternalMessageInfo

func (m *CreateRequest) GetDocument() []byte {
	if m != nil {
		return m.Document
	}
	return nil
}

func (m *CreateRequest) GetOperation() string {
	if m != nil {
		return m.Operation
	}
	return ""
}

func (m *CreateRequest) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type ReadRequest struct {
	Find                 []byte       `protobuf:"bytes,1,opt,name=find,proto3" json:"find,omitempty"`
	Operation            string       `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
	Options              *ReadOptions `protobuf:"bytes,3,opt,name=options,proto3" json:"options,omitempty"`
	Meta                 *Meta        `protobuf:"bytes,4,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}     `json:"-"`
	XXX_unrecognized     []byte       `json:"-"`
	XXX_sizecache        int32        `json:"-"`
}

func (m *ReadRequest) Reset()         { *m = ReadRequest{} }
func (m *ReadRequest) String() string { return proto.CompactTextString(m) }
func (*ReadRequest) ProtoMessage()    {}
func (*ReadRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{1}
}

func (m *ReadRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReadRequest.Unmarshal(m, b)
}
func (m *ReadRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReadRequest.Marshal(b, m, deterministic)
}
func (m *ReadRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReadRequest.Merge(m, src)
}
func (m *ReadRequest) XXX_Size() int {
	return xxx_messageInfo_ReadRequest.Size(m)
}
func (m *ReadRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ReadRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ReadRequest proto.InternalMessageInfo

func (m *ReadRequest) GetFind() []byte {
	if m != nil {
		return m.Find
	}
	return nil
}

func (m *ReadRequest) GetOperation() string {
	if m != nil {
		return m.Operation
	}
	return ""
}

func (m *ReadRequest) GetOptions() *ReadOptions {
	if m != nil {
		return m.Options
	}
	return nil
}

func (m *ReadRequest) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type ReadOptions struct {
	Select               map[string]int32 `protobuf:"bytes,1,rep,name=select,proto3" json:"select,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Sort                 map[string]int32 `protobuf:"bytes,2,rep,name=sort,proto3" json:"sort,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"varint,2,opt,name=value,proto3"`
	Skip                 int64            `protobuf:"varint,3,opt,name=skip,proto3" json:"skip,omitempty"`
	Limit                int64            `protobuf:"varint,4,opt,name=limit,proto3" json:"limit,omitempty"`
	Distinct             string           `protobuf:"bytes,5,opt,name=distinct,proto3" json:"distinct,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *ReadOptions) Reset()         { *m = ReadOptions{} }
func (m *ReadOptions) String() string { return proto.CompactTextString(m) }
func (*ReadOptions) ProtoMessage()    {}
func (*ReadOptions) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{2}
}

func (m *ReadOptions) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ReadOptions.Unmarshal(m, b)
}
func (m *ReadOptions) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ReadOptions.Marshal(b, m, deterministic)
}
func (m *ReadOptions) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ReadOptions.Merge(m, src)
}
func (m *ReadOptions) XXX_Size() int {
	return xxx_messageInfo_ReadOptions.Size(m)
}
func (m *ReadOptions) XXX_DiscardUnknown() {
	xxx_messageInfo_ReadOptions.DiscardUnknown(m)
}

var xxx_messageInfo_ReadOptions proto.InternalMessageInfo

func (m *ReadOptions) GetSelect() map[string]int32 {
	if m != nil {
		return m.Select
	}
	return nil
}

func (m *ReadOptions) GetSort() map[string]int32 {
	if m != nil {
		return m.Sort
	}
	return nil
}

func (m *ReadOptions) GetSkip() int64 {
	if m != nil {
		return m.Skip
	}
	return 0
}

func (m *ReadOptions) GetLimit() int64 {
	if m != nil {
		return m.Limit
	}
	return 0
}

func (m *ReadOptions) GetDistinct() string {
	if m != nil {
		return m.Distinct
	}
	return ""
}

type UpdateRequest struct {
	Find                 []byte   `protobuf:"bytes,1,opt,name=find,proto3" json:"find,omitempty"`
	Operation            string   `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
	Update               []byte   `protobuf:"bytes,3,opt,name=update,proto3" json:"update,omitempty"`
	Meta                 *Meta    `protobuf:"bytes,4,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateRequest) Reset()         { *m = UpdateRequest{} }
func (m *UpdateRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateRequest) ProtoMessage()    {}
func (*UpdateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{3}
}

func (m *UpdateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateRequest.Unmarshal(m, b)
}
func (m *UpdateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateRequest.Marshal(b, m, deterministic)
}
func (m *UpdateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateRequest.Merge(m, src)
}
func (m *UpdateRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateRequest.Size(m)
}
func (m *UpdateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateRequest proto.InternalMessageInfo

func (m *UpdateRequest) GetFind() []byte {
	if m != nil {
		return m.Find
	}
	return nil
}

func (m *UpdateRequest) GetOperation() string {
	if m != nil {
		return m.Operation
	}
	return ""
}

func (m *UpdateRequest) GetUpdate() []byte {
	if m != nil {
		return m.Update
	}
	return nil
}

func (m *UpdateRequest) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type DeleteRequest struct {
	Find                 []byte   `protobuf:"bytes,1,opt,name=find,proto3" json:"find,omitempty"`
	Operation            string   `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
	Meta                 *Meta    `protobuf:"bytes,3,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DeleteRequest) Reset()         { *m = DeleteRequest{} }
func (m *DeleteRequest) String() string { return proto.CompactTextString(m) }
func (*DeleteRequest) ProtoMessage()    {}
func (*DeleteRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{4}
}

func (m *DeleteRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteRequest.Unmarshal(m, b)
}
func (m *DeleteRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteRequest.Marshal(b, m, deterministic)
}
func (m *DeleteRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteRequest.Merge(m, src)
}
func (m *DeleteRequest) XXX_Size() int {
	return xxx_messageInfo_DeleteRequest.Size(m)
}
func (m *DeleteRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteRequest.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteRequest proto.InternalMessageInfo

func (m *DeleteRequest) GetFind() []byte {
	if m != nil {
		return m.Find
	}
	return nil
}

func (m *DeleteRequest) GetOperation() string {
	if m != nil {
		return m.Operation
	}
	return ""
}

func (m *DeleteRequest) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type AggregateRequest struct {
	Pipeline             []byte   `protobuf:"bytes,1,opt,name=pipeline,proto3" json:"pipeline,omitempty"`
	Operation            string   `protobuf:"bytes,2,opt,name=operation,proto3" json:"operation,omitempty"`
	Meta                 *Meta    `protobuf:"bytes,3,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AggregateRequest) Reset()         { *m = AggregateRequest{} }
func (m *AggregateRequest) String() string { return proto.CompactTextString(m) }
func (*AggregateRequest) ProtoMessage()    {}
func (*AggregateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{5}
}

func (m *AggregateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AggregateRequest.Unmarshal(m, b)
}
func (m *AggregateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AggregateRequest.Marshal(b, m, deterministic)
}
func (m *AggregateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AggregateRequest.Merge(m, src)
}
func (m *AggregateRequest) XXX_Size() int {
	return xxx_messageInfo_AggregateRequest.Size(m)
}
func (m *AggregateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AggregateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AggregateRequest proto.InternalMessageInfo

func (m *AggregateRequest) GetPipeline() []byte {
	if m != nil {
		return m.Pipeline
	}
	return nil
}

func (m *AggregateRequest) GetOperation() string {
	if m != nil {
		return m.Operation
	}
	return ""
}

func (m *AggregateRequest) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type Response struct {
	Status               int32    `protobuf:"varint,1,opt,name=status,proto3" json:"status,omitempty"`
	Error                string   `protobuf:"bytes,2,opt,name=error,proto3" json:"error,omitempty"`
	Result               []byte   `protobuf:"bytes,3,opt,name=result,proto3" json:"result,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Response) Reset()         { *m = Response{} }
func (m *Response) String() string { return proto.CompactTextString(m) }
func (*Response) ProtoMessage()    {}
func (*Response) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{6}
}

func (m *Response) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Response.Unmarshal(m, b)
}
func (m *Response) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Response.Marshal(b, m, deterministic)
}
func (m *Response) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Response.Merge(m, src)
}
func (m *Response) XXX_Size() int {
	return xxx_messageInfo_Response.Size(m)
}
func (m *Response) XXX_DiscardUnknown() {
	xxx_messageInfo_Response.DiscardUnknown(m)
}

var xxx_messageInfo_Response proto.InternalMessageInfo

func (m *Response) GetStatus() int32 {
	if m != nil {
		return m.Status
	}
	return 0
}

func (m *Response) GetError() string {
	if m != nil {
		return m.Error
	}
	return ""
}

func (m *Response) GetResult() []byte {
	if m != nil {
		return m.Result
	}
	return nil
}

type Meta struct {
	Project              string   `protobuf:"bytes,1,opt,name=project,proto3" json:"project,omitempty"`
	DbType               string   `protobuf:"bytes,2,opt,name=dbType,proto3" json:"dbType,omitempty"`
	Col                  string   `protobuf:"bytes,3,opt,name=col,proto3" json:"col,omitempty"`
	Token                string   `protobuf:"bytes,4,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Meta) Reset()         { *m = Meta{} }
func (m *Meta) String() string { return proto.CompactTextString(m) }
func (*Meta) ProtoMessage()    {}
func (*Meta) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{7}
}

func (m *Meta) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Meta.Unmarshal(m, b)
}
func (m *Meta) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Meta.Marshal(b, m, deterministic)
}
func (m *Meta) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Meta.Merge(m, src)
}
func (m *Meta) XXX_Size() int {
	return xxx_messageInfo_Meta.Size(m)
}
func (m *Meta) XXX_DiscardUnknown() {
	xxx_messageInfo_Meta.DiscardUnknown(m)
}

var xxx_messageInfo_Meta proto.InternalMessageInfo

func (m *Meta) GetProject() string {
	if m != nil {
		return m.Project
	}
	return ""
}

func (m *Meta) GetDbType() string {
	if m != nil {
		return m.DbType
	}
	return ""
}

func (m *Meta) GetCol() string {
	if m != nil {
		return m.Col
	}
	return ""
}

func (m *Meta) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type AllRequest struct {
	Col                  string   `protobuf:"bytes,1,opt,name=col,proto3" json:"col,omitempty"`
	Document             []byte   `protobuf:"bytes,2,opt,name=document,proto3" json:"document,omitempty"`
	Operation            string   `protobuf:"bytes,3,opt,name=operation,proto3" json:"operation,omitempty"`
	Find                 []byte   `protobuf:"bytes,4,opt,name=find,proto3" json:"find,omitempty"`
	Update               []byte   `protobuf:"bytes,5,opt,name=update,proto3" json:"update,omitempty"`
	Type                 string   `protobuf:"bytes,6,opt,name=type,proto3" json:"type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *AllRequest) Reset()         { *m = AllRequest{} }
func (m *AllRequest) String() string { return proto.CompactTextString(m) }
func (*AllRequest) ProtoMessage()    {}
func (*AllRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{8}
}

func (m *AllRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_AllRequest.Unmarshal(m, b)
}
func (m *AllRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_AllRequest.Marshal(b, m, deterministic)
}
func (m *AllRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_AllRequest.Merge(m, src)
}
func (m *AllRequest) XXX_Size() int {
	return xxx_messageInfo_AllRequest.Size(m)
}
func (m *AllRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_AllRequest.DiscardUnknown(m)
}

var xxx_messageInfo_AllRequest proto.InternalMessageInfo

func (m *AllRequest) GetCol() string {
	if m != nil {
		return m.Col
	}
	return ""
}

func (m *AllRequest) GetDocument() []byte {
	if m != nil {
		return m.Document
	}
	return nil
}

func (m *AllRequest) GetOperation() string {
	if m != nil {
		return m.Operation
	}
	return ""
}

func (m *AllRequest) GetFind() []byte {
	if m != nil {
		return m.Find
	}
	return nil
}

func (m *AllRequest) GetUpdate() []byte {
	if m != nil {
		return m.Update
	}
	return nil
}

func (m *AllRequest) GetType() string {
	if m != nil {
		return m.Type
	}
	return ""
}

type BatchRequest struct {
	Batchrequest         []*AllRequest `protobuf:"bytes,1,rep,name=batchrequest,proto3" json:"batchrequest,omitempty"`
	Meta                 *Meta         `protobuf:"bytes,2,opt,name=meta,proto3" json:"meta,omitempty"`
	XXX_NoUnkeyedLiteral struct{}      `json:"-"`
	XXX_unrecognized     []byte        `json:"-"`
	XXX_sizecache        int32         `json:"-"`
}

func (m *BatchRequest) Reset()         { *m = BatchRequest{} }
func (m *BatchRequest) String() string { return proto.CompactTextString(m) }
func (*BatchRequest) ProtoMessage()    {}
func (*BatchRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{9}
}

func (m *BatchRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_BatchRequest.Unmarshal(m, b)
}
func (m *BatchRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_BatchRequest.Marshal(b, m, deterministic)
}
func (m *BatchRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_BatchRequest.Merge(m, src)
}
func (m *BatchRequest) XXX_Size() int {
	return xxx_messageInfo_BatchRequest.Size(m)
}
func (m *BatchRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_BatchRequest.DiscardUnknown(m)
}

var xxx_messageInfo_BatchRequest proto.InternalMessageInfo

func (m *BatchRequest) GetBatchrequest() []*AllRequest {
	if m != nil {
		return m.Batchrequest
	}
	return nil
}

func (m *BatchRequest) GetMeta() *Meta {
	if m != nil {
		return m.Meta
	}
	return nil
}

type FaaSRequest struct {
	Params               []byte   `protobuf:"bytes,1,opt,name=params,proto3" json:"params,omitempty"`
	Timeout              int64    `protobuf:"varint,2,opt,name=timeout,proto3" json:"timeout,omitempty"`
	Engine               string   `protobuf:"bytes,3,opt,name=engine,proto3" json:"engine,omitempty"`
	Function             string   `protobuf:"bytes,4,opt,name=function,proto3" json:"function,omitempty"`
	Token                string   `protobuf:"bytes,5,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *FaaSRequest) Reset()         { *m = FaaSRequest{} }
func (m *FaaSRequest) String() string { return proto.CompactTextString(m) }
func (*FaaSRequest) ProtoMessage()    {}
func (*FaaSRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_ad098daeda4239f7, []int{10}
}

func (m *FaaSRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_FaaSRequest.Unmarshal(m, b)
}
func (m *FaaSRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_FaaSRequest.Marshal(b, m, deterministic)
}
func (m *FaaSRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_FaaSRequest.Merge(m, src)
}
func (m *FaaSRequest) XXX_Size() int {
	return xxx_messageInfo_FaaSRequest.Size(m)
}
func (m *FaaSRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_FaaSRequest.DiscardUnknown(m)
}

var xxx_messageInfo_FaaSRequest proto.InternalMessageInfo

func (m *FaaSRequest) GetParams() []byte {
	if m != nil {
		return m.Params
	}
	return nil
}

func (m *FaaSRequest) GetTimeout() int64 {
	if m != nil {
		return m.Timeout
	}
	return 0
}

func (m *FaaSRequest) GetEngine() string {
	if m != nil {
		return m.Engine
	}
	return ""
}

func (m *FaaSRequest) GetFunction() string {
	if m != nil {
		return m.Function
	}
	return ""
}

func (m *FaaSRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func init() {
	proto.RegisterType((*CreateRequest)(nil), "proto.CreateRequest")
	proto.RegisterType((*ReadRequest)(nil), "proto.ReadRequest")
	proto.RegisterType((*ReadOptions)(nil), "proto.ReadOptions")
	proto.RegisterMapType((map[string]int32)(nil), "proto.ReadOptions.SelectEntry")
	proto.RegisterMapType((map[string]int32)(nil), "proto.ReadOptions.SortEntry")
	proto.RegisterType((*UpdateRequest)(nil), "proto.UpdateRequest")
	proto.RegisterType((*DeleteRequest)(nil), "proto.DeleteRequest")
	proto.RegisterType((*AggregateRequest)(nil), "proto.AggregateRequest")
	proto.RegisterType((*Response)(nil), "proto.Response")
	proto.RegisterType((*Meta)(nil), "proto.Meta")
	proto.RegisterType((*AllRequest)(nil), "proto.AllRequest")
	proto.RegisterType((*BatchRequest)(nil), "proto.BatchRequest")
	proto.RegisterType((*FaaSRequest)(nil), "proto.FaaSRequest")
}

func init() { proto.RegisterFile("server.proto", fileDescriptor_ad098daeda4239f7) }

var fileDescriptor_ad098daeda4239f7 = []byte{
	// 718 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x54, 0xcd, 0x52, 0x13, 0x41,
	0x10, 0x76, 0x93, 0xdd, 0x40, 0x3a, 0xa1, 0xc0, 0x91, 0xc2, 0x54, 0x8a, 0x12, 0x2a, 0xa7, 0x1c,
	0x34, 0x0a, 0xfe, 0xa0, 0xde, 0x00, 0xf5, 0x66, 0x49, 0x0d, 0x7a, 0xd6, 0xc9, 0xa6, 0x09, 0x0b,
	0x9b, 0x9d, 0x71, 0x66, 0x96, 0x2a, 0x7c, 0x03, 0x2f, 0x9e, 0x7d, 0x0d, 0x1f, 0xcc, 0x77, 0xb0,
	0xe6, 0x67, 0x37, 0x1b, 0x4c, 0xa0, 0x28, 0x4e, 0x99, 0x2f, 0xdb, 0x5f, 0xff, 0x7d, 0xdd, 0x0d,
	0x6d, 0x85, 0xf2, 0x02, 0xe5, 0x40, 0x48, 0xae, 0x39, 0x89, 0xec, 0x4f, 0xef, 0x0c, 0x56, 0x0e,
	0x25, 0x32, 0x8d, 0x14, 0xbf, 0xe7, 0xa8, 0x34, 0xe9, 0xc2, 0xf2, 0x88, 0xc7, 0xf9, 0x04, 0x33,
	0xdd, 0x09, 0xb6, 0x83, 0x7e, 0x9b, 0x96, 0x98, 0x6c, 0x42, 0x93, 0x0b, 0x94, 0x4c, 0x27, 0x3c,
	0xeb, 0xd4, 0xb6, 0x83, 0x7e, 0x93, 0x4e, 0xff, 0x20, 0x5b, 0x10, 0x4e, 0x50, 0xb3, 0x4e, 0x7d,
	0x3b, 0xe8, 0xb7, 0x76, 0x5b, 0x2e, 0xce, 0xe0, 0x23, 0x6a, 0x46, 0xed, 0x87, 0xde, 0xaf, 0x00,
	0x5a, 0x14, 0xd9, 0xa8, 0x08, 0x45, 0x20, 0x3c, 0x49, 0xb2, 0x91, 0x0f, 0x63, 0xdf, 0x37, 0x84,
	0x78, 0x0c, 0x4b, 0x5c, 0x98, 0x97, 0xf2, 0x51, 0x88, 0x8f, 0x62, 0xdc, 0x7e, 0x72, 0x5f, 0x68,
	0x61, 0x52, 0x26, 0x14, 0x2e, 0x4a, 0xe8, 0x4f, 0xcd, 0x25, 0xe4, 0x99, 0xe4, 0x15, 0x34, 0x14,
	0xa6, 0x18, 0x9b, 0xca, 0xeb, 0xfd, 0xd6, 0xee, 0xa3, 0xff, 0xbd, 0x0f, 0x8e, 0xad, 0xc1, 0xfb,
	0x4c, 0xcb, 0x4b, 0xea, 0xad, 0xc9, 0x33, 0x08, 0x15, 0x97, 0xba, 0x53, 0xb3, 0xac, 0xcd, 0x79,
	0x2c, 0x2e, 0x3d, 0xc7, 0x5a, 0x9a, 0xd2, 0xd5, 0x79, 0x22, 0x6c, 0x15, 0x75, 0x6a, 0xdf, 0x64,
	0x1d, 0xa2, 0x34, 0x99, 0x24, 0xda, 0xe6, 0x5b, 0xa7, 0x0e, 0x58, 0x3d, 0x12, 0xa5, 0x93, 0x2c,
	0xd6, 0x9d, 0xc8, 0xf6, 0xa3, 0xc4, 0xdd, 0x37, 0xd0, 0xaa, 0xa4, 0x43, 0xd6, 0xa0, 0x7e, 0x8e,
	0x97, 0xb6, 0x9d, 0x4d, 0x6a, 0x9e, 0xc6, 0xe5, 0x05, 0x4b, 0x73, 0xb4, 0x9d, 0x8c, 0xa8, 0x03,
	0x6f, 0x6b, 0xaf, 0x83, 0xee, 0x1e, 0x34, 0xcb, 0x9c, 0x6e, 0x43, 0xec, 0xfd, 0x80, 0x95, 0x2f,
	0x62, 0x54, 0x19, 0x98, 0xdb, 0xab, 0xb8, 0x01, 0x8d, 0xdc, 0xba, 0xb0, 0xe5, 0xb7, 0xa9, 0x47,
	0x37, 0xeb, 0x35, 0x84, 0x95, 0x77, 0x98, 0xe2, 0x5d, 0x62, 0xdf, 0x38, 0xa4, 0x13, 0x58, 0xdb,
	0x1f, 0x8f, 0x25, 0x8e, 0x67, 0x77, 0x42, 0x24, 0x02, 0xd3, 0x24, 0xc3, 0x62, 0x27, 0x0a, 0x7c,
	0xd7, 0x70, 0x47, 0xb0, 0x4c, 0x51, 0x09, 0x9e, 0x29, 0x34, 0x7d, 0x51, 0x9a, 0xe9, 0x5c, 0xd9,
	0x20, 0x11, 0xf5, 0xc8, 0x88, 0x81, 0x52, 0x72, 0xe9, 0xdd, 0x3b, 0x60, 0xac, 0x25, 0xaa, 0x3c,
	0xd5, 0x45, 0x17, 0x1d, 0xea, 0x7d, 0x83, 0xd0, 0xf8, 0x27, 0x1d, 0x58, 0x12, 0x92, 0x9f, 0xb9,
	0x69, 0x36, 0xbc, 0x02, 0x1a, 0xe6, 0x68, 0xf8, 0xf9, 0x52, 0xa0, 0x77, 0xe8, 0x91, 0x19, 0x83,
	0x98, 0xa7, 0xd6, 0x5d, 0x93, 0x9a, 0xa7, 0x89, 0xac, 0xf9, 0x39, 0x66, 0x56, 0x92, 0x26, 0x75,
	0xa0, 0xf7, 0x3b, 0x00, 0xd8, 0x4f, 0xd3, 0xa2, 0x3b, 0x9e, 0x16, 0x4c, 0x69, 0xd5, 0x1b, 0x52,
	0xbb, 0xee, 0x86, 0xd4, 0xaf, 0xf6, 0xab, 0x10, 0x34, 0xac, 0x08, 0x3a, 0x1d, 0x97, 0x68, 0x66,
	0x5c, 0x08, 0x84, 0xda, 0x14, 0xd1, 0xb0, 0x4e, 0xec, 0xbb, 0x77, 0x02, 0xed, 0x03, 0xa6, 0xe3,
	0xd3, 0x22, 0xb7, 0x97, 0xd0, 0x1e, 0x1a, 0x2c, 0x1d, 0xf6, 0x7b, 0x7d, 0xdf, 0xeb, 0x30, 0x2d,
	0x82, 0xce, 0x98, 0x95, 0xb2, 0xd5, 0x16, 0xc9, 0xf6, 0x33, 0x80, 0xd6, 0x07, 0xc6, 0x8e, 0x8b,
	0x38, 0x1b, 0xd0, 0x10, 0x4c, 0xb2, 0x89, 0xf2, 0xf3, 0xe1, 0x91, 0x11, 0x41, 0x27, 0x13, 0xe4,
	0xb9, 0x6b, 0x44, 0x9d, 0x16, 0xd0, 0x30, 0x30, 0x1b, 0x9b, 0x89, 0x72, 0x4d, 0xf0, 0xc8, 0xf4,
	0xee, 0x24, 0xcf, 0x62, 0xdb, 0x1e, 0xd7, 0xf5, 0x12, 0x4f, 0xe5, 0x88, 0x2a, 0x72, 0xec, 0xfe,
	0xad, 0x01, 0x1c, 0x0b, 0x16, 0xe3, 0x61, 0xca, 0xf3, 0x11, 0xd9, 0x81, 0x86, 0xbb, 0xe8, 0x64,
	0xdd, 0xe7, 0x3d, 0x73, 0xe0, 0xbb, 0xab, 0xe5, 0x79, 0x72, 0x63, 0xd7, 0xbb, 0x47, 0x9e, 0x40,
	0x68, 0x8e, 0x15, 0xa9, 0x5e, 0xd3, 0x6b, 0xcc, 0x77, 0xa0, 0xe1, 0x4e, 0x40, 0x19, 0x61, 0xe6,
	0x22, 0x2c, 0xa0, 0xb8, 0xcd, 0x2d, 0x29, 0x33, 0x8b, 0x3c, 0x8f, 0xb2, 0x07, 0xcd, 0x72, 0x11,
	0xc9, 0xc3, 0x42, 0xb1, 0x2b, 0xab, 0x39, 0x8f, 0xf8, 0x14, 0x22, 0x3b, 0x03, 0xe4, 0x81, 0xff,
	0x56, 0x9d, 0x88, 0x05, 0xe5, 0x1f, 0xb2, 0x34, 0x2d, 0xcb, 0xaf, 0x08, 0x3b, 0xc7, 0xfc, 0xe0,
	0x05, 0x6c, 0xc5, 0x7c, 0x32, 0x50, 0xa6, 0xe5, 0xb9, 0xd0, 0x18, 0x9f, 0xba, 0xf7, 0x57, 0x26,
	0x12, 0x67, 0x7d, 0xb0, 0x3a, 0xd5, 0xe3, 0xc8, 0xfc, 0x71, 0x14, 0x0c, 0x1b, 0xf6, 0xcb, 0xf3,
	0x7f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x83, 0xb1, 0x3c, 0x70, 0x86, 0x07, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// SpaceCloudClient is the client API for SpaceCloud service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type SpaceCloudClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*Response, error)
	Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*Response, error)
	Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*Response, error)
	Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Response, error)
	Aggregate(ctx context.Context, in *AggregateRequest, opts ...grpc.CallOption) (*Response, error)
	Batch(ctx context.Context, in *BatchRequest, opts ...grpc.CallOption) (*Response, error)
	Call(ctx context.Context, in *FaaSRequest, opts ...grpc.CallOption) (*Response, error)
}

type spaceCloudClient struct {
	cc *grpc.ClientConn
}

func NewSpaceCloudClient(cc *grpc.ClientConn) SpaceCloudClient {
	return &spaceCloudClient{cc}
}

func (c *spaceCloudClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.SpaceCloud/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceCloudClient) Read(ctx context.Context, in *ReadRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.SpaceCloud/Read", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceCloudClient) Update(ctx context.Context, in *UpdateRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.SpaceCloud/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceCloudClient) Delete(ctx context.Context, in *DeleteRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.SpaceCloud/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceCloudClient) Aggregate(ctx context.Context, in *AggregateRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.SpaceCloud/Aggregate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceCloudClient) Batch(ctx context.Context, in *BatchRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.SpaceCloud/Batch", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *spaceCloudClient) Call(ctx context.Context, in *FaaSRequest, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/proto.SpaceCloud/Call", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SpaceCloudServer is the server API for SpaceCloud service.
type SpaceCloudServer interface {
	Create(context.Context, *CreateRequest) (*Response, error)
	Read(context.Context, *ReadRequest) (*Response, error)
	Update(context.Context, *UpdateRequest) (*Response, error)
	Delete(context.Context, *DeleteRequest) (*Response, error)
	Aggregate(context.Context, *AggregateRequest) (*Response, error)
	Batch(context.Context, *BatchRequest) (*Response, error)
	Call(context.Context, *FaaSRequest) (*Response, error)
}

func RegisterSpaceCloudServer(s *grpc.Server, srv SpaceCloudServer) {
	s.RegisterService(&_SpaceCloud_serviceDesc, srv)
}

func _SpaceCloud_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceCloudServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.SpaceCloud/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceCloudServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceCloud_Read_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReadRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceCloudServer).Read(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.SpaceCloud/Read",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceCloudServer).Read(ctx, req.(*ReadRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceCloud_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceCloudServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.SpaceCloud/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceCloudServer).Update(ctx, req.(*UpdateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceCloud_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceCloudServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.SpaceCloud/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceCloudServer).Delete(ctx, req.(*DeleteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceCloud_Aggregate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AggregateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceCloudServer).Aggregate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.SpaceCloud/Aggregate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceCloudServer).Aggregate(ctx, req.(*AggregateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceCloud_Batch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BatchRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceCloudServer).Batch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.SpaceCloud/Batch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceCloudServer).Batch(ctx, req.(*BatchRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _SpaceCloud_Call_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FaaSRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SpaceCloudServer).Call(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.SpaceCloud/Call",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SpaceCloudServer).Call(ctx, req.(*FaaSRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _SpaceCloud_serviceDesc = grpc.ServiceDesc{
	ServiceName: "proto.SpaceCloud",
	HandlerType: (*SpaceCloudServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _SpaceCloud_Create_Handler,
		},
		{
			MethodName: "Read",
			Handler:    _SpaceCloud_Read_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _SpaceCloud_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _SpaceCloud_Delete_Handler,
		},
		{
			MethodName: "Aggregate",
			Handler:    _SpaceCloud_Aggregate_Handler,
		},
		{
			MethodName: "Batch",
			Handler:    _SpaceCloud_Batch_Handler,
		},
		{
			MethodName: "Call",
			Handler:    _SpaceCloud_Call_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "server.proto",
}
