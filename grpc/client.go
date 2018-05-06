package grpc

import (
	"google.golang.org/grpc"
)

type CreateClientServiceFunc func(connection *grpc.ClientConn) interface{}

type Client struct {
	connection *grpc.ClientConn
	address string
	service interface{}
}

func NewClient(address string, serviceCreator CreateClientServiceFunc) (*Client, error) {
	connection, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{connection: connection, address: address, service: serviceCreator(connection)}, nil
}

func (c* Client) Address() string {
	return c.address
}

func (c* Client) Service() interface{} {
	return c.service
}

func (c *Client) Close() {
	if c.connection != nil {
		c.connection.Close()
	}
}



