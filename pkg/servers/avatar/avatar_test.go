package avatar

import (
	"testing"
)

func Test_Avatar(t *testing.T) {
	t.Logf("Test_Avatar")
	a := NewAvatar()
	a.Run()

}
