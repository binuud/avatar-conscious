package avatarServer

import (
	"context"
	"errors"
	"io"
	"net/http"

	"github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	"github.com/gorilla/websocket"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"
)

var websocketUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true }, // Allow all origins
}

// WebsocketHandler()
// grpc proto generated code does not support bidirectional stream
// this is a websocket proxy, that listens on a path
// and proxies to avataLive server via a grpc client
// #TODO handle all connection/disconnect scenarios
func WebsocketHandler(w http.ResponseWriter, r *http.Request) {
	// Upgrade the HTTP connection to a WebSocket

	logrus.Info("Websocket - Received connection to Avatar")

	ws, err := websocketUpgrader.Upgrade(w, r, nil)
	if err != nil {
		logrus.Printf("\n Failed to upgrade websocket: %v", err)
		return
	}
	defer ws.Close()

	// Initialize gRPC Client Conn pointing to your running gRPC server (from step 1)
	conn, err := grpc.NewClient("localhost:9090", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logrus.Fatalf("Could not connect to internal gRPC system: %v", err)
	}
	defer conn.Close()
	grpcClient := avatar.NewAvatarServiceClient(conn)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	grpcStream, err := grpcClient.Control(ctx)
	if err != nil {
		logrus.Errorf("Failed to establish gRPC stream connection: %v", err)
		return
	}

	// Channel to coordinate clean connection shutdown between the two loops
	done := make(chan struct{})
	// logrus.Infof("Handling request %v", r)

	// Loop A: Read from gRPC Stream -> Write to WebSocket Client
	go func() {
		defer close(done)
		for {
			grpcMsg, err := grpcStream.Recv()
			if errors.Is(err, io.EOF) {
				logrus.Errorf("Backend gRPC server cleanly closed connection stream (EOF)")
				return
			}
			if err != nil {
				logrus.Errorf("Error reading from gRPC server stream: %v", err)
				return
			}
			logrus.Info("Websocket <-- Received data from grpc")
			jsonBytes, err := protojson.Marshal(grpcMsg)
			if err != nil {
				logrus.Errorf("Failed to marshal gRPC proto to JSON: %v", err)
				return
			}

			// Send the raw JSON byte array over the WebSocket as a TextMessage
			err = ws.WriteMessage(websocket.TextMessage, jsonBytes)
			if err != nil {
				logrus.Errorf("WebSocket write error: %v", err)
				return
			}

		}
	}()

	// Loop B: Read from WebSocket Client -> Write to gRPC Stream
	for {

		var clientMsg avatar.ControlRequest
		// Read incoming JSON message from ui/client and unpack directly into Proto Struct
		err := readControlRequest(ws, &clientMsg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logrus.Errorf("WebSocket connection encountered unexpected error: %v", err)
			} else {
				logrus.Println("WebSocket connection closed normally by ui or client")
			}
			break
		}

		// Forward unmarshaled payload to backend gRPC Stream
		logrus.Infof("Received message on websocket %s --> grpc", clientMsg.Control.Meta.Name)
		if err := grpcStream.Send(&clientMsg); err != nil {
			logrus.Errorf("Failed to send text payload to gRPC service: %v", err)
			break
		}
	}

	// Tell gRPC server we are done writing data and teardown context pipeline
	_ = grpcStream.CloseSend()
	<-done
}

func readControlRequest(ws *websocket.Conn, clientMsg *avatar.ControlRequest) error {

	messageType, payload, err := ws.ReadMessage()
	if err != nil {
		return err
	}

	// (Optional) Ensure it's a text message if you expect JSON strings
	if messageType == websocket.TextMessage {

		err := protojson.Unmarshal(payload, clientMsg)
		if err != nil {
			logrus.Errorf("Failed to unmarshal proto JSON: %v", err)
			return err
		}
	}
	return nil
}
