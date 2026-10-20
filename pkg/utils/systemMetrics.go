package utils

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	pb "github.com/binuud/conscious/gen/go/v1/conscious"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
)

var __IsGpuMetricsAvailable = false

func NewSystemMetrics() (*pb.SystemMetrics, error) {
	// Sample CPU usage over 250ms so we get a meaningful percentage.
	pcts, err := cpu.Percent(250*time.Millisecond, false)
	if err != nil {
		return nil, err
	}
	cpuPct := 0.0
	if len(pcts) > 0 {
		cpuPct = pcts[0]
	}

	vm, err := mem.VirtualMemory()
	if err != nil {
		return nil, err
	}

	ramUsage := (float64(vm.Used) / float64(vm.Total) * 100)

	gpuUsage := 0.0
	if __IsGpuMetricsAvailable {
		gpuUsage, err = GetGpuUsage()
		if err != nil {
			logrus.Fatalf("Cannot parse gpuinfo output: %v\n", err)
		}
	}

	return &pb.SystemMetrics{
		CpuUsage:      cpuPct,
		TotalRamBytes: uint64(vm.Used),
		UsedRamBytes:  uint64(vm.Total),
		RamUsage:      ramUsage,
		GpuUsage:      gpuUsage,
		CreatedAt:     time.Now().Unix(),
	}, nil
}

// GetGpuUsage()
// only works on macos
// finds the total gpuusage on all gpu chips
func GetGpuUsage() (float64, error) {

	// gpu info built from project
	// https://github.com/andersrennermalm/gpuinfo.git
	gpuUsage := 0.0
	cmd := exec.Command("gpuinfo", "-p")
	cmdOut, err := cmd.Output()
	if err != nil {
		// This will catch if the command doesn't exist or fails to run
		return gpuUsage, fmt.Errorf("gpuinfo binary missing, gpuinfo only for macos: %v\n", err)
	} else {
		outputStr := strings.TrimSpace(string(cmdOut))
		gpuUsage, err = strconv.ParseFloat(outputStr, 64)
		if err != nil {
			// unable to parse gpuinfo, file may be corrupted, or version unsupported
			return gpuUsage, fmt.Errorf("Cannot parse gpuinfo output: %v\n", err)
		}
	}
	return gpuUsage, nil

}

// IsGpuMetricsAvailable()
// checks if os is macos, and if gpuinfo binary is available in path
// then it sets variable __IsGpuMetricsAvailable to true and returns true
func IsGpuMetricsAvailable() (isAvailable bool) {
	isAvailable = false
	os := runtime.GOOS

	if os == "darwin" {
		logrus.Infof("Operating system: %s", os)
		gpuUsage, err := GetGpuUsage()
		if err == nil {
			logrus.Infof("CHECK:IsGpuMetricsAvailable: able to get gpuusage on macos %f", gpuUsage)
			isAvailable = true
		} else {
			logrus.Infof("GPU Info is not available")
		}
	}
	__IsGpuMetricsAvailable = isAvailable
	return
}
