package main

import (
	"flag"
	"log"
	"os"

	"google.golang.org/grpc"

	"github.com/chapdast/simple_message_broker/server"
)

func main() {

	fRunServer := flag.Bool("server", false, "Run as a Server")
	fRunServerShort := flag.Bool("s", false, "Short for --server")
	fPort := flag.Int64("port", 1960, "Port Number")

	flag.Parse()
	port := *fPort
	asServer := *fRunServer || *fRunServerShort

	if asServer {
		if err := server.Run(port, []grpc.ServerOption{}); err != nil {
			log.Fatalf("error runnig server on port %d,\n%s\n", port, err.Error())
		}
		os.Exit(0)
	}

}
