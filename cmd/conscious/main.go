package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	pbAvatar "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	pbConscious "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	pbSettings "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	"github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/binuud/avatar-conscious/pkg/servers"
	"github.com/binuud/avatar-conscious/pkg/servers/avatarServer"
	llmProxyServer "github.com/binuud/avatar-conscious/pkg/servers/llmProxyServer"
	"github.com/binuud/avatar-conscious/pkg/servers/staticFileServer"
	"github.com/binuud/avatar-conscious/pkg/settings"

	// mqttServer "github.com/dronasys-com/brain-nodes/pkg/servers/mqttBridge"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// to know that the binary is running within a docker container
var buildEnv = "local" // "local" - for normal builds || "docker" - for docker builds

var (
	server_grpc_port = flag.Int("grpc_port", 9090, "The Brain Backend Server GRPC port, no token required unsecured")
	server_http_port = flag.Int("http_port", 9080, "The Brain Backend Server HTTP port, no token required unsecured")
	verbose          = flag.Bool("v", false, "Enable verbose (debug) output")
)

func serveStaticPaths() *http.ServeMux {

	httpMux := http.NewServeMux()

	// OPEN_API_FILES_PATH is the folder where the swagger files for apps are stored
	SERVABLE_DIRECTORY := settings.Global.GetValue("SERVABLE_DIRECTORY")

	if !strings.Contains(SERVABLE_DIRECTORY, "/dataVolume/static") {
		logrus.Panic("Protection fail, only directory having /dataVolume/static can be served")
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

	// avatar server
	err = pbAvatar.RegisterAvatarServiceHandler(ctx, gwMux, conn)
	if err != nil {
		log.Fatalln(err)
	}

	// systemMetricsServer
	err = pbConscious.RegisterMetricsHandler(ctx, gwMux, conn)
	if err != nil {
		log.Fatalln(err)
	}

	// settings server
	err = pbSettings.RegisterSettingsHandler(ctx, gwMux, conn)
	if err != nil {
		log.Fatalln(err)
	}

	httpMux := serveStaticPaths()

	// for bidirectional stream on avatar live server
	httpMux.HandleFunc("/api/avatar/ws", avatarServer.WebsocketHandler)
	// mqttServer.InitMQTT()

	// handle llm proxt on path /llmHandler
	httpMux.Handle("/llmHandler/", http.StripPrefix("/llmHandler", llmProxyServer.GetLlmProxyHandler()))
	logrus.Infof("Server directory %s -> %s", "/llmHandler/", "llm proxy")

	// handle websocket
	// httpMux.HandleFunc("/ws", mqttServer.WebsocketHandler)

	// link the concious gateway server to httpmux
	httpMux.Handle("/api/", http.StripPrefix("/api", gwMux))
	logrus.Infof("Server directory %s -> %s", "/api/avatar/ws", "avatar control websocket")
	logrus.Infof("Server directory %s -> %s", "/api/", "conscious server")
	logrus.Infof("Server directory %s -> %s", "/api/", "metrics server")
	logrus.Infof("Server directory %s -> %s", "/api/", "avatar store server")
	logrus.Infof("Server directory %s -> %s", "/api/", "settings server")

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

func main() {

	flag.Parse()

	if *verbose {
		logrus.SetLevel(logrus.DebugLevel)
	} else {
		logrus.SetLevel(logrus.InfoLevel)
	}

	flag.Parse()
	logrus.Println("Welcome to life")
	bootstrap.Bootstrap(buildEnv)

	go func() {
		servers.Bootstrap_GRPC_Services(*server_grpc_port)
	}()

	runHtppServer()

}
