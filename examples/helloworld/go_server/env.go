package main

import (
	"github.com/yaoxh6/CustomRPC/rpc"
	"github.com/yaoxh6/CustomRPC/rpc/transport"
)

func InitServerEnv() (*rpc.CustomService, error) {
	simpleT := transport.SimpleTransport{
		Client: nil,
		Conn:   nil,
	}
	err := simpleT.Setup("tcp", "127.0.0.1:8888")
	if err != nil {
		return nil, err
	}
	codec := rpc.JsonCodec{}
	h, err := rpc.CreateCustomService(&simpleT, &codec)
	if err != nil {
		return nil, err
	}
	return h, nil
}
