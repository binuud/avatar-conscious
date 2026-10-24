package avatarServer

import (
	"context"
	"io"
	"net"
	"testing"
	"time"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func init() {
	// to use an in memory grpc listener
	// instead of launching server
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	pb.RegisterAvatarServiceServer(s, &avatarGrpcServer{})
	go func() {
		if err := s.Serve(lis); err != nil {
			panic(err)
		}
	}()
}

func Test_AvatarServer_Start(t *testing.T) {
	t.Logf("Test_AvatarServer_Start")

	a := NewAvatarGrpcServer()

	got, err := a.Start(context.Background(), &pb.StartRequest{
		Name: "election-officer",
	})

	if err != nil {
		t.Errorf("Test_AvatarServer_Start: cannot start avatar %v", err)
	}

	logrus.Infof("Got avatar %v", got)

}

func Test_AvatarServer_Control(t *testing.T) {
	t.Logf("Test_AvatarServer_Control with bidirectional streaming")

	ctx, cancel := context.WithTimeout(context.Background(), 150*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(ctx, "bufnet", grpc.WithContextDialer(bufDialer), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Fatalf("Failed to dial bufnet: %v", err)
	}
	defer conn.Close()

	client := pb.NewAvatarServiceClient(conn)
	stream, err := client.Control(ctx)
	if err != nil {
		t.Fatalf("Avatar Connect Stream failed: %v", err)
	}

	// requests := []string{"Hello", "How are you?", "Goodbye"}
	// waitc := make(chan struct{})

	// trying a block stream
	// doing a sequence of send and receive to simulate ui interaction

	// send start request
	req := &pb.ControlRequest{
		Control: &pb.AvatarControl{
			Action: pb.AvatarControl_Run,
			Meta: &pb.AvatarMeta{
				Name: "election-officer",
			},
		},
	}

	err = stream.Send(req)
	if err != nil {
		t.Fatalf("Failed to Send start request message %v", err)
	}

	count := 0
	waitc := make(chan struct{})
	for {
		t.Logf("Running step %d", count)
		resp, err := stream.Recv()
		if err == io.EOF || count > 50 {
			break
		}
		if err != nil {
			t.Fatalf("Unexpected error on Recv during start send : %v", err)
		}
		logrus.Infof("Step run %d, stepStatus: %v status: %v", resp.Control.CurrentStep, resp.Control.StepStatuses[resp.Control.CurrentStep].Status, resp.Control.Status)
		if resp.Control.Status != pb.AvatarControl_Active {
			// avatar run probably over
			logrus.Infof("Avatar run complete")
			close(waitc)
			break
		}
	}

	// Close the send stream
	if err := stream.CloseSend(); err != nil {
		t.Fatalf("Failed to CloseSend: %v", err)
	}

	// // Wait for receiver goroutine to finish
	select {
	case <-waitc:
	// 	// Test completed successfully
	case <-ctx.Done():
		t.Fatal("Test timed out waiting for stream responses")
	}

}
