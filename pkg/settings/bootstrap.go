package settings

import (
	pbSettings "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	"github.com/binuud/avatar-conscious/pkg/global"
	"github.com/sirupsen/logrus"
)

const (
	SYSTEM_PROMPT = "You are a helpful assistant." +
		"You must use the tools provided to answer questions." +
		"Very polite, friendly with children. Do not use curse words." +
		"Do not use bad words. Do not use any smilies, emojis, or symbols in your responses from now on." +
		"If tool call response if a filename or a filename with path, give output as \"\render [filename]\""
)

func Bootstrap() {

	// override defaults localhost
	// when run from within docker
	localhost := "localhost"
	if global.IsDocker {
		localhost = "host.docker.internal"
		logrus.Infof("build for docker, using localhost : %s", localhost)
	} else {
		logrus.Infof("build for non docker, using localhost : %s", localhost)
	}
	// initialize the global settings
	Global = NewSettings("global")

	Global.Add(&pbSettings.Config{
		Name:        "INITIAL_CONFIGURED",
		Group:       "_HIDDEN",
		Default:     "NO",
		Description: "Is the initial configuration done",
		Type:        pbSettings.Config_STRING,
		Weight:      11,
	})

	Global.Add(&pbSettings.Config{
		Name:        "LLM_PROVIDER",
		Group:       "LLM",
		Default:     "MLX-SERVE",
		Description: "Select the LLM Provider",
		Options:     []string{"CLAUDE", "DOCKER-MODEL-RUNNER", "CHATGPT", "MLX-SERVE", "OLLAMA", "OPENROUTER", "OTHERS"},
		Type:        pbSettings.Config_MULTI_OPTION,
		Weight:      10,
	})

	// set global settings
	Global.Add(&pbSettings.Config{
		Name:        "LLM_ENDPOINT",
		Group:       "LLM",
		Default:     "http://" + localhost + ":11234/v1",
		Description: "LLM Endpoint, can be local or cloud, no trailing slash",
		Type:        pbSettings.Config_STRING,
		Weight:      11,
	})

	Global.Add(&pbSettings.Config{
		Name:        "LLM_API_KEY",
		Group:       "LLM",
		Default:     "dummy",
		Description: "LLM API Key, use use third party llm providers",
		Type:        pbSettings.Config_STRING,
		Weight:      12,
	})

	Global.Add(&pbSettings.Config{
		Name:        "LLM_SYSTEM_PROMPT",
		Group:       "LLM",
		Default:     SYSTEM_PROMPT,
		Description: "System Prompt used in chat interface and streaming chat",
		Type:        pbSettings.Config_MULTI_LINE_STRING,
		Weight:      20,
		NumLines:    6,
	})

	Global.Add(&pbSettings.Config{
		Name:        "MQTT_ENDPOINT",
		Group:       "EVENTS",
		Default:     "tcp://" + localhost + ":1883",
		Description: "Default mqtt broker used for messaging and events",
		Type:        pbSettings.Config_STRING,
		Weight:      11,
	})

	Global.Add(&pbSettings.Config{
		Name:        "LLM_MODEL",
		Group:       "MODELS",
		Default:     "lmstudio-community/Qwen3.5-9B-MLX-4bit",
		Description: "Default model used for Chat Generation, Avatars, Tools and Agents",
		Options:     []string{},
		Type:        pbSettings.Config_DATALIST,
		Weight:      11,
	})

	Global.Add(&pbSettings.Config{
		Name:        "VISION_MODEL",
		Group:       "MODELS",
		Default:     "lmstudio-community/Qwen3.5-9B-MLX-4bit",
		Description: "Default model used for Image Description Query",
		Options:     []string{},
		Type:        pbSettings.Config_DATALIST,
		Weight:      12,
	})

	Global.Add(&pbSettings.Config{
		Name:        "IMAGE_MODEL",
		Group:       "MODELS",
		Default:     "Runpod/FLUX.2-klein-4B-mflux-4bit",
		Description: "Default model used for Image Generation",
		Options:     []string{},
		Type:        pbSettings.Config_DATALIST,
		Weight:      13,
	})

	Global.Add(&pbSettings.Config{
		Name:        "TTS_MODEL",
		Group:       "MODELS",
		Default:     "mlx-community/Qwen3-TTS-12Hz-0.6B-Base-8bit",
		Description: "Default Text To Speech model",
		Options:     []string{},
		Type:        pbSettings.Config_DATALIST,
		Weight:      14,
	})

	Global.Add(&pbSettings.Config{
		Name:        "SERVABLE_DIRECTORY",
		Group:       "SERVER",
		Default:     "./dataVolume/static",
		Description: "The folder containing assets, which can be served by the static http server",
		Type:        pbSettings.Config_STRING,
		Weight:      11,
	})

	Global.Add(&pbSettings.Config{
		Name:        "FILE_STORE_PATH",
		Group:       "SERVER",
		Default:     "./dataVolume",
		Description: "The folder containing assets, which can be served by the static http server",
		Type:        pbSettings.Config_STRING,
		Weight:      12,
	})

	Global.Add(&pbSettings.Config{
		Name:        "AVATAR_STORE_DIRECTORY",
		Group:       "SERVER",
		Default:     "./avatar-store",
		Description: "The folder containing avatar-store md files, cache and static resources",
		Type:        pbSettings.Config_STRING,
		Weight:      21,
	})

	// load from storage, if file exists
	Global.Load()

	// override settings with env variable
	Global.SetWithEnv()

	Global.Save()
}
