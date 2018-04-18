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
	grpcEndpointAddress string
	serviceRegistrator HttpGatewayServiceRegistrator
}


func New(address string, sr ServiceRegistrator) *Server {
	return &Server{grpcServer: grpc.NewServer(), address: address, serviceRegistrator: sr}
}

func NewHttpGateway(address, grpcEndpointAddress string, gsr HttpGatewayServiceRegistrator) *HttpGatewayServer {
	return &HttpGatewayServer{address: address, grpcEndpointAddress: grpcEndpointAddress, serviceRegistrator: gsr}
}

func (s* Server) Address() string {
	return s.address
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

func (s* HttpGatewayServer) Address() string {
	return s.address
}

func (s* HttpGatewayServer) GrpcEndpointAddress() string {
	return s.grpcEndpointAddress
}


func (s *HttpGatewayServer) Run() error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := s.serviceRegistrator.RegisterGateway(ctx, mux, s.grpcEndpointAddress, opts); err != nil {
		return fmt.Errorf("failed to register http grpc gateway server, %s", err)
	}

	if err := http.ListenAndServe(s.address, mux); err != nil {
		return fmt.Errorf("http grpc gateway server failed to listen and serve: %s", err)
	}

	return nil
}