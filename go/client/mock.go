package main

import (
	"context"
	"github.com/ymmt2005/grpc-tutorial/go/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type mockComputeClient struct {
	deepthought.ComputeClient
	mockBoot  func(ctx context.Context, in *deepthought.BootRequest) (deepthought.Compute_BootClient, error)
	mockInfer func(ctx context.Context, in *deepthought.InferRequest) (*deepthought.InferResponse, error)
}

func (m *mockComputeClient) Boot(ctx context.Context, in *deepthought.BootRequest, opts ...grpc.CallOption) (deepthought.Compute_BootClient, error) {
	return m.mockBoot(ctx, in)
}

func (m *mockComputeClient) Infer(ctx context.Context, in *deepthought.InferRequest, opts ...grpc.CallOption) (*deepthought.InferResponse, error) {
	return m.mockInfer(ctx, in)
}

type mockCompute_BootClient struct {
	grpc.ClientStream
	mockRecv func() (*deepthought.BootResponse, error)
}

func (m *mockCompute_BootClient) Recv() (*deepthought.BootResponse, error) {
	return m.mockRecv()
}
func (m *mockCompute_BootClient) Header() (metadata.MD, error) {
	return nil, nil
}
func (m *mockCompute_BootClient) Trailer() metadata.MD {
	return nil
}
func (m *mockCompute_BootClient) CloseSend() error {
	return nil
}
func (m *mockCompute_BootClient) Context() context.Context {
	return nil
}
func (m *mockCompute_BootClient) SendMsg(i interface{}) error {
	return nil
}
func (m *mockCompute_BootClient) RecvMsg(i interface{}) error {
	return nil
}

type canceledError struct{}

func (c *canceledError) Error() string {
	return ""
}

func (c *canceledError) GRPCStatus() *status.Status {
	return status.New(codes.Canceled, "Timeout")
}
