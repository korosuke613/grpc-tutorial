package main

import (
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ymmt2005/grpc-tutorial/go/deepthought" // protoc で自動生成されたパッケージ
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"net"
	"net/http"
	"os"
	"time"
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
				grpc_prometheus.StreamServerInterceptor,
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(extractFields)),
				grpc_zap.UnaryServerInterceptor(log),
				grpc_prometheus.UnaryServerInterceptor,
			),
		),
	)

	// 実装した Server を登録
	deepthought.RegisterComputeServer(serv, &Server{})

	// After all your registrations, make sure all of the Prometheus metrics are initialized.
	grpc_prometheus.Register(serv)
	// Register Prometheus metrics handler.
	http.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":8081", nil); err != nil {
			panic(err)
		}
	}()

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
