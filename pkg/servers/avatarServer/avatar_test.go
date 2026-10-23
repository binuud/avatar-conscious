package avatarServer

import (
	"testing"
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
