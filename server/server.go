package server

import (
	"context"
	"fmt"

	"github.com/chapdast/simple_message_broker/proto"
)

type channelData struct {
	channel     chan []byte
	subscribers uint32
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
	s.channels[req.ChannelId].subscribers++
	for {
		select {
		case m := <-s.channels[req.ChannelId].channel:
			srv.SendMsg(m)
		}
	}
	s.channels[req.ChannelId].subscribers--
	return nil
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
