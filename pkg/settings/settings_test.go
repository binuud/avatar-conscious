package settings_test

import (
	"fmt"
	"testing"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/sirupsen/logrus"
)

func Test_Settings(t *testing.T) {
	// Step 1: Setup
	t.Run("Step=CreateSettingsAndSave", func(t *testing.T) {
		_create_settings_and_save(t)
	})

	// Step 2: Depends on Step 1
	t.Run("Step=LoadSettings", func(t *testing.T) {
		_load_settings(t)
	})

	// Step 3: Depends on Step 2
	t.Run("Step=EnvVariables", func(t *testing.T) {
		_env_var_settings(t)
	})

	// Step 3: Depends on Step 3
	t.Run("Step=EnvVariablesShouldNotPersist", func(t *testing.T) {
		_env_var_settings_should_not_persist(t)
	})

}

func _init_settings_entries(s *settings.Settings) {
	s.Add(&pb.Config{
		Name:  "llm-origin",
		Group: "llm",
		Value: "local",
		Type:  pb.Config_STRING,
	})
	s.Add(&pb.Config{
		Name:  "user",
		Group: "llm",
		Value: "test-user",
		Type:  pb.Config_STRING,
	})
	s.Add(&pb.Config{
		Name:  "enabled",
		Group: "llm",
		Value: "true",
		Type:  pb.Config_BOOLEAN,
	})
	s.Add(&pb.Config{
		Name:  "endpoint",
		Group: "database",
		Value: "localhost:9080",
		Type:  pb.Config_STRING,
	})
	s.Add(&pb.Config{
		Name:  "emtpy_test",
		Group: "database",
		Value: "",
		Type:  pb.Config_STRING,
	})
}

func _create_settings_and_save(t *testing.T) {

	t.Logf("Test_Settings_Save")
	s := settings.NewSettings("test-settings")
	_init_settings_entries(s)

	// check if able to get settings by name
	got := s.Get("endpoint")
	if got == nil {
		t.Errorf("error - not able to retreive settings entry")
	} else {
		if got.Group == "database" && got.Value == "localhost:9080" {
			t.Logf("able to retreive settings entry")
		} else {
			t.Errorf("error - not able to retreive settings entry, entry corrupted")
		}
		// update an settings entry by value
		// retreive same and check
		s.Update("endpoint", "localhost:9090")
		got = s.Get("endpoint")
		if got == nil {
			t.Errorf("error - not able to retreive settings entry")
		} else {
			if got.Group == "database" && got.Value == "localhost:9090" {
				t.Logf("able to update settings entry")
			} else {
				t.Errorf("error - not able to update settings entry, entry corrupted")
			}
		}
	}

	err := s.Save()
	if err != nil {
		t.Errorf("error - cannot save settings %v", err)
	} else {
		t.Logf("saved settings file")
	}

}

func _load_settings(t *testing.T) {

	b := settings.NewSettings("test-settings")
	got := b.Get("endpoint")
	if got != nil {
		t.Errorf("error - expected nil, but settings entry present after NewSettings")
	}

	err := b.Load()
	if err != nil {
		t.Errorf("error - cannot load settings %v", err)
	}

	got = b.Get("endpoint")
	if got != nil {
		t.Errorf("error - expected nil, but settings entry present after NewSettings and load")
	}

	_init_settings_entries(b)
	err = b.Load()
	if err != nil {
		t.Errorf("error - cannot load settings %v", err)
	}

	got = b.Get("endpoint")
	if got == nil {
		t.Errorf("error - not able to retreive settings entry, after loading")
	} else {
		if got.Group == "database" {
			t.Logf("able to retreive settings entry after save, data verified")
		} else {
			t.Errorf("error - not able to retreive settings entry, fter loading, entry corrupted")
		}
	}

	groups := b.GetAllByGroup()
	fmt.Printf("%v ", groups)
	if groups == nil {
		t.Errorf("error - not able to retreive settings by groups")
	} else {
		fmt.Printf("%v ", groups)
		if _, ok := groups["llm"]; ok {
			if len(groups["llm"].Configs) == 3 {
				t.Logf("able to retreive settings entry by group for llm")
			}
		} else {
			t.Errorf("error - llm group not found")
		}
		if _, ok := groups["database"]; ok {
			if len(groups["database"].Configs) == 1 {
				t.Logf("able to retreive settings entry by group for database")
			}
		} else {
			t.Errorf("error - database group not found")
		}
	}
}

func _env_var_settings(t *testing.T) {

	s := settings.NewSettings("test-settings")
	_init_settings_entries(s)
	err := s.Load()
	if err != nil {
		t.Errorf("error - cannot load settings %v", err)
	}

	// Try setting a non existing ENV entry in config
	s.SetWithEnv()

	got := s.Get("SETTINGS_IGNORE_TESTCASE_ENDPOINT")
	if got != nil {
		t.Errorf("error - able to retreive settings entry")
	} else {
		t.Logf("Success, not able to set unavailable config entry with env variables")
	}

	// Try setting the config entry, after adding it to config
	s.Add(&pb.Config{
		Name:  "SETTINGS_IGNORE_TESTCASE_ENDPOINT",
		Group: "database",
		Value: "SAMPLE-555",
		Type:  pb.Config_STRING,
	})

	s.SetWithEnv()

	got = s.Get("SETTINGS_IGNORE_TESTCASE_ENDPOINT")
	if got == nil {
		t.Errorf("error - not able to retreive settings entry")
	} else {
		if s.GetValue("SETTINGS_IGNORE_TESTCASE_ENDPOINT") == "Sample123" {
			t.Logf("Success, able to set env entry")
		} else {
			t.Errorf("error - unable to retreive env config entry")
		}
		if got.EnvValue != "Sample123" {
			t.Errorf("error - wrong env value set")
		}
		if got.Value != "SAMPLE-555" {
			t.Errorf("error - value field corrupted by env set")
		}
	}

	s.Save()

}

func _env_var_settings_should_not_persist(t *testing.T) {

	s := settings.NewSettings("test-settings")
	_init_settings_entries(s)
	s.Add(&pb.Config{
		Name:  "SETTINGS_IGNORE_TESTCASE_ENDPOINT",
		Group: "database",
		Value: "SAMPLE-555",
		Type:  pb.Config_STRING,
	})
	err := s.Load()
	if err != nil {
		t.Errorf("error - cannot load settings %v", err)
	}
	logrus.Infof("settings %v", s)
	// now check the enviroment entry
	// this entry should be present but the env value should not be saved
	got := s.Get("SETTINGS_IGNORE_TESTCASE_ENDPOINT")
	if got == nil {
		// the config entry should be there, but value should be empty
		t.Errorf("error - not able to retreive settings entry")
	} else {
		if s.GetValue("SETTINGS_IGNORE_TESTCASE_ENDPOINT") == "Sample123" {
			t.Errorf("error - failed, env variable for persisted")
		}
		if got.EnvValue == "Sample123" {
			t.Errorf("error - env value persisted on save - this should not happen")
		}
		if got.Value != "SAMPLE-555" {
			t.Errorf("error - value field corrupted by env set")
		}
	}
}
