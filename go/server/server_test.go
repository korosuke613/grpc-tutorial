package main

import (
	"context"
	"github.com/ymmt2005/grpc-tutorial/go/deepthought"
	"reflect"
	"testing"
	"time"
)

func Test_Boot(t *testing.T) {
	type fields struct {
		UnimplementedComputeServer deepthought.UnimplementedComputeServer
	}
	type args struct {
		req    *deepthought.BootRequest
		second time.Duration
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    deepthought.BootResponse
		wantErr bool
	}{
		{
			name: "I THINK THEREFORE I AM.と返す",
			fields: fields{
				UnimplementedComputeServer: deepthought.UnimplementedComputeServer{},
			},
			args: args{
				req:    &deepthought.BootRequest{},
				second: time.Second * 2,
			},
			want: deepthought.BootResponse{
				Message: "I THINK THEREFORE I AM.",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedComputeServer: tt.fields.UnimplementedComputeServer,
			}
			ctx, cancel := context.WithCancel(context.Background())
			go func(cancel func()) {
				time.Sleep(tt.args.second)
				cancel()
			}(cancel)

			stream := mockCompute_BootServer{
				ctx: ctx,
			}

			if err := s.Boot(tt.args.req, stream); (err != nil) != tt.wantErr {
				t.Errorf("Boot() error = %v, wantErr %v", err, tt.wantErr)
			}
			if result.Message != tt.want.Message {
				t.Errorf("result.Message = %v, want %v", result.Message, tt.want.Message)
			}
		})
	}
}

func Test_Infer(t *testing.T) {
	type fields struct {
		UnimplementedComputeServer deepthought.UnimplementedComputeServer
	}
	type args struct {
		ctx context.Context
		req *deepthought.InferRequest
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *deepthought.InferResponse
		wantErr bool
	}{
		{
			name:   "42を答える",
			fields: fields{UnimplementedComputeServer: deepthought.UnimplementedComputeServer{}},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					go func(cancel func()) {
						time.Sleep(2500 * time.Millisecond)
						cancel()
					}(cancel)
					return ctx
				}(),
				req: &deepthought.InferRequest{Query: "Life"},
			},
			want:    &deepthought.InferResponse{Answer: 42},
			wantErr: false,
		},
		{
			name:   "DeadLineに到達する",
			fields: fields{UnimplementedComputeServer: deepthought.UnimplementedComputeServer{}},
			args: args{
				ctx: func() context.Context {
					ctx, cancel := context.WithCancel(context.Background())
					clientDeadline := time.Now().Add(time.Duration(100) * time.Millisecond)
					ctx, cancel = context.WithDeadline(ctx, clientDeadline)
					go func(cancel func()) {
						time.Sleep(2500 * time.Millisecond)
						cancel()
					}(cancel)

					return ctx
				}(),
				req: &deepthought.InferRequest{Query: "Life"},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Server{
				UnimplementedComputeServer: tt.fields.UnimplementedComputeServer,
			}
			got, err := s.Infer(tt.args.ctx, tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Infer() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Infer() got = %v, want %v", got, tt.want)
			}
		})
	}
}
