package settings

import (
	"fmt"
	"os"
	"strings"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	SETTINGS_DIRECTORY = "configs"
)

var (
	ALLOWED_ENV_VARS = []string{
		"LLM_ENDPOINT",
		"LLM_API_KEY",
		"MQTT_ENDPOINT",
		"LLM_MODEL",
		"IMAGE_MODEL",
		"VISION_MODEL",
		"FILE_STORE_PATH",
		"SERVABLE_DIRECTORY",
		"SETTINGS_IGNORE_TESTCASE_ENDPOINT",
	}
)

type Settings struct {
	*pb.Configs
}

var Global *Settings

func NewSettings(name string) *Settings {
	mySettings := &Settings{
		&pb.Configs{
			Entry: make(map[string]*pb.Config, 0),
			Name:  name,
		},
	}
	return mySettings
}

// // NewLoadSettings
// // return settings object by loading save from storage
// // if storage file exists
// func NewLoadSettings(name string) *settings {
// 	mySettings := &settings{
// 		&pb.Configs{
// 			Entry: make(map[string]*pb.Config, 0),
// 			Name:  name,
// 		},
// 	}
// 	// check if settings file present
// 	// then load the same
// 	err := mySettings.Load()
// 	if err != nil {
// 		logrus.Errorf("config file is missing %v", err)
// 	}
// 	return mySettings
// }

// SetIfNotExists
// The config entry should exist in settings
// If the Value is emtpy, set it with new value
// if not leave the old value untouched
func (s *Settings) SetIfEmtpy(name string, value string) {
	existing := s.Get(name)
	if existing != nil {
		if existing.GetValue() == "" {
			existing.Value = value
		}
	}
}

// SetWithEnv
// Override existing value of the config entry
// with the value provided in environment
// values in environment always overrides config file
func (s *Settings) SetWithEnv() {
	for _, envName := range ALLOWED_ENV_VARS {
		if value, exists := os.LookupEnv(envName); exists {
			s.UpdateEnv(envName, value)
		}
	}
}

func (s *Settings) Add(entry *pb.Config) {
	s.Configs.Entry[entry.Name] = entry
}

func (s *Settings) Update(name string, value string) error {
	if _, ok := s.Configs.Entry[name]; ok {
		s.Configs.Entry[name].Value = value
		return nil
	}
	return fmt.Errorf("no settings entry found for %s", name)
}

func (s *Settings) UpdateEnv(name string, value string) error {
	if _, ok := s.Configs.Entry[name]; ok {
		logrus.Debugf("Config entry override with env %s:%s", name, value)
		s.Configs.Entry[name].EnvValue = value
		return nil
	}
	return fmt.Errorf("no settings entry found for %s", name)
}

// Get()
// gets the config entry by name
// if the entry is absent returns a null value
func (s *Settings) Get(name string) *pb.Config {
	if entry, ok := s.Configs.Entry[name]; ok {
		return entry
	}
	return nil
}

// Get()
// gets the config entry's value by name
// if env is set, return the env first
// if value is set, return value first
// if default is set, return the default value
// if the entry is absent returns a empty string
func (s *Settings) GetValue(name string) string {
	if entry, ok := s.Configs.Entry[name]; ok {
		if entry.EnvValue != "" {
			return entry.EnvValue
		}
		if entry.Value != "" {
			return entry.Value
		}
		return entry.Default
	}
	return ""
}

// GetAllByGroup()
// group all config entries, by their group name
// return the new map, this can be used for rendering
// settings on ui on on cli
func (s *Settings) GetAllByGroup() map[string]*pb.ConfigList {

	var groups = make(map[string]*pb.ConfigList, 0)

	for _, entry := range s.Configs.Entry {
		// logrus.Printf("%s - %s, %s", entry.Group, key, entry.Description)

		if _, ok := groups[entry.Group]; ok {
			groups[entry.Group].Configs = append(groups[entry.Group].Configs, entry)
		} else {
			groups[entry.Group] = &pb.ConfigList{
				Configs: make([]*pb.Config, 0),
			}
			groups[entry.Group].Configs = append(groups[entry.Group].Configs, entry)
		}

	}

	return groups
}

// Save
// Remove all values populated from environment variables
// and save cleaned config
// I do not want env variables to be pushed to repo or leaked
func (s *Settings) Save() error {

	if s.Name == "" || strings.HasPrefix(s.Name, "_") {
		return fmt.Errorf("Internal settings cannot be saved")
	}

	// clone the configs, since all proto messages are by reference
	myClonedMsg := proto.Clone(s.Configs)
	cleanedConfigs := myClonedMsg.(*pb.Configs)
	// emtpt all env entries before saving
	for _, entry := range cleanedConfigs.Entry {
		entry.EnvValue = ""
	}

	configFileName := s.Name + ".config.yaml"
	filePath := datalayer.GetFileAbsolutePath(SETTINGS_DIRECTORY, configFileName)
	// Save list of configs
	if err := datalayer.SaveProtoYAML(filePath, cleanedConfigs); err != nil {
		logrus.Errorf("Cannot save config file %s %v", filePath, err)
		return err
	}
	logrus.Infof("Saved Settings to file %s", configFileName)
	return nil
}

// Load
// Read the config file from storage
// Update the values from storage to the current settings
// Do not add entries, only updation to existing entries happens here
func (s *Settings) Load() error {

	filePath := datalayer.GetFileAbsolutePath(SETTINGS_DIRECTORY, s.Name+".config.yaml")
	newConfigs := &pb.Configs{}

	// note: only the proto msg has to be sent, not the wrapper
	err := datalayer.ReadProtoYAML(filePath, newConfigs)
	if err != nil {
		return err
	}

	for _, entry := range newConfigs.Entry {
		err := s.Update(entry.Name, entry.Value)
		if err != nil {
			logrus.Errorf("config entry in storage, file %s, entry %s is invalid or possible config injection", s.Name, entry.Name)
		}
	}

	return nil
}
