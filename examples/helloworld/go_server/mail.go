package main

import (
	"context"
	log "github.com/hyahm/golog"

	pb "github.com/yaoxh6/CustomRPC/examples/helloworld/helloworld"
)

type s2s struct{}

func (s *s2s) SayHello(ctx context.Context, helloRequest *pb.HelloRequest) (*pb.HelloReply, error) {
	log.Info(helloRequest)
	return &pb.HelloReply{Message: "HelloReplyContent"}, nil
}

func (s *s2s) SayHello2(ctx context.Context, helloRequest *pb.HelloRequest2) (*pb.HelloReply2, error) {
	log.Info(helloRequest)
	return &pb.HelloReply2{ReplyNum:helloRequest.Num, Res: true}, nil
}