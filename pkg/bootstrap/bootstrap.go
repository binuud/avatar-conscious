package bootstrap

import (
	"runtime"

	pbSettings "github.com/binuud/avatar-conscious/gen/go/v1/settings"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/binuud/avatar-conscious/pkg/slashCommands"
	"github.com/binuud/avatar-conscious/pkg/tools"
	"github.com/binuud/avatar-conscious/pkg/utils"
	"github.com/sirupsen/logrus"
)

var __IsDarwin = false
var __IsGpuMetricsAvailable = false

// BootStrap()
// calls Bootstrap in each package as required
// #TODO, check if this can be replaced by init
// #BINU - this gives more control than init, but need to explore init in packages
func Bootstrap() {

	if runtime.GOOS == "darwin" {
		__IsGpuMetricsAvailable = true
		__IsDarwin = true
	}

	logrus.Infof("Runtime OS %s : %s", runtime.GOOS, runtime.GOARCH)
	// initiate all slash commands
	slashCommands.Bootstrap()

	// create a GlobalCommonTools, for LLM to use
	tools.Bootstrap()

	// set global settings
	settings.Global.Set(&pbSettings.Config{
		Name:        "LLM_ENDPOINT",
		Group:       "LLM",
		Value:       utils.GetEnv("LLM_HOST", "http://localhost:11234"),
		Description: "LLM Endpoint, can be local or cloud, no trailing slash",
		Type:        pbSettings.Config_STRING,
	})

	settings.Global.Set(&pbSettings.Config{
		Name:        "LLM_API_KEY",
		Group:       "LLM",
		Value:       utils.GetEnv("LLM_API_KEY", "dummy"),
		Description: "LLM API Key, use use third party llm providers",
		Type:        pbSettings.Config_STRING,
	})

	settings.Global.Set(&pbSettings.Config{
		Name:        "LLM_MODEL",
		Group:       "MODELS",
		Value:       utils.GetEnv("LLM_MODEL", "lmstudio-community/Qwen3.5-9B-MLX-4bit"),
		Description: "Default model used for Chat Generation, Avatars, Tools and Agents",
		Type:        pbSettings.Config_STRING,
	})

	settings.Global.Set(&pbSettings.Config{
		Name:        "VISION_MODEL",
		Group:       "MODELS",
		Value:       utils.GetEnv("VISION_MODEL", "lmstudio-community/Qwen3.5-9B-MLX-4bit"),
		Description: "Default model used for Image Description Query",
		Type:        pbSettings.Config_STRING,
	})

	settings.Global.Set(&pbSettings.Config{
		Name:        "IMAGE_MODEL",
		Group:       "MODELS",
		Value:       utils.GetEnv("IMAGE_MODEL", "Runpod/FLUX.2-klein-4B-mflux-4bit"),
		Description: "Default model used for Image Generation",
		Type:        pbSettings.Config_STRING,
	})

	settings.Global.Set(&pbSettings.Config{
		Name:        "MQTT_ENDPOINT",
		Group:       "EVENTS",
		Value:       utils.GetEnv("MQTT_HOST", "tcp://localhost:1883"),
		Description: "Default mqtt broker used for messaging and events",
		Type:        pbSettings.Config_STRING,
	})

	settings.Global.Save()
}
