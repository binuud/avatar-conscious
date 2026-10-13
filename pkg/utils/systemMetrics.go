package utils

import (
	"os/exec"
	"strconv"
	"strings"
	"time"

	pb "github.com/binuud/conscious/gen/go/v1/conscious"
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

	// gpu info built from project
	// https://github.com/andersrennermalm/gpuinfo.git
	cmd := exec.Command("./gpuinfo", "-p")
	cmdOut, err := cmd.Output()
	gpuUsage := 0.0
	if err != nil {
		// This will catch if the command doesn't exist or fails to run
		logrus.Fatalf("Error executing command, gpuinfo binary missing, only for macos: %v\n", err)
	} else {
		outputStr := strings.TrimSpace(string(cmdOut))
		gpuUsage, err = strconv.ParseFloat(outputStr, 64)
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
