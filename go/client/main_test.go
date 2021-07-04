package main

import (
	"context"
	"github.com/ymmt2005/grpc-tutorial/go/deepthought"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func Test_callBoot(t *testing.T) {
	type args struct {
		cc deepthought.ComputeClient
		d  time.Duration
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "return hoge and 2021-07-04 22:54:07 +0000 UTC",
			args: args{
				cc: &mockComputeClient{
					mockBoot: func(ctx context.Context, in *deepthought.BootRequest) (deepthought.Compute_BootClient, error) {
						return &mockCompute_BootClient{
							mockRecv: func() (*deepthought.BootResponse, error) {
								select {
								case <-ctx.Done():
									return nil, &canceledError{}
								case <-time.After(2 * time.Second):
									return &deepthought.BootResponse{
										Message: "hoge",
										Ts:      timestamppb.New(time.Date(2021, 7, 4, 22, 54, 07, 0, &time.Location{})),
									}, nil
								}
							},
						}, nil
					},
				},
				d: 3 * time.Second,
			},
			want:    "Boot: hoge, 2021-07-04 22:54:07 +0000 UTC\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := printTester{t: t}
			p.setupPrintTest()
			if err := callBoot(tt.args.cc, tt.args.d); (err != nil) != tt.wantErr {
				t.Errorf("callInfer() error = %v, wantErr %v", err, tt.wantErr)
			}
			actual := p.donePrintTest()
			if actual != tt.want {
				t.Errorf("print() = %s, want %s", actual, tt.want)
			}
		})
	}
}

func Test_callInfer(t *testing.T) {
	type args struct {
		cc deepthought.ComputeClient
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "return 100",
			args: args{
				cc: &mockComputeClient{
					mockInfer: func(ctx context.Context, in *deepthought.InferRequest) (*deepthought.InferResponse, error) {
						return &deepthought.InferResponse{Answer: 100}, nil
					},
				},
			},
			want:    "Infer: 100\n",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := printTester{t: t}
			p.setupPrintTest()
			if err := callInfer(tt.args.cc); (err != nil) != tt.wantErr {
				t.Errorf("callInfer() error = %v, wantErr %v", err, tt.wantErr)
			}
			actual := p.donePrintTest()
			if actual != tt.want {
				t.Errorf("print() = %s, want %s", actual, tt.want)
			}
		})
	}
}
