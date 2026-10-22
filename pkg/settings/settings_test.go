package settings

import (
	"fmt"
	"testing"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/settings"
)

func Test_Settings_Save(t *testing.T) {

	t.Logf("Test_Settings_Save")
	s := NewSettings("test-settings")
	s.Set(&pb.Config{
		Name:  "llm-origin",
		Group: "llm",
		Value: "local",
		Type:  pb.Config_STRING,
	})
	s.Set(&pb.Config{
		Name:  "user",
		Group: "llm",
		Value: "test-user",
		Type:  pb.Config_STRING,
	})
	s.Set(&pb.Config{
		Name:  "enabled",
		Group: "llm",
		Value: "true",
		Type:  pb.Config_BOOLEAN,
	})
	s.Set(&pb.Config{
		Name:  "endpoint",
		Group: "database",
		Value: "localhost:9080",
		Type:  pb.Config_STRING,
	})

	// check if able to get settings by name
	got := s.Get("endpoint")
	if got == nil {
		t.Errorf("not able to retreive settings entry")
	} else {
		if got.Group == "database" && got.Value == "localhost:9080" {
			t.Logf("able to retreive settings entry")
		} else {
			t.Errorf("not able to retreive settings entry, entry corrupted")
		}
		// update an settings entry by value
		// retreive same and check
		s.Update("endpoint", "localhost:9090")
		got = s.Get("endpoint")
		if got == nil {
			t.Errorf("not able to retreive settings entry")
		} else {
			if got.Group == "database" && got.Value == "localhost:9090" {
				t.Logf("able to update settings entry")
			} else {
				t.Errorf("not able to update settings entry, entry corrupted")
			}
		}
	}

	groups := s.GetAllByGroup()
	fmt.Printf("%v ", groups)
	if groups == nil {
		t.Errorf("not able to retreive settings by groups")
	} else {
		fmt.Printf("%v ", groups)
		if _, ok := groups["llm"]; ok {
			if len(groups["llm"].Configs) == 3 {
				t.Logf("able to retreive settings entry by group for llm")
			}
		} else {
			t.Errorf("llm group not found")
		}
		if _, ok := groups["database"]; ok {
			if len(groups["database"].Configs) == 1 {
				t.Logf("able to retreive settings entry by group for database")
			}
		} else {
			t.Errorf("database group not found")
		}
	}

	err := s.Save()
	if err != nil {
		t.Errorf("cannot save settings %v", err)
	}

	b := NewSettings("test-settings")
	err = b.Load()
	if err != nil {
		t.Errorf("cannot load settings %v", err)
	}
	got = b.Get("endpoint")
	if got == nil {
		t.Errorf("not able to retreive settings entry, after loading")
	} else {
		if got.Group == "database" {
			t.Logf("able to retreive settings entry")
		} else {
			t.Errorf("not able to retreive settings entry, fter loading, entry corrupted")
		}
	}

}

func Test_Settings_SaveInternal(t *testing.T) {

	t.Logf("Test_Settings_SaveInternal")
	// TODO
}
