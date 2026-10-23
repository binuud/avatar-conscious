package avatarServer

import (
	"context"
	"testing"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	"github.com/sirupsen/logrus"
)

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
