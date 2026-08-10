package consciousServer

import (
	"encoding/gob"
	"os"

	datalayer "github.com/binuud/conscious/pkg/dataLayer"
	"github.com/sirupsen/logrus"
)

func (m *memory) DumpMemory() {
	logrus.Info("-----------------------------------------------------------------")
	logrus.Infof("Name: %s, Model: %s, Temp: %d", m.Name, m.Model, m.Temperature)
	logrus.Infof("System Prompt: %s", m.SystemPrompt)
	for _, msg := range m.Messages {
		var role string
		var content string
		switch {
		case msg.OfUser != nil:
			role = "User"
			content = msg.OfUser.Content.OfString.Value
		case msg.OfAssistant != nil:
			role = "Assistant"
			content = msg.OfAssistant.Content.OfString.Value
		case msg.OfSystem != nil:
			role = "System"
			content = msg.OfSystem.Content.OfString.Value
		case msg.OfDeveloper != nil:
			role = "Developer"
			content = msg.OfDeveloper.Content.OfString.Value
		case msg.OfFunction != nil:
			role = "Function"
			content = msg.OfFunction.Content.Value
		case msg.OfTool != nil:
			role = "Tool"
			content = msg.OfTool.Content.OfString.Value
		}

		logrus.Infof("Name: %s, content: %s", role, content)
	}
	logrus.Info("-----------------------------------------------------------------")
}

// internal function to load from model file
func (m *memory) Load(fileName string) error {

	file, err := os.Open(datalayer.GetFileAbsolutePath(fileName))
	if err != nil {
		return err
	}
	defer file.Close() // Ensure the file is closed

	decoder := gob.NewDecoder(file)
	err = decoder.Decode(m) // Decode into the pointer
	if err != nil {
		return err
	}
	return nil

}

func (m *memory) Save(fileName string) error {

	fullFileName := datalayer.GetFileName(fileName)

	f, err := os.Create(fullFileName)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := gob.NewEncoder(f)

	// save the file and return error if present

	// we need not save tool chain
	tools := m.Tools
	m.Tools = nil
	err = enc.Encode(m)
	m.Tools = tools

	if err != nil {
		logrus.Errorf("Cannot save file %s %v", fullFileName, err)
		return err
	}
	logrus.Infof("Saved memory to file %s", fullFileName)
	return nil
}
