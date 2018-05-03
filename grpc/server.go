package grpc

import (
	"net"
	"fmt"
	"net/http"

	"google.golang.org/grpc"
	"golang.org/x/net/context"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
)

type ServiceRegister interface {
	Register(s *grpc.Server)
}

type GatewayServiceRegister interface {
	RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error
}

type Server struct {
	grpcServer *grpc.Server
	address string
}

type HttpGatewayServer struct {
	address string
	grpcEndpointAddress string
	healthCheckEndpoint string

	context context.Context
	cancel context.CancelFunc
	mux *runtime.ServeMux
	opts []grpc.DialOption
}


func NewServer(address string, sr ServiceRegister) *Server {
	grpcServer := grpc.NewServer()
	sr.Register(grpcServer)
	return &Server{grpcServer: grpcServer, address: address}
}

func NewHttpGatewayServer(address, grpcEndpointAddress string, gsr GatewayServiceRegister, healthCheckEndpoint string) (*HttpGatewayServer, error) {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}

	if err := gsr.RegisterGateway(ctx, mux, grpcEndpointAddress, opts); err != nil {
		return nil, fmt.Errorf("failed to register http grpc gateway server, %s", err)
	}

	return &HttpGatewayServer{address: address, grpcEndpointAddress: grpcEndpointAddress, healthCheckEndpoint: healthCheckEndpoint, context: ctx, cancel: cancel, mux: mux, opts: opts}, nil
}

func (s* Server) Address() string {
	return s.address
}

func (s *Server) Run() error {
	lis, err := net.Listen("tcp", s.address)
	if err != nil {
		return fmt.Errorf("gRPC failed to listen on tcp port: %s", err)
	}

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

func (s *HttpGatewayServer) HealthCheck() error {
	req, err := http.NewRequest("GET", fmt.Sprintf("http://%s%s", s.address, s.healthCheckEndpoint), nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("response status: %s", resp.Status)
	}
	return nil
}

func (s *HttpGatewayServer) Run() error {
	defer s.cancel()

	if err := http.ListenAndServe(s.address, s.mux); err != nil {
		return fmt.Errorf("http grpc gateway server failed to listen and serve: %s", err)
	}

	return nil
}

