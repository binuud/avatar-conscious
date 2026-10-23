package avatarServer

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_Avatar(t *testing.T) {
	t.Logf("Test_Avatar")
	a := NewAvatarFromStore("election-officer.md")
	a.Build()
	a.Save()
	a.Run()
}

func Test_PreBuilt_Avatar(t *testing.T) {
	t.Logf("Test_PreBuilt_Avatar")
	a := NewAvatar("election-officer")
	a.Print(nil)
	a.Run()
}

func Test_Avatar_Load(t *testing.T) {
	t.Logf("Test_Avatar_Load")
	a := NewAvatar("sample")
	if a == nil {
		t.Errorf("cannot load avatar ReadAvatarFromStore")
	}
	logrus.Infof("%v", a)
	a.Print(nil)

}
