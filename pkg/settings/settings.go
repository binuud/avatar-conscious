package settings

import (
	"fmt"
	"strings"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	datalayer "github.com/binuud/avatar-conscious/pkg/dataLayer"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

const (
	SETTINGS_DIRECTORY = "configs"
)

type settings struct {
	Configs *pb.Configs
	Name    string
}

var Global = NewSettings("global")

func NewSettings(name string) *settings {
	mySettings := &settings{
		Configs: &pb.Configs{
			Entry: make(map[string]*pb.Config, 0),
		},
		Name: name,
	}
	// check if settings file present
	// then load the same
	mySettings.Load()
	return mySettings
}

func (s *settings) Set(entry *pb.Config) {
	s.Configs.Entry[entry.Name] = entry
}

func (s *settings) Update(name string, value string) error {
	if _, ok := s.Configs.Entry[name]; ok {
		s.Configs.Entry[name].Value = value
		return nil
	}
	return fmt.Errorf("no settings entry found for %s", name)
}

// Get()
// gets the config entry by name
// if the entry is absent returns a null value
func (s *settings) Get(name string) *pb.Config {
	if entry, ok := s.Configs.Entry[name]; ok {
		return entry
	}
	return nil
}

// GetAllByGroup()
// group all config entries, by their group name
// return the new map, this can be used for rendering
// settings on ui on on cli
func (s *settings) GetAllByGroup() map[string]*pb.ConfigList {

	var groups = make(map[string]*pb.ConfigList, 0)

	for key, value := range s.Configs.Entry {
		logrus.Printf("%s - %s, %s", value.Group, key, value.Description)

		if _, ok := groups[value.Group]; ok {
			groups[value.Group].Configs = append(groups[value.Group].Configs, value)
		} else {
			groups[value.Group] = &pb.ConfigList{
				Configs: make([]*pb.Config, 0),
			}
			groups[value.Group].Configs = append(groups[value.Group].Configs, value)
		}

	}

	return groups
}

func (s *settings) Save() error {
	if s.Name == "" || strings.HasPrefix(s.Name, "_") {
		return fmt.Errorf("Internal settings cannot be saved")
	}

	filePath := datalayer.GetFileAbsolutePath(SETTINGS_DIRECTORY, s.Name+".config")
	// Save list of configs
	if err := datalayer.SaveBinary(filePath, s); err != nil {
		logrus.Errorf("Cannot save avatar file %s %v", filePath, err)
		return err
	}
	logrus.Infof("Saved Settings to file %s", filePath)
	return nil
}

func (s *settings) Load() error {
	filePath := datalayer.GetFileAbsolutePath(SETTINGS_DIRECTORY, s.Name+".config")

	// Read AvatarWrapper back (Go automatically infers type T from your variable assignment)
	mySettings, err := datalayer.ReadBinary[settings](filePath)
	if err != nil {
		return err
	}

	// deep clone object, to avoid mutex copy issues
	s.Configs = proto.Clone(mySettings.Configs).(*pb.Configs)
	return nil
}
