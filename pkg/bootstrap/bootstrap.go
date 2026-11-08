package bootstrap

import (
	"runtime"

	"github.com/binuud/avatar-conscious/pkg/global"
	"github.com/binuud/avatar-conscious/pkg/settings"
	"github.com/binuud/avatar-conscious/pkg/slashCommands"
	"github.com/binuud/avatar-conscious/pkg/tools"
	"github.com/sirupsen/logrus"
)

// BootStrap()
// calls Bootstrap in each package as required
// #TODO, check if this can be replaced by init
// #BINU - this gives more control than init, but need to explore init in packages
func Bootstrap(buildEnv ...string) {

	if len(buildEnv) > 0 && buildEnv[0] == "docker" {
		global.IsDocker = true
	}

	if runtime.GOOS == "darwin" {
		global.IsGpuMetricsAvailable = true
		global.IsDarwin = true
	}

	logrus.Infof("Runtime OS %s : %s", runtime.GOOS, runtime.GOARCH)

	// create the required storage directorues
	CreateStorageDirStructure()

	// initialize settings to default values
	// read from env and change config
	settings.Bootstrap()

	// initiate all slash commands
	slashCommands.Bootstrap()

	// create a GlobalCommonTools, for LLM to use
	tools.Bootstrap()

}
