package transport

type Package struct {
	ServiceName string
	Data []byte
}

type Option func(t Transport) error

type Transport interface {
	Poll() (*Package, error)
	Send(pak *Package, opts ...Option) error
}