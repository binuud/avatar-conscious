package bootstrap

import (
	"runtime"

	"github.com/binuud/avatar-conscious/pkg/slashCommands"
)

var __IsDarwin = false
var __IsGpuMetricsAvailable = false

func Bootstrap() {

	if runtime.GOOS == "darwin" {
		__IsGpuMetricsAvailable = true
		__IsDarwin = true
	}

	slashCommands.Bootstrap()
}
