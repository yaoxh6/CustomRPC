package rpc

import (
	"context"
	"github.com/yaoxh6/CustomRPC/rpc/transport"
)

type Service interface {
	Register(serviceDesc interface{}, serviceImpl interface{}) error
	Serve() error
	Close(chan struct{}) error
}

func CreateCustomService(tt transport.Transport, c Codec) (*CustomService, error) {
	h := &CustomService{
		trans: tt,
		d: c,
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.ctx = context.WithValue(ctx, ContextCustomService, h)
	h.cancel = cancel
	return h, nil
}

func Serve(h *CustomService) error {
	return h.Serve()
}
