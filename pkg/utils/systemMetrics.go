package utils

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"time"

	pb "github.com/binuud/avatar-conscious/gen/go/v1/conscious"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
	"github.com/sirupsen/logrus"
)

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

	gpuUsage := 0
	gpuMem := 0
	if __IsGpuMetricsAvailable {
		gpuUsage, gpuMem, err = GetGpuMetrics()
		if err != nil {
			logrus.Fatalf("Cannot parse GetGpuMetrics output: %v\n", err)
		}
	}

	return &pb.SystemMetrics{
		CpuUsage:      cpuPct,
		TotalRamBytes: uint64(vm.Used),
		UsedRamBytes:  uint64(vm.Total),
		RamUsage:      ramUsage,
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
