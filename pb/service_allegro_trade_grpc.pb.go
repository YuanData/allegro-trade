// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v3.21.12
// source: service_allegro_trade.proto

package pb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

const (
	AllegroTrade_CreateMember_FullMethodName = "/pb.AllegroTrade/CreateMember"
	AllegroTrade_UpdateMember_FullMethodName = "/pb.AllegroTrade/UpdateMember"
	AllegroTrade_LoginMember_FullMethodName  = "/pb.AllegroTrade/LoginMember"
	AllegroTrade_VerifyEmail_FullMethodName  = "/pb.AllegroTrade/VerifyEmail"
)

// AllegroTradeClient is the client API for AllegroTrade service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AllegroTradeClient interface {
	CreateMember(ctx context.Context, in *CreateMemberRequest, opts ...grpc.CallOption) (*CreateMemberResponse, error)
	UpdateMember(ctx context.Context, in *UpdateMemberRequest, opts ...grpc.CallOption) (*UpdateMemberResponse, error)
	LoginMember(ctx context.Context, in *LoginMemberRequest, opts ...grpc.CallOption) (*LoginMemberResponse, error)
	VerifyEmail(ctx context.Context, in *VerifyEmailRequest, opts ...grpc.CallOption) (*VerifyEmailResponse, error)
}

type allegroTradeClient struct {
	cc grpc.ClientConnInterface
}

func NewAllegroTradeClient(cc grpc.ClientConnInterface) AllegroTradeClient {
	return &allegroTradeClient{cc}
}

func (c *allegroTradeClient) CreateMember(ctx context.Context, in *CreateMemberRequest, opts ...grpc.CallOption) (*CreateMemberResponse, error) {
	out := new(CreateMemberResponse)
	err := c.cc.Invoke(ctx, AllegroTrade_CreateMember_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *allegroTradeClient) UpdateMember(ctx context.Context, in *UpdateMemberRequest, opts ...grpc.CallOption) (*UpdateMemberResponse, error) {
	out := new(UpdateMemberResponse)
	err := c.cc.Invoke(ctx, AllegroTrade_UpdateMember_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *allegroTradeClient) LoginMember(ctx context.Context, in *LoginMemberRequest, opts ...grpc.CallOption) (*LoginMemberResponse, error) {
	out := new(LoginMemberResponse)
	err := c.cc.Invoke(ctx, AllegroTrade_LoginMember_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *allegroTradeClient) VerifyEmail(ctx context.Context, in *VerifyEmailRequest, opts ...grpc.CallOption) (*VerifyEmailResponse, error) {
	out := new(VerifyEmailResponse)
	err := c.cc.Invoke(ctx, AllegroTrade_VerifyEmail_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AllegroTradeServer is the server API for AllegroTrade service.
// All implementations must embed UnimplementedAllegroTradeServer
// for forward compatibility
type AllegroTradeServer interface {
	CreateMember(context.Context, *CreateMemberRequest) (*CreateMemberResponse, error)
	UpdateMember(context.Context, *UpdateMemberRequest) (*UpdateMemberResponse, error)
	LoginMember(context.Context, *LoginMemberRequest) (*LoginMemberResponse, error)
	VerifyEmail(context.Context, *VerifyEmailRequest) (*VerifyEmailResponse, error)
	mustEmbedUnimplementedAllegroTradeServer()
}

// UnimplementedAllegroTradeServer must be embedded to have forward compatible implementations.
type UnimplementedAllegroTradeServer struct {
}

func (UnimplementedAllegroTradeServer) CreateMember(context.Context, *CreateMemberRequest) (*CreateMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateMember not implemented")
}
func (UnimplementedAllegroTradeServer) UpdateMember(context.Context, *UpdateMemberRequest) (*UpdateMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateMember not implemented")
}
func (UnimplementedAllegroTradeServer) LoginMember(context.Context, *LoginMemberRequest) (*LoginMemberResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoginMember not implemented")
}
func (UnimplementedAllegroTradeServer) VerifyEmail(context.Context, *VerifyEmailRequest) (*VerifyEmailResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyEmail not implemented")
}
func (UnimplementedAllegroTradeServer) mustEmbedUnimplementedAllegroTradeServer() {}

// UnsafeAllegroTradeServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AllegroTradeServer will
// result in compilation errors.
type UnsafeAllegroTradeServer interface {
	mustEmbedUnimplementedAllegroTradeServer()
}

func RegisterAllegroTradeServer(s grpc.ServiceRegistrar, srv AllegroTradeServer) {
	s.RegisterService(&AllegroTrade_ServiceDesc, srv)
}

func _AllegroTrade_CreateMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AllegroTradeServer).CreateMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AllegroTrade_CreateMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AllegroTradeServer).CreateMember(ctx, req.(*CreateMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AllegroTrade_UpdateMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AllegroTradeServer).UpdateMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AllegroTrade_UpdateMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AllegroTradeServer).UpdateMember(ctx, req.(*UpdateMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AllegroTrade_LoginMember_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginMemberRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AllegroTradeServer).LoginMember(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AllegroTrade_LoginMember_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AllegroTradeServer).LoginMember(ctx, req.(*LoginMemberRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AllegroTrade_VerifyEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AllegroTradeServer).VerifyEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: AllegroTrade_VerifyEmail_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AllegroTradeServer).VerifyEmail(ctx, req.(*VerifyEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AllegroTrade_ServiceDesc is the grpc.ServiceDesc for AllegroTrade service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AllegroTrade_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "pb.AllegroTrade",
	HandlerType: (*AllegroTradeServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMember",
			Handler:    _AllegroTrade_CreateMember_Handler,
		},
		{
			MethodName: "UpdateMember",
			Handler:    _AllegroTrade_UpdateMember_Handler,
		},
		{
			MethodName: "LoginMember",
			Handler:    _AllegroTrade_LoginMember_Handler,
		},
		{
			MethodName: "VerifyEmail",
			Handler:    _AllegroTrade_VerifyEmail_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "service_allegro_trade.proto",
}
