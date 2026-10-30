package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"time"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/binuud/avatar-conscious/pkg/global"
	"github.com/sirupsen/logrus"
)

func NewSystemMetrics() (*pb.SystemMetrics, error) {
	// Sample CPU usage over 250ms so we get a meaningful percentage.
	if runtime.GOOS != "darwin" {
		return nil, fmt.Errorf("Metrics not implemented on non macos operating systems")
	}

	if global.IsDocker {
		return nil, fmt.Errorf("Metrics not implemented on docker container")
	}

	cpuUsage, err := GetCpuUsage()
	if err != nil {
		logrus.Fatalf("Cannot parse Cpu metrics output: %v\n", err)
	}

	memUsage, err := GetCpuUsage()
	if err != nil {
		logrus.Fatalf("Cannot parse Cpu metrics output: %v\n", err)
	}

	gpuUsage, gpuMem, err := GetGpuMetrics()
	if err != nil {
		logrus.Fatalf("Cannot parse GetGpuMetrics output: %v\n", err)
	}

	return &pb.SystemMetrics{
		CpuUsage:      cpuUsage,
		TotalRamBytes: 0,
		UsedRamBytes:  0,
		RamUsage:      memUsage,
		GpuUsage:      float64(gpuUsage),
		GpuMemUsed:    int64(gpuMem),
		CreatedAt:     time.Now().Unix(),
	}, nil
}

// GetGpuMetrics()
// runs ioreg to capture the current GPU metrics percentage on M-series chips.
// returns gpusage as percentage, and total memory used by GPU in bytes
func GetGpuMetrics() (gpuUsage int, gpuMem int, err error) {

	gpuUsage = 0
	gpuMem = 0

	// Query the AGXAccelerator service for the performance metrics dictionary
	cmd := exec.Command("ioreg", "-r", "-c", "AGXAccelerator", "-d", "2")
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return gpuUsage, gpuMem, fmt.Errorf("failed running ioreg: %w", err)
	}

	// Regex to extract the "Device Utilization %" integer value
	re := regexp.MustCompile(`"Device Utilization %"\s*=\s*(\d+)`)
	matches := re.FindStringSubmatch(out.String())

	if len(matches) < 2 {
		return gpuUsage, gpuMem, fmt.Errorf("could not parse GPU utilization; confirm this is an Apple Silicon Mac")
	}
	// Convert parsed string back to an integer
	gpuUsage, err = strconv.Atoi(matches[1])
	if err != nil {
		return
	}

	// Regex to extract the "Device Utilization %" integer value
	re = regexp.MustCompile(`"In use system memory"\s*=\s*(\d+)`)
	matches = re.FindStringSubmatch(out.String())

	if len(matches) < 2 {
		return
	}

	// Convert parsed string back to an integer
	gpuMem, err = strconv.Atoi(matches[1])
	if err != nil {
		return
	}

	return
}

func GetMemPressure() (memPressure float64, err error) {

	memPressure = 0.0

	cmd := exec.Command("memory_pressure", "-Q")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		return
	}

	// Regex to extract the "Device Utilization %" integer value
	re := regexp.MustCompile(`free percentage: (\d+)%`)
	matches := re.FindStringSubmatch(out.String())
	if len(matches) < 2 {
		err = fmt.Errorf("could not parse memory pressure; confirm this is an Apple Silicon Mac")
		return
	}

	// Convert parsed string back to an float
	memPressure, err = strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return
	}
	memPressure = 100 - memPressure

	return

}

func GetCpuUsage() (cpuUsage float64, err error) {

	cpuUsage = 0.0
	cmd := exec.Command("top", "-l", "1", "-n", "0", "-s", "1")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err = cmd.Run(); err != nil {
		return
	}
	// logrus.Infof("Received output %s", out.String())
	// Regex to extract the "Device Utilization %" integer value
	re := regexp.MustCompile(`CPU usage:.*, (\d+.\d+)% idle`)
	matches := re.FindStringSubmatch(out.String())

	if len(matches) < 2 {
		err = fmt.Errorf("could not parse CPU utilization; confirm this is an Apple Silicon Mac")
		return
	}

	// Convert parsed string back to an float
	cpuUsage, err = strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return
	}
	cpuUsage = 100 - cpuUsage

	return
}
