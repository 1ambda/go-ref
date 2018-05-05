// Code generated by protoc-gen-go. DO NOT EDIT.
// source: location.proto

package pb

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

type LocationContext struct {
	SessionId string `protobuf:"bytes,1,opt,name=sessionId" json:"sessionId,omitempty"`
	Country   string `protobuf:"bytes,2,opt,name=country" json:"country,omitempty"`
}

func (m *LocationContext) Reset()                    { *m = LocationContext{} }
func (m *LocationContext) String() string            { return proto.CompactTextString(m) }
func (*LocationContext) ProtoMessage()               {}
func (*LocationContext) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{0} }

func (m *LocationContext) GetSessionId() string {
	if m != nil {
		return m.SessionId
	}
	return ""
}

func (m *LocationContext) GetCountry() string {
	if m != nil {
		return m.Country
	}
	return ""
}

type LocationRequest struct {
	LocationContext *LocationContext `protobuf:"bytes,1,opt,name=locationContext" json:"locationContext,omitempty"`
}

func (m *LocationRequest) Reset()                    { *m = LocationRequest{} }
func (m *LocationRequest) String() string            { return proto.CompactTextString(m) }
func (*LocationRequest) ProtoMessage()               {}
func (*LocationRequest) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{1} }

func (m *LocationRequest) GetLocationContext() *LocationContext {
	if m != nil {
		return m.LocationContext
	}
	return nil
}

type LocationResponse struct {
	LocationContext *LocationContext `protobuf:"bytes,1,opt,name=locationContext" json:"locationContext,omitempty"`
}

func (m *LocationResponse) Reset()                    { *m = LocationResponse{} }
func (m *LocationResponse) String() string            { return proto.CompactTextString(m) }
func (*LocationResponse) ProtoMessage()               {}
func (*LocationResponse) Descriptor() ([]byte, []int) { return fileDescriptor1, []int{2} }

func (m *LocationResponse) GetLocationContext() *LocationContext {
	if m != nil {
		return m.LocationContext
	}
	return nil
}

func init() {
	proto.RegisterType((*LocationContext)(nil), "pb.LocationContext")
	proto.RegisterType((*LocationRequest)(nil), "pb.LocationRequest")
	proto.RegisterType((*LocationResponse)(nil), "pb.LocationResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Location service

type LocationClient interface {
	AddSession(ctx context.Context, in *LocationRequest, opts ...grpc.CallOption) (*LocationResponse, error)
}

type locationClient struct {
	cc *grpc.ClientConn
}

func NewLocationClient(cc *grpc.ClientConn) LocationClient {
	return &locationClient{cc}
}

func (c *locationClient) AddSession(ctx context.Context, in *LocationRequest, opts ...grpc.CallOption) (*LocationResponse, error) {
	out := new(LocationResponse)
	err := grpc.Invoke(ctx, "/pb.Location/AddSession", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Location service

type LocationServer interface {
	AddSession(context.Context, *LocationRequest) (*LocationResponse, error)
}

func RegisterLocationServer(s *grpc.Server, srv LocationServer) {
	s.RegisterService(&_Location_serviceDesc, srv)
}

func _Location_AddSession_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LocationRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LocationServer).AddSession(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Location/AddSession",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LocationServer).AddSession(ctx, req.(*LocationRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Location_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Location",
	HandlerType: (*LocationServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AddSession",
			Handler:    _Location_AddSession_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "location.proto",
}

func init() { proto.RegisterFile("location.proto", fileDescriptor1) }

var fileDescriptor1 = []byte{
	// 179 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0xcb, 0xc9, 0x4f, 0x4e,
	0x2c, 0xc9, 0xcc, 0xcf, 0xd3, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x2a, 0x48, 0x52, 0xf2,
	0xe4, 0xe2, 0xf7, 0x81, 0x8a, 0x3a, 0xe7, 0xe7, 0x95, 0xa4, 0x56, 0x94, 0x08, 0xc9, 0x70, 0x71,
	0x16, 0xa7, 0x16, 0x17, 0x67, 0xe6, 0xe7, 0x79, 0xa6, 0x48, 0x30, 0x2a, 0x30, 0x6a, 0x70, 0x06,
	0x21, 0x04, 0x84, 0x24, 0xb8, 0xd8, 0x93, 0xf3, 0x4b, 0xf3, 0x4a, 0x8a, 0x2a, 0x25, 0x98, 0xc0,
	0x72, 0x30, 0xae, 0x52, 0x00, 0xc2, 0xa8, 0xa0, 0xd4, 0xc2, 0xd2, 0xd4, 0xe2, 0x12, 0x21, 0x5b,
	0x2e, 0xfe, 0x1c, 0x54, 0xd3, 0xc1, 0x06, 0x72, 0x1b, 0x09, 0xeb, 0x15, 0x24, 0xe9, 0xa1, 0x59,
	0x1c, 0x84, 0xae, 0x56, 0x29, 0x90, 0x4b, 0x00, 0x61, 0x62, 0x71, 0x41, 0x7e, 0x5e, 0x71, 0x2a,
	0x85, 0x46, 0x1a, 0xb9, 0x72, 0x71, 0xc0, 0xd4, 0x08, 0x59, 0x72, 0x71, 0x39, 0xa6, 0xa4, 0x04,
	0x43, 0xbc, 0x26, 0x84, 0xa2, 0x1f, 0xea, 0x01, 0x29, 0x11, 0x54, 0x41, 0x88, 0x1b, 0x94, 0x18,
	0x92, 0xd8, 0xc0, 0x21, 0x68, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0x15, 0xfa, 0x8e, 0x80, 0x53,
	0x01, 0x00, 0x00,
}