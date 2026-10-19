package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strings"

	pbAvatar "github.com/binuud/conscious/gen/go/v1/avatar"
	pbConscious "github.com/binuud/conscious/gen/go/v1/conscious"
	"github.com/binuud/conscious/pkg/servers/avatarStoreServer"
	"github.com/binuud/conscious/pkg/servers/consciousServer"
	llmProxyServer "github.com/binuud/conscious/pkg/servers/llmProxyServer"
	"github.com/binuud/conscious/pkg/servers/staticFileServer"
	"github.com/binuud/conscious/pkg/servers/systemMetricsServer"

	"github.com/binuud/conscious/pkg/utils"

	// mqttServer "github.com/dronasys-com/brain-nodes/pkg/servers/mqttBridge"

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

func serveStaticPaths() *http.ServeMux {

	httpMux := http.NewServeMux()

	// OPEN_API_FILES_PATH is the folder where the swagger files for apps are stored
	SERVABLE_DIRECTORY := utils.GetEnv("SERVABLE_DIRECTORY", "/dataVolume/gen")

	if !strings.Contains(SERVABLE_DIRECTORY, "/dataVolume/gen") {
		logrus.Panic("Protection fail, only directory having /dataVolume/gen can be served")
	}

	SERVABLE_UI_DIRECTORY := os.Getenv("SERVABLE_UI_DIRECTORY")
	if SERVABLE_UI_DIRECTORY == "" {
		SERVABLE_UI_DIRECTORY = "/ui"
	}

	if !strings.Contains(SERVABLE_UI_DIRECTORY, "/ui") {
		logrus.Panic("Protection fail, only directory having /ui can be served")
	}

	staticFileServer.StaticServeAiPath(httpMux, "/static/", SERVABLE_DIRECTORY)
	// staticFileServer.SPAStaticServePath(httpMux, "/", SERVABLE_UI_DIRECTORY)

	return httpMux
}

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

	// consciousServer
	err = pbConscious.RegisterConsciousHandler(ctx, gwMux, conn)
	if err != nil {
		log.Fatalln(err)
	}

	// avatar store server
	err = pbAvatar.RegisterAvatarStoreServiceHandler(ctx, gwMux, conn)
	if err != nil {
		log.Fatalln(err)
	}

	// systemMetricsServer
	err = pbConscious.RegisterMetricsHandler(ctx, gwMux, conn)
	if err != nil {
		log.Fatalln(err)
	}

	httpMux := serveStaticPaths()

	// mqttServer.InitMQTT()
	// httpMux.HandleFunc("/ws", mqttServer.WebsocketHandler)

	// handle llm proxt on path /llmHandler
	httpMux.Handle("/llmHandler/", http.StripPrefix("/llmHandler", llmProxyServer.GetLlmProxyHandler()))
	logrus.Infof("Server directory %s -> %s", "/llmHandler/", "llm proxy")

	// link the concious gateway server to httpmux
	httpMux.Handle("/api/", http.StripPrefix("/api", gwMux))
	logrus.Infof("Server directory %s -> %s", "/api/", "conscious server")
	logrus.Infof("Server directory %s -> %s", "/api/", "metrics server")
	logrus.Infof("Server directory %s -> %s", "/api/", "avatar store server")

	// Configure CORS options as needed
	corsEnabledMux := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"}, // Replace "*" with specific origin(s) for production
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"}, // []string{"Authorization", "Content-Type", "Accept"},
		AllowCredentials: true,
		Debug:            true, // Set to true to trace CORS errors in console
	})

	corsEnabledHandler := corsEnabledMux.Handler(httpMux)

	gwServer := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", *server_http_port),
		Handler: corsEnabledHandler,
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

	// GRPC Conscious Server
	pbConscious.RegisterConsciousServer(insecureServer, consciousServer.NewGrpcConsciousServer())

	// GRPC Avatar Store Server
	pbAvatar.RegisterAvatarStoreServiceServer(insecureServer, avatarStoreServer.NewAvatarStoreGrpcServer())

	// GRPC SystemMetrics Server
	pbConscious.RegisterMetricsServer(insecureServer, systemMetricsServer.NewGrpcSystemMetricsServer())

	// Register reflection service on gRPC server.
	reflection.Register(insecureServer)

	logrus.Infof("GRPC server starting  :(%d)", *server_grpc_port)
	logrus.Fatalln(insecureServer.Serve(insecureConn))

}

func main() {
	logrus.Println("Welcome to life")

	LLM_HOST := utils.GetEnv("LLM_HOST", "http://localhost:12434")

	logrus.Printf("LLM Endpoint %v", LLM_HOST)
	go func() {
		runGRPCServer()
	}()

	runHtppServer()

}
