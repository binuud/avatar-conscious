package avatarStoreServer_test

import (
	"context"
	"testing"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	"github.com/binuud/avatar-conscious/pkg/bootstrap"
	"github.com/binuud/avatar-conscious/pkg/servers/avatarStoreServer"
	"github.com/sirupsen/logrus"
)

func init() {
	bootstrap.Bootstrap()
	logrus.SetLevel(logrus.DebugLevel)
}

func Test_Avatar_Store_List(t *testing.T) {
	t.Logf("Test_Avatar_Store_List")
	s := avatarStoreServer.NewGrpcAvatarStoreServer()

	got, err := s.List(context.Background(), &pb.AvatarStoreListRequest{
		Limit: 100,
	})
	if err != nil {
		t.Errorf("Test_Avatar_Store_List: cannot list avatars from store : %v", err)
	}

	if len(got.Avatars) == 0 {
		t.Errorf("Test_Avatar_Store_List: return zero avatars from store : %v", err)
	} else {
		t.Logf("listing from avatar store works, got %d files", len(got.Avatars))
		avatar := got.Avatars[0]
		if avatar.Meta.Uuid == "" || avatar.Meta.Name == "" {
			t.Errorf("Test_Avatar_Store_List: return zero avatars from store : %v", err)
		} else {
			logrus.Infof("Avatar file %s", avatar.Meta.Description)
		}
	}

}

func Test_Avatar_Sync(t *testing.T) {
	t.Logf("Test_Avatar_Sync")
	s := avatarStoreServer.NewAvatarStore()
	if s == nil {
		t.Errorf("error Test_Avatar_Sync: cannot create avatarStore")
		return
	}
	err := s.Sync()
	if err != nil {
		t.Errorf("Test_Avatar_Sync: cannot sync avatar : %v ", err)
	}
}
