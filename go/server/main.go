package main

import (
	"fmt"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc/keepalive"
	"net"
	"os"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	// protoc で自動生成されたパッケージ
	"github.com/ymmt2005/grpc-tutorial/go/deepthought"
	"google.golang.org/grpc"
)

const portNumber = 13333

func extractFields(fullMethod string, req interface{}) map[string]interface{} {
	ret := make(map[string]interface{})

	switch args := req.(type) {
	case *deepthought.InferRequest:
		ret["Query"] = args.Query
	case *deepthought.BootRequest:
		ret["Silent"] = args.Silent
	default:
		return nil
	}

	return ret
}

func main() {
	kep := keepalive.EnforcementPolicy{
		MinTime: 60 * time.Second,
	}

	log, _ := zap.NewProduction()
	serv := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(kep),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(extractFields)),
				grpc_zap.StreamServerInterceptor(log),
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(extractFields)),
				grpc_zap.UnaryServerInterceptor(log),
			),
		),
	)

	// 実装した Server を登録
	deepthought.RegisterComputeServer(serv, &Server{})

	// 待ち受けソケットを作成
	l, err := net.Listen("tcp", fmt.Sprintf(":%d", portNumber))
	if err != nil {
		fmt.Println("failed to listen:", err)
		os.Exit(1)
	}

	// gRPC サーバーでリクエストの受付を開始
	// l は Close されてから戻るので、main 関数での Close は不要
	serv.Serve(l)
}
