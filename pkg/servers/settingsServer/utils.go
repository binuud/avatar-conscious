package settingsServer

import (
	"github.com/binuud/avatar-conscious/pkg/llm"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/sirupsen/logrus"
)

func prefillFields(settings *settings.Settings) error {
	models, err := llm.MakeLLMGetModels()
	if err != nil {
		logrus.Errorf("Cannot get models %v", err)
		return err
	}
	modelList := []string{}
	for _, model := range models {
		modelList = append(modelList, model.Id)
	}
	logrus.Info("Prefilled Global Settings Model list")

	llmEntry := settings.Get("LLM_MODEL")
	if llmEntry != nil {
		llmEntry.Options = modelList
	}

	visionEntry := settings.Get("VISION_MODEL")
	if visionEntry != nil {
		visionEntry.Options = modelList
	}

	imageEntry := settings.Get("IMAGE_MODEL")
	if imageEntry != nil {
		imageEntry.Options = modelList
	}

	ttsEntry := settings.Get("TTS_MODEL")
	if ttsEntry != nil {
		ttsEntry.Options = modelList
	}

	return nil
}
