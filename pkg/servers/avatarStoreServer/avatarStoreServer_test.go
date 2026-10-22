package avatarStoreServer

import (
	"context"
	"testing"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/avatar"
	"github.com/sirupsen/logrus"
)

func Test_Avatar_Store_List(t *testing.T) {
	t.Logf("Test_Avatar_Store_List")
	s := NewAvatarStoreGrpcServer()

	got, err := s.List(context.Background(), &pb.AvatarStoreListRequest{
		Limit: 100,
	})
	if err != nil {
		t.Errorf("Test_Avatar_Store_List: cannot list avatars from store : %v \n", err)
	}

	if len(got.Avatars) == 0 {
		t.Errorf("Test_Avatar_Store_List: return zero avatars from store : %v \n", err)
	} else {
		t.Logf("listing from avatar store works, got %d files", len(got.Avatars))
		avatar := got.Avatars[0]
		if avatar.Meta.Uuid == "" || avatar.Meta.Name == "" {
			t.Errorf("Test_Avatar_Store_List: return zero avatars from store : %v \n", err)
		} else {
			logrus.Infof("Avatar file %s", avatar.Meta.Description)
		}
	}

}
