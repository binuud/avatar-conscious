package avatarServer

import (
	"testing"

	"github.com/sirupsen/logrus"
)

func Test_Avatar_Election_Officer(t *testing.T) {
	t.Log("Test_Avatar")
	a := NewAvatarFromStore("election-officer")
	a.Build()
	a.Save()
	a.Start()
	a.Run()
}

func Test_Avatar_By_Step_Election_Officer(t *testing.T) {
	t.Log("Test_Avatar")
	a := NewAvatarFromStore("election-officer")
	a.Build()
	a.Save()
	a.Start()
	for i, step := range a.Avatar.Steps {
		err := a.RunStep(i, step)
		if err != nil {
			t.Errorf("Test_Avatar_By_Step_Election_Officer failed %v", err)
		}
	}
}

func Test_PreBuilt_Avatar(t *testing.T) {
	t.Log("Test_PreBuilt_Avatar")
	a := NewAvatar("election-officer")
	a.Print(nil)
	a.Start()
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

func Test_Avatar_CountryGuesser(t *testing.T) {
	t.Logf("Test_Avatar_CountryGuesser")
	a := NewAvatar("sample")
	if a == nil {
		t.Errorf("cannot load avatar ReadAvatarFromStore")
	}
	logrus.Infof("Test_Avatar_CountryGuesser avatar data %v", a)
	a.Build()
	logrus.Infof("Build success")
	a.Save()
	logrus.Infof("Avatar saved")
	a.Start()
	a.Run()
}
