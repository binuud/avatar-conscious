package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	gw "github.com/binuud/conscious/gen/go/v1/brain"
	"github.com/binuud/conscious/pkg/servers/consciousServer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

var (
	server_grpc_port = flag.Int("grpc_port", 9090, "The Brain Backend Server GRPC port, no token required unsecured")
	server_http_port = flag.Int("http_port", 9080, "The Brain Backend Server HTTP port, no token required unsecured")
)

func runHtppServer() {

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create a client connection to the gRPC server we just started
	// This is where the gRPC-Gateway proxies the requests
	conn, err := grpc.NewClient(
		fmt.Sprintf("0.0.0.0:%d", *server_grpc_port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.Fatalln("Failed to dial server:", err)
	}

	gwMux := runtime.NewServeMux()

	// AppStore
	err = gw.RegisterConsciousHandler(ctx, gwMux, conn)
	if err != nil {
		log.Fatalln(err)
	}

	// Configure CORS options as needed
	corsEnabledMux := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace "*" with specific origin(s) for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"}, // []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
	})

	corsHandler := corsEnabledMux.Handler(gwMux)

	gwServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", *server_http_port),
		Handler: corsHandler,
	}

	logrus.Infof("HTTP and websocket server starting  :(%d)", *server_http_port)
	// Start HTTP server (and proxy calls to gRPC server endpoint)
	log.Fatalln(gwServer.ListenAndServe())
}

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

	logrus.Infof("GRPC server starting  :(%d)", *server_grpc_port)
	logrus.Fatalln(insecureServer.Serve(insecureConn))

}

func main() {
	logrus.Println("Welcome to life")

	LLM_ENDPOINT := os.Getenv("LLM_ENDPOINT")
	if LLM_ENDPOINT == "" {
		LLM_ENDPOINT = "localhost:12434"
	}

	logrus.Printf("LLM Endpoint %v", LLM_ENDPOINT)
	go func() {
		runGRPCServer()
	}()

	runHtppServer()

}
