package server

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/chapdast/simple_message_broker/proto"
)

type channelData struct {
	channel     chan []byte
	subscribers uint32
	srvs        []proto.MessageBroker_SubscribeServer
}
type Server struct {
	channels map[string]*channelData
	proto.UnimplementedMessageBrokerServer
}

var s proto.MessageBrokerServer = &Server{}

func (s *Server) CreateChannel(ctx context.Context, req *proto.CreateChannelRequest) (
	*proto.CreateChannelResponse, error) {
	if _, ok := s.channels[req.ChannelId]; ok {
		return &proto.CreateChannelResponse{
			Status: proto.CreateChannelResponse_DUPLICATE,
		}, nil
	}
	s.channels[req.ChannelId] = &channelData{
		channel:     make(chan []byte, 10),
		subscribers: 0,
	}
	return &proto.CreateChannelResponse{
		Status: proto.CreateChannelResponse_OK,
	}, nil
}

func (s *Server) Subscribe(req *proto.SubscribeRequest, srv proto.MessageBroker_SubscribeServer) error {
	if _, ok := s.channels[req.ChannelId]; !ok {
		return fmt.Errorf("unknown Channel")
	}

	index := len(s.channels[req.ChannelId].srvs)
	s.channels[req.ChannelId].srvs = append(s.channels[req.ChannelId].srvs, srv)
	s.channels[req.ChannelId].subscribers += 1
	for {
		select {
		case <-srv.Context().Done():
			s.channels[req.ChannelId].subscribers -= 1
			s.channels[req.ChannelId].srvs[index] = nil
			return nil
		case m := <-s.channels[req.ChannelId].channel:
			for _, x := range s.channels[req.ChannelId].srvs {
				if x != nil {
					x.SendMsg(&proto.SubscribeResponse{
						Message: m,
					})
				}
			}

		}
	}
}

func (s *Server) CloseChannel(ctx context.Context, req *proto.CloseChannelRequest) (
	*proto.CloseChannelResponse, error) {

	if _, ok := s.channels[req.ChannelId]; !ok {
		return &proto.CloseChannelResponse{
			Result: proto.ChannelResult_UNKOWN_CHANNEL,
		}, nil
	}
	close(s.channels[req.ChannelId].channel)
	delete(s.channels, req.ChannelId)
	return &proto.CloseChannelResponse{
		Result: proto.ChannelResult_OK,
	}, nil
}

func (s *Server) Publish(ctx context.Context, req *proto.PublishRequest) (
	*proto.PublishResponse, error) {
	if _, ok := s.channels[req.ChannelId]; !ok {
		return &proto.PublishResponse{
			Result: proto.ChannelResult_UNKOWN_CHANNEL,
		}, nil
	}

	go func() {
		s.channels[req.ChannelId].channel <- req.Message
	}()

	return &proto.PublishResponse{
		Result: proto.ChannelResult_OK,
	}, nil
}

func newServer() *Server {
	return &Server{
		channels: make(map[string]*channelData),
	}
}

func Run(port int64, opts []grpc.ServerOption) error {
	var intenalOpts []grpc.ServerOption
	if opts != nil {
		intenalOpts = append(intenalOpts, opts...)
	}
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))

	if err != nil {
		return fmt.Errorf("can't listen on port %d, %s", port, err.Error())
	}

	grpcServer := grpc.NewServer(intenalOpts...)
	proto.RegisterMessageBrokerServer(grpcServer, newServer())
	log.Println("RUNNING SERVER ON: ", lis.Addr())
	return grpcServer.Serve(lis)

}
