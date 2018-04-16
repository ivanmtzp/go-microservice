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
	serviceRegistrator ServiceRegistrator
	address string
	httpGatewayServiceRegistrator HttpGatewayServiceRegistrator
	httpGatewayAddress string
}

func NewServer(address string, sr ServiceRegistrator) *Server {
	return &Server{grpcServer: grpc.NewServer(), address: address, serviceRegistrator: sr}
}

func (s *Server) WithHttpGateway(address string, httpGatewaySR HttpGatewayServiceRegistrator) *Server {
	s.httpGatewayAddress = address
	s.httpGatewayServiceRegistrator = httpGatewaySR
	return s
}


func (s *Server)Run() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("gRPC failed to listen on tcp port: %s", err)
	}

	s.serviceRegistrator.register(s.grpcServer)

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("gRPC failed to serve: %s", err)
	}
	return nil
}

func (s *Server)RunHttpGateway() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	s.httpGatewayServiceRegistrator.register(ctx, mux, s.address, opts)

	if err := http.ListenAndServe(s.httpGatewayAddress, mux); err != nil {
		return fmt.Errorf("http grpc gateway server failed to listen and serve: %s", err)
	}

	return nil
}