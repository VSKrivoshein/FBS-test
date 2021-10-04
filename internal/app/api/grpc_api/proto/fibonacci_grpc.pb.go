// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package api

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

// FibonacciClient is the client API for Fibonacci service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type FibonacciClient interface {
	CalcFibonacciSequence(ctx context.Context, in *CalcFibonacciSequenceReq, opts ...grpc.CallOption) (*CalcFibonacciSequenceResponse, error)
}

type fibonacciClient struct {
	cc grpc.ClientConnInterface
}

func NewFibonacciClient(cc grpc.ClientConnInterface) FibonacciClient {
	return &fibonacciClient{cc}
}

func (c *fibonacciClient) CalcFibonacciSequence(ctx context.Context, in *CalcFibonacciSequenceReq, opts ...grpc.CallOption) (*CalcFibonacciSequenceResponse, error) {
	out := new(CalcFibonacciSequenceResponse)
	err := c.cc.Invoke(ctx, "/Fibonacci/CalcFibonacciSequence", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// FibonacciServer is the server API for Fibonacci service.
// All implementations must embed UnimplementedFibonacciServer
// for forward compatibility
type FibonacciServer interface {
	CalcFibonacciSequence(context.Context, *CalcFibonacciSequenceReq) (*CalcFibonacciSequenceResponse, error)
	mustEmbedUnimplementedFibonacciServer()
}

// UnimplementedFibonacciServer must be embedded to have forward compatible implementations.
type UnimplementedFibonacciServer struct {
}

func (UnimplementedFibonacciServer) CalcFibonacciSequence(context.Context, *CalcFibonacciSequenceReq) (*CalcFibonacciSequenceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalcFibonacciSequence not implemented")
}
func (UnimplementedFibonacciServer) mustEmbedUnimplementedFibonacciServer() {}

// UnsafeFibonacciServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to FibonacciServer will
// result in compilation errors.
type UnsafeFibonacciServer interface {
	mustEmbedUnimplementedFibonacciServer()
}

func RegisterFibonacciServer(s grpc.ServiceRegistrar, srv FibonacciServer) {
	s.RegisterService(&Fibonacci_ServiceDesc, srv)
}

func _Fibonacci_CalcFibonacciSequence_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalcFibonacciSequenceReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(FibonacciServer).CalcFibonacciSequence(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Fibonacci/CalcFibonacciSequence",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(FibonacciServer).CalcFibonacciSequence(ctx, req.(*CalcFibonacciSequenceReq))
	}
	return interceptor(ctx, in, info, handler)
}

// Fibonacci_ServiceDesc is the grpc.ServiceDesc for Fibonacci service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Fibonacci_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Fibonacci",
	HandlerType: (*FibonacciServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CalcFibonacciSequence",
			Handler:    _Fibonacci_CalcFibonacciSequence_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "fibonacci.proto",
}
