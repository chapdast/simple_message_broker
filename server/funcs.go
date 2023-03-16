package server

import (
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/chapdast/simple_message_broker/proto"
)

func newServer() *Server {
	return &Server{
		channels: make(map[string]chan []byte),
	}
}

func Run(port int64, opts []grpc.ServerOption) error {
	var intenalOpts []grpc.ServerOption
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		return fmt.Errorf("Can't listen on port %d, %s\n", port, err.Error())
	}

	grpcServer := grpc.NewServer(intenalOpts...)
	proto.RegisterMessageBrokerServer(grpcServer, newServer())

	return grpcServer.Serve(lis)

}
