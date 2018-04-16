package grpc

import (
	"net"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type ServiceRegistrator interface {
	register (s *grpc.Server)
}

type HttpGatewayServiceRegistrator interface {
	register (ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

type Server struct {
	grpcServer *grpc.Server
	address string
}

func New(address string) *Server {
	return &Server{grpcServer: grpc.NewServer(), address: address}
}


func (s *Server) Run(sr ServiceRegistrator) error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("gRPC failed to listen on tcp port: %s", err)
	}

	sr.register(s.grpcServer)

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("gRPC failed to serve: %s", err)
	}
	return nil
}

func (s *Server) RunHttpGateway(address string, sr HttpGatewayServiceRegistrator) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	sr.register(ctx, mux, s.address, opts)

	if err := http.ListenAndServe(address, mux); err != nil {
		return fmt.Errorf("http grpc gateway server failed to listen and serve: %s", err)
	}

	return nil
}