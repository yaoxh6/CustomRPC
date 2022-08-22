package rpc

import (
	"context"
	log "github.com/hyahm/golog"

	"github.com/yaoxh6/CustomRPC/rpc/transport"
)

type Service interface {
	Register(serviceDesc interface{}, serviceImpl interface{}) error
	Serve() error
	Close(chan struct{}) error
}

func CreateCustomService(tt transport.Transport) (*CustomService, error) {
	h := &CustomService{
		trans: tt,
	}

	ctx, cancel := context.WithCancel(context.Background())
	h.ctx = context.WithValue(ctx, ContextCustomService, h)
	h.cancel = cancel

	err := h.initService()
	if err != nil {
		log.Error(err.Error())
		return nil, err
	}
	return h, nil
}

func Serve(h *CustomService) error {
	return h.Serve()
}
