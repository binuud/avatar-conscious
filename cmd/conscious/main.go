package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	gw "github.com/binuud/conscious/gen/go/v1/brain"
	"github.com/binuud/conscious/pkg/servers/consciousServer"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	server_grpc_port = flag.Int("grpc_port", 9090, "The Brain Backend Server GRPC port, no token required unsecured")
	server_http_port = flag.Int("http_port", 9080, "The Brain Backend Server HTTP port, no token required unsecured")
)

func runGRPCServer() {

	insecureConn, err := net.Listen("tcp", fmt.Sprintf(":%d", *server_grpc_port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	insecureServer := grpc.NewServer()

	// GRPC App Store Server
	gw.RegisterConsciousServer(insecureServer, consciousServer.NewGrpcConsciousServer())

	// Register reflection service on gRPC server.
	reflection.Register(insecureServer)

	logrus.Info("GRPC server starting...")
	logrus.Fatalln(insecureServer.Serve(insecureConn))

}

func main() {
	logrus.Println("Welcome to life")

	LLM_ENDPOINT := os.Getenv("LLM_ENDPOINT")
	if LLM_ENDPOINT == "" {
		LLM_ENDPOINT = "localhost:12434"
	}

	logrus.Printf("LLM Endpoint %v", LLM_ENDPOINT)
	runGRPCServer()
	// consciousServer.MakeLLMRequest()
}
