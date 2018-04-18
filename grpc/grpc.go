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
	Register (s *grpc.Server)
}

type HttpGatewayServiceRegistrator interface {
	RegisterGateway (ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error)
}

type Server struct {
	grpcServer *grpc.Server
	address string
	serviceRegistrator ServiceRegistrator
}

type HttpGatewayServer struct {
	address string
	serviceRegistrator HttpGatewayServiceRegistrator
}


func New(address string, sr ServiceRegistrator) *Server {
	return &Server{grpcServer: grpc.NewServer(), address: address, serviceRegistrator: sr}
}

func NewHttpGateway(address string, gsr HttpGatewayServiceRegistrator) *HttpGatewayServer {
	return &HttpGatewayServer{address: address, serviceRegistrator: gsr}
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("gRPC failed to listen on tcp port: %s", err)
	}

	s.serviceRegistrator.Register(s.grpcServer)

	if err := s.grpcServer.Serve(lis); err != nil {
		return fmt.Errorf("gRPC failed to serve: %s", err)
	}
	return nil
}

func (s *HttpGatewayServer) Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	s.serviceRegistrator.RegisterGateway(ctx, mux, s.address, opts)

	if err := http.ListenAndServe(s.address, mux); err != nil {
		return fmt.Errorf("http grpc gateway server failed to listen and serve: %s", err)
	}

	return nil
}