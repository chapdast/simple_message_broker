package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

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
	fInteractive := flag.Bool("i", false, "interactive mode")

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

	// CLI MODE

	if *fInteractive {
		for {
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("Enter Command (? for help): ")
			input, err := reader.ReadString('\n')
			if err != nil {
				log.Printf("error %s", err)
				continue
			}
			input = strings.ReplaceAll(input, "\n", "")
			cArgs := strings.Split(input, " ")
			if len(cArgs) == 0 {
				continue
			}
			cmd := strings.ToLower(cArgs[0])
			if cmd == "p" || cmd == "publish" {
				if len(cArgs) < 3 {
					printHelp()
					continue
				}
				if err := cli.Publish(ctx, cArgs[1], []byte(strings.Join(cArgs[2:], " "))); err != nil {
					log.Println(err)
					continue
				}
			} else if cmd == "c" || cmd == "create" {

				if len(cArgs) != 2 {
					printHelp()
					continue
				}
				if err := cli.Create(ctx, cArgs[1]); err != nil {
					log.Println(err)
					continue
				}
			} else if cmd == "s" || cmd == "sub" {
				if len(cArgs) != 2 {
					printHelp()
					continue
				}
				outCtx, canceled := context.WithCancel(ctx)
				out := make(chan *proto.SubscribeResponse, 1)
				go func() {
					for {
						select {
						case <-outCtx.Done():
							return
						case m := <-out:
							fmt.Printf("> %s\n", m.Message)
						}
					}
				}()

				if err := cli.Subscribe(ctx, cArgs[1], out); err != nil {
					canceled()
					log.Println(err)
				}
			} else if cmd == "?" || cmd == "help" {
				printHelp()
			} else {
				fmt.Println("INVALID COMMAND")
				printHelp()
			}
		}
	}
}

func printHelp() {
	fmt.Println("Avaliable Commands Are:\nsub(s) channel_name,\npublish(p) channel_name message,\ncreate(c) channel_name,\nclose(d) channel_name,")

}
