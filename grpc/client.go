package grpc

import (
	"google.golang.org/grpc"
)

type CreateClientServiceFunc func(clientConn *grpc.ClientConn) interface{}

type Client struct {
	clientConn *grpc.ClientConn
	address string
	service interface{}
}

func NewClient(address string, serviceCreator CreateClientServiceFunc) (*Client, error) {
	clientConn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{clientConn: clientConn, address: address, service: serviceCreator(clientConn)}, nil
}

func (c* Client) Address() string {
	return c.address
}

func (c* Client) Service() interface{} {
	return c.service
}

func (c *Client) Close() {
	c.clientConn.Close()
}



