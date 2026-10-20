package main

import (
	"fmt"
	"runtime"

	"github.com/binuud/avatar-conscious/pkg/utils"
)

func main() {
	os := runtime.GOOS
	fmt.Println("Operating system:", os)
	gpuusage, gpuMem, err := utils.GetGpuMetrics()
	if err == nil {
		fmt.Println("GPU usage:", gpuusage)
		fmt.Println("GPU memory:", gpuMem)
	}

}
