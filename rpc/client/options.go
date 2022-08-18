package client

import (
	"time"
)

type Option func(o *Options)

type Options struct {
	CallOptions CallOptions
}

type CallOptions struct {
	Target          string
	TargetHashKey   *uint64
	ReqIdShift      int
	CallTimeout     time.Duration
}

const (
	OptionsDefaultTimeout = 5 * time.Second
)

func NewOptions(options ...Option) Options {
	opts := Options{
		CallOptions: CallOptions{
			Target:          "",
			TargetHashKey:   nil,
			ReqIdShift:      0,
			CallTimeout:     OptionsDefaultTimeout,
		},
	}

	for _, o := range options {
		o(&opts)
	}
	return opts
}