package settings

import (
	pbSettings "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	"github.com/binuud/avatar-conscious/pkg/global"
	"github.com/sirupsen/logrus"
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

	// set global settings
	Global.Add(&pbSettings.Config{
		Name:        "LLM_ENDPOINT",
		Group:       "LLM",
		Default:     "http://" + localhost + ":11234/v1",
		Description: "LLM Endpoint, can be local or cloud, no trailing slash",
		Type:        pbSettings.Config_STRING,
	})

	Global.Add(&pbSettings.Config{
		Name:        "LLM_API_KEY",
		Group:       "LLM",
		Default:     "dummy",
		Description: "LLM API Key, use use third party llm providers",
		Type:        pbSettings.Config_STRING,
	})

	Global.Add(&pbSettings.Config{
		Name:        "MQTT_ENDPOINT",
		Group:       "EVENTS",
		Default:     "tcp://" + localhost + ":1883",
		Description: "Default mqtt broker used for messaging and events",
		Type:        pbSettings.Config_STRING,
	})

	Global.Add(&pbSettings.Config{
		Name:        "LLM_MODEL",
		Group:       "MODELS",
		Default:     "lmstudio-community/Qwen3.5-9B-MLX-4bit",
		Description: "Default model used for Chat Generation, Avatars, Tools and Agents",
		Type:        pbSettings.Config_STRING,
	})

	Global.Add(&pbSettings.Config{
		Name:        "VISION_MODEL",
		Group:       "MODELS",
		Default:     "lmstudio-community/Qwen3.5-9B-MLX-4bit",
		Description: "Default model used for Image Description Query",
		Type:        pbSettings.Config_STRING,
	})

	Global.Add(&pbSettings.Config{
		Name:        "IMAGE_MODEL",
		Group:       "MODELS",
		Default:     "Runpod/FLUX.2-klein-4B-mflux-4bit",
		Description: "Default model used for Image Generation",
		Type:        pbSettings.Config_STRING,
	})

	Global.Add(&pbSettings.Config{
		Name:        "TTS_MODEL",
		Group:       "MODELS",
		Default:     "mlx-community/Qwen3-TTS-12Hz-0.6B-Base-8bit",
		Description: "Default Text To Speech model",
		Type:        pbSettings.Config_STRING,
	})

	Global.Add(&pbSettings.Config{
		Name:        "SERVABLE_DIRECTORY",
		Group:       "SERVER",
		Default:     "/dataVolume/static",
		Description: "The folder containing assets, which can be served by the static http server",
		Type:        pbSettings.Config_STRING,
	})

	Global.Add(&pbSettings.Config{
		Name:        "FILE_STORE_PATH",
		Group:       "SERVER",
		Default:     "/dataVolume",
		Description: "The folder containing assets, which can be served by the static http server",
		Type:        pbSettings.Config_STRING,
	})

	// load from storage, if file exists
	Global.Load()

	// override settings with env variable
	Global.SetWithEnv()

	Global.Save()
}
