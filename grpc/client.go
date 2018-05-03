package grpc

import (
	"google.golang.org/grpc"
)

type Client struct {
	clientConn *grpc.ClientConn
	address string
}

func NewClient(address string) (*Client, error) {
	clientConn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return &Client{clientConn: clientConn, address: address}, nil
}

func (c *Client) Close() {
	c.clientConn.Close()
}

func (c *Client) Connection() *grpc.ClientConn {
	return c.clientConn
}

func (c* Client) Address() string {
	return c.address
}


