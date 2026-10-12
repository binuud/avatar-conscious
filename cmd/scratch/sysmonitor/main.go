package main

import (
	"bytes"
	"fmt"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

func main() {
	fmt.Println("Gathering macOS System Metrics...")

	// 1. Get CPU Usage
	cpuPercent, err := cpu.Percent(time.Second, false)
	if err != nil {
		log.Fatalf("Failed to get CPU stats: %v", err)
	}

	// 2. Get Memory Usage
	vMem, err := mem.VirtualMemory()
	if err != nil {
		log.Fatalf("Failed to get Memory stats: %v", err)
	}

	// 3. Get GPU Usage (macOS specific)
	// gpuPercent := getMacGPUUsage()

	// Print Results
	fmt.Println("\n--- SYSTEM METRICS ---")
	fmt.Printf("CPU Usage:    %.2f%%\n", cpuPercent[0])
	fmt.Printf("Memory Usage: %.2f%% (Used: %.2f GB / Total: %.2f GB)\n",
		vMem.UsedPercent,
		float64(vMem.Used)/(1024*1024*1024),
		float64(vMem.Total)/(1024*1024*1024))
	// fmt.Printf("GPU Usage:    %s\n", gpuPercent)
}

// getMacGPUUsage executes powermetrics to read Apple Silicon / Intel GPU utilization.
// Note: This requires root/sudo privileges due to macOS security policies.
func getMacGPUUsage() string {
	// We take a brief 1-sample snapshot using powermetrics
	cmd := exec.Command("sudo", "powermetrics", "-i", "200", "-n", "1", "-s", "gpu_power")
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		// Fallback for non-sudo or older architectures: try system_profiler info
		return "Unknown (Run with 'sudo' to read real-time Apple Silicon GPU metrics)"
	}

	output := out.String()

	// Parse Apple Silicon "GPU Active Resency" (e.g., "GPU Active: 14%")
	reAppleSilicon := regexp.MustCompile(`GPU Active:\s+(\d+)%`)
	match := reAppleSilicon.FindStringSubmatch(output)
	if len(match) > 1 {
		return match[1] + "%"
	}

	// Alternative parse for Intel Macs/AMD configs inside powermetrics
	if strings.Contains(output, "GPU") {
		reIntel := regexp.MustCompile(`(?i)GPU\s+[\w\s]+:\s+(\d+)%`)
		matchIntel := reIntel.FindStringSubmatch(output)
		if len(matchIntel) > 1 {
			return matchIntel[1] + "%"
		}
	}

	return "0% (Idle)"
}
