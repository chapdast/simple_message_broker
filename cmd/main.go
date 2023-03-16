package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/chapdast/simple_message_broker/client"
	"github.com/chapdast/simple_message_broker/proto"
	"github.com/chapdast/simple_message_broker/server"
)

func main() {

	fRunServer := flag.Bool("server", false, "Run as a Server")
	fRunServerShort := flag.Bool("s", false, "Short for --server")
	fPort := flag.Int64("port", 19601, "Port Number")

	//client specific
	fAddress := flag.String("server_address", "localhost:1960", "destination server address")

	flag.Parse()
	port := *fPort
	asServer := *fRunServer || *fRunServerShort

	// RUN AS SERVER
	if asServer {
		if err := server.Run(port, []grpc.ServerOption{
			grpc.UnaryInterceptor(
				func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
					fmt.Println("SERVER", req, info)
					return handler(ctx, req)
				},
			),
		}); err != nil {
			log.Fatalf("error running server on port %d,\n%s\n", port, err.Error())
		}
		os.Exit(0)
	}

	//RUN AS CLIENT
	cli, closer, err := client.New(*fAddress, []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer closer()

	ctx := context.Background()
	cname := "chapdast@public"
	if err := cli.Create(ctx, cname); err != nil {
		log.Fatal(err)
	}
	outCtx, canceled := context.WithCancel(ctx)
	defer canceled()
	out := make(chan *proto.SubscribeResponse, 1)
	go func() {
		for {
			select {
			case <-outCtx.Done():
				return
			case m := <-out:
				fmt.Printf("message: %s", m.Message)
			}
		}
	}()

	if err := cli.Publish(ctx, cname, []byte("Hello Message")); err != nil {
		log.Fatal(err)
	}

	if err := cli.Subscribe(ctx, cname, out); err != nil {
		log.Fatal(err)
		canceled()
	}

	if err := cli.Publish(ctx, cname, []byte("Hello Message")); err != nil {
		log.Fatal(err)
	}
	// http.ListenAndServe(fmt.Sprintf("localhost:%d", port), )

}
