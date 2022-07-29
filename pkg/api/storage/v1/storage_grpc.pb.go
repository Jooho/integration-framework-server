// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.2
// source: proto/v1/storage.proto

package storage

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

// StorageClient is the client API for Storage service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StorageClient interface {
	GetStorageParams(ctx context.Context, in *GetStorageParamsRequest, opts ...grpc.CallOption) (*GetStorageParamResponse, error)
	GetRenderedStorageManifest(ctx context.Context, in *CreateStorageRequest, opts ...grpc.CallOption) (*CreateStorageResponse, error)
}

type storageClient struct {
	cc grpc.ClientConnInterface
}

func NewStorageClient(cc grpc.ClientConnInterface) StorageClient {
	return &storageClient{cc}
}

func (c *storageClient) GetStorageParams(ctx context.Context, in *GetStorageParamsRequest, opts ...grpc.CallOption) (*GetStorageParamResponse, error) {
	out := new(GetStorageParamResponse)
	err := c.cc.Invoke(ctx, "/api.Storage/GetStorageParams", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *storageClient) GetRenderedStorageManifest(ctx context.Context, in *CreateStorageRequest, opts ...grpc.CallOption) (*CreateStorageResponse, error) {
	out := new(CreateStorageResponse)
	err := c.cc.Invoke(ctx, "/api.Storage/GetRenderedStorageManifest", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StorageServer is the server API for Storage service.
// All implementations must embed UnimplementedStorageServer
// for forward compatibility
type StorageServer interface {
	GetStorageParams(context.Context, *GetStorageParamsRequest) (*GetStorageParamResponse, error)
	GetRenderedStorageManifest(context.Context, *CreateStorageRequest) (*CreateStorageResponse, error)
	mustEmbedUnimplementedStorageServer()
}

// UnimplementedStorageServer must be embedded to have forward compatible implementations.
type UnimplementedStorageServer struct {
}

func (UnimplementedStorageServer) GetStorageParams(context.Context, *GetStorageParamsRequest) (*GetStorageParamResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetStorageParams not implemented")
}
func (UnimplementedStorageServer) GetRenderedStorageManifest(context.Context, *CreateStorageRequest) (*CreateStorageResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetRenderedStorageManifest not implemented")
}
func (UnimplementedStorageServer) mustEmbedUnimplementedStorageServer() {}

// UnsafeStorageServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StorageServer will
// result in compilation errors.
type UnsafeStorageServer interface {
	mustEmbedUnimplementedStorageServer()
}

func RegisterStorageServer(s grpc.ServiceRegistrar, srv StorageServer) {
	s.RegisterService(&Storage_ServiceDesc, srv)
}

func _Storage_GetStorageParams_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStorageParamsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).GetStorageParams(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/GetStorageParams",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).GetStorageParams(ctx, req.(*GetStorageParamsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Storage_GetRenderedStorageManifest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateStorageRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StorageServer).GetRenderedStorageManifest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Storage/GetRenderedStorageManifest",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StorageServer).GetRenderedStorageManifest(ctx, req.(*CreateStorageRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Storage_ServiceDesc is the grpc.ServiceDesc for Storage service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Storage_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "api.Storage",
	HandlerType: (*StorageServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetStorageParams",
			Handler:    _Storage_GetStorageParams_Handler,
		},
		{
			MethodName: "GetRenderedStorageManifest",
			Handler:    _Storage_GetRenderedStorageManifest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/v1/storage.proto",
}