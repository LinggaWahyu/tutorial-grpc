package main

import (
	"log"
	"net"

	appproto "github.com/Xanvial/tutorial-grpc/proto"
	"github.com/Xanvial/tutorial-grpc/server/handler"
	"github.com/Xanvial/tutorial-grpc/server/interceptor"
	"github.com/Xanvial/tutorial-grpc/server/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// use this product class Usecase as param for grpc handler
	productData := usecase.NewProductUsecase()

	// put grpc class initialization and interceptor here
	productHandler := handler.NewProductHandler(productData)
	grpcInterceptor := interceptor.NewGRPCInterceptor()
	log.Println("productHandler:", productHandler) // just to avoid compile error, remove this after implementing other codes

	// grpc unary interceptors, will be executed from first element to the last
	unaryInterceptors := []grpc.UnaryServerInterceptor{
		grpcInterceptor.LoggingInterceptor(),
		grpcInterceptor.MetadataInterceptor(),
	}

	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(unaryInterceptors...),
	)

	// register server using reflection
	reflection.Register(grpcServer)

	appproto.RegisterProductServiceServer(grpcServer, productHandler)

	log.Println("start listening on port: 9000")
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %s", err)
	}
}
