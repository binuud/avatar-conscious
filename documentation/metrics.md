# Quick summary
top -l 1 -n 0 | grep PhysMem

# System-wide pressure and free percentage
memory_pressure -Q

# Page-level breakdown (page size is 16384 on Apple Silicon, 4096 on Intel)
vm_stat

# Total installed RAM
sysctl -n hw.memsize

# Top memory consumers
top -o mem -l 1 -n 10 -stats pid,command,mem

# Overall CPU usage (user / sys / idle), sampled once
top -l 2 -n 0 -s 1 | grep "CPU usage" | tail -1

# Per-process, sorted by CPU
top -o cpu -l 1 -n 10 -stats pid,command,cpu

# Sum of all process CPU % (can exceed 100 on multi-core)
ps -A -o %cpu | awk '{s+=$1} END {print s "%"}'

ioreg -r -d 1 -w 0 -c IOAccelerator


ioreg -r -d 1 -w 0 -c IOAccelerator | grep -o '"Device Utilization %"=[0-9]*\|"In use system memory"=[0-9]*'


# Commands checked in macos to get some useful live metrics

system_profiler SPDisplaysDataType

top -l 2 -stats "pid,cpu,power,command"
ioreg -r -d 1 -c IOAccelerator
 sysctl -n machdep.cpu.brand_string
sw_vers -productVersion
 ioreg -rc AGXAccelerator -d1
 sysctl hw.perflevel0.physicalcpu hw.perflevel1.physicalcpu
 sysctl -n hw.perflevel0.name      
 sysctl -n hw.perflevel1.name

 sysctl machdep.cpu.brand_string
 systemstats  --concise-dashboard json 
 system_profiler SPHardwareDataType SPMemoryDataType -json
top -l 2 -n 0 -s 1 | grep "CPU usage" | tail -1
top -o cpu -l 1 -n 10 -stats pid,command,cpu
ps -A -o %cpu | awk '{s+=$1} END {print s "%"}'
top -l 1 -n 0 | grep PhysMemcon