package main

import (
	// (一部抜粋)
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"

	"google.golang.org/grpc/health"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"

	pb "golang-grpc-starting/genproto/hello"
)

type myServer struct {
	pb.UnimplementedGreetingServiceServer
}

func NewMyServer() *myServer {
	return &myServer{}
}

func (s *myServer) Hello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	// metadata 取り出し
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		log.Println(md)
	}
	// リクエストからnameフィールドを取り出して
	// "Hello, [名前]!"というレスポンスを返す
	log.Println("Hello, " + req.GetName())
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.GetName()),
	}, nil
}

func main() {
	// 1. 8080番portのLisnterを作成
	port := 8080
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		panic(err)
	}

	// 2. gRPCサーバーを作成
	s := grpc.NewServer(
		grpc.UnaryInterceptor(myUnaryServerInterceptor1),
	)

	// 3. gRPCサーバーにGreetingServiceを登録
	pb.RegisterGreetingServiceServer(s, NewMyServer())

	// HealthCheck
	healthSrv := health.NewServer()
	healthpb.RegisterHealthServer(s, healthSrv)
	healthSrv.SetServingStatus("mygrpc", healthpb.HealthCheckResponse_SERVING)

	// 4. サーバーリフレクションの設定
	//   grpcurl のため
	reflection.Register(s)

	// 5. 作成したgRPCサーバーを、8080番ポートで稼働させる
	go func() {
		log.Printf("start gRPC server port: %v", port)
		s.Serve(listener)
	}()

	// 4.Ctrl+Cが入力されたらGraceful shutdownされるようにする
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("stopping gRPC server...")
	s.GracefulStop()
}
