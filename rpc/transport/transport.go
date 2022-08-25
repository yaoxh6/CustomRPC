package transport

import (
	"bufio"
	"fmt"
	"net"
)

type Package struct {
	ServiceName string
	Data []byte
}

type Transport interface {
	Setup(network string, address string) error
	Recv() (*Package, error)
	Send(pak *Package) error
}

type SimpleTransport struct {
	Client net.Listener
	Conn net.Conn
	isClose bool
	network string
	address string
}

func (s *SimpleTransport) Setup (network string, address string) error {
	var err error
	s.Client, err = net.Listen(network, address)
	if err != nil {
		return err
	}
	s.Conn, err = s.Client.Accept() // 监听客户端的连接请求
	if err != nil {
		fmt.Println("Accept() failed, err: ", err)
		return err
	}
	s.isClose = false
	s.network = network
	s.address = address
	return nil
}

func (s *SimpleTransport) Recv() (*Package, error)  {
	var err error
	if s.isClose {
		s.Setup(s.network, s.address)
	}
	reader := bufio.NewReader(s.Conn)
	var buf [128]byte
	n, err := reader.Read(buf[:]) // 读取数据
	if !s.isClose && err != nil {
		fmt.Println("read from client failed, err: ", err)
		s.Conn.Close()
		s.Client.Close()
		s.isClose = true
		return s.Recv()
	} else if s.isClose && err != nil {
		return nil, err
	}
	recvStr := buf[:n]

	res := &Package {
		ServiceName: "",
		Data:        recvStr,
	}
	return res, nil
}

func (s *SimpleTransport) Send (pak *Package) error {
	_, err := s.Conn.Write(pak.Data)
	if err != nil {
		return err
	}
	return nil
}