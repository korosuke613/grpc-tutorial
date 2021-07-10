package main

import (
	"context"
	"github.com/ymmt2005/grpc-tutorial/go/deepthought"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type mockCompute_BootServer struct {
	ctx context.Context
	deepthought.Compute_BootServer
	grpc.ServerStream
}

var result deepthought.BootResponse

func (m mockCompute_BootServer) SetHeader(metadata.MD) error {
	panic("implement me")
}

func (m mockCompute_BootServer) SetTrailer(metadata.MD) {
	panic("implement me")
}

func (m mockCompute_BootServer) SendHeader(metadata.MD) error {
	panic("implement me")
}

func (m mockCompute_BootServer) Send(response *deepthought.BootResponse) error {
	result = deepthought.BootResponse{
		Message: response.Message,
		Ts:      response.Ts,
	}
	return nil
}

func (m mockCompute_BootServer) Context() context.Context {
	return m.ctx
}

func (m mockCompute_BootServer) SendMsg(i interface{}) error {
	panic("implement me")
}

func (m mockCompute_BootServer) RecvMsg(i interface{}) error {
	panic("impliment me")
}
