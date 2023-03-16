package client

import (
	"context"
	"fmt"
	"io"
	"log"

	"google.golang.org/grpc"

	"github.com/chapdast/simple_message_broker/proto"
)

type Client struct {
	cli proto.MessageBrokerClient
}

func New(serverAddress string, opts []grpc.DialOption) (*Client, func(), error) {
	var internalOpts []grpc.DialOption
	if opts != nil {
		internalOpts = append(internalOpts, opts...)
	}
	conn, err := grpc.Dial(serverAddress, internalOpts...)
	if err != nil {
		return nil, nil,
			fmt.Errorf("error dialing to server address, %s\n%s", serverAddress, err.Error())
	}
	cli := proto.NewMessageBrokerClient(conn)
	log.Println("CREATE CLIENT to ", serverAddress)
	return &Client{
			cli: cli,
		}, func() {
			conn.Close()
		}, nil
}

func (c *Client) Create(ctx context.Context, name string) error {
	log.Println("CREATE CHANNEL", name)
	resp, err := c.cli.CreateChannel(ctx, &proto.CreateChannelRequest{
		ChannelId: name,
	})
	if err != nil {
		return err
	}
	if resp.Status != proto.CreateChannelResponse_OK {
		return fmt.Errorf("error %s", resp.Status)
	}
	return nil
}

func (c *Client) Subscribe(ctx context.Context, name string, out chan *proto.SubscribeResponse) error {
	log.Println("SUB CHANNEL", name)
	stream, err := c.cli.Subscribe(ctx, &proto.SubscribeRequest{
		ChannelId: name,
	})
	if err != nil {
		return err
	}
	errChan := make(chan error)
	go func() {
		for {
			in, err := stream.Recv()
			if err != nil && err != io.EOF {
				log.Println("ERR", err)
				errChan <- err
				continue
			}
			if in != nil {
				out <- in
			}
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("user canceled recieving")
		case e := <-errChan:
			return fmt.Errorf("error on messages, %s", e)
		}
	}
}

func (c *Client) Publish(ctx context.Context, name string, message []byte) error {
	log.Println("PUBLISH CHANNEL", name)
	resp, err := c.cli.Publish(ctx, &proto.PublishRequest{
		ChannelId: name,
		Message:   message,
	})
	if err != nil {
		return err
	}

	if resp.Result != proto.ChannelResult_OK {
		return fmt.Errorf("publish error %s", resp.Result)
	}
	return nil
}
