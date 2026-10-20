package utils

import "runtime"

var __IsDarwin = false
var __IsGpuMetricsAvailable = false

func Bootstrap() {
	os := runtime.GOOS

	if os == "darwin" {
		__IsGpuMetricsAvailable = true
		__IsDarwin = true
	}

}
