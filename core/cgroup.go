package core

import (
	"log/slog"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var cgroupDir = "/sys/fs/cgroup"
var cgroupOldDir = "/cgroup" // Red Hat 6.9

func cpuLimit(pid int) error {
	cgroupCpuDir := filepath.Join(cgroupDir, "/cpu/fip-agent")

	const periodStr = "100000"
	n1, _ := strconv.Atoi(periodStr)
	n2, _ := strconv.Atoi(strings.Trim(cfg.Cgroup.MaxCpuUsage, "%"))
	quotaStr := strconv.Itoa(n1 * n2 / 100)

	if err := checkMakeDir(cgroupCpuDir); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(cgroupCpuDir, "cpu.cfs_quota_us"), []byte(quotaStr), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(cgroupCpuDir, "cpu.cfs_period_us"), []byte(periodStr), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(cgroupCpuDir, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return err
	}
	return nil
}

func memoryLimit(pid int) error {
	cgroupMemoryDir := filepath.Join(cgroupDir, "/memory/fip-agent")

	if err := checkMakeDir(cgroupMemoryDir); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(cgroupMemoryDir, "memory.limit_in_bytes"), []byte(cfg.Cgroup.MaxMemoryUsage), 0644); err != nil {
		return err
	}
	if err := os.WriteFile(filepath.Join(cgroupMemoryDir, "cgroup.procs"), []byte(strconv.Itoa(pid)), 0644); err != nil {
		return err
	}
	return nil
}

func JoinCgroup() {
	if IsLinux && (isDirExist(cgroupDir) || isDirExist(cgroupOldDir)) {
		if !isDirExist(cgroupDir) && isDirExist(cgroupOldDir) {
			cgroupDir = cgroupOldDir
		}

		pid := os.Getpid()
		slog.Info("Current process", "pid", pid)
		if err := cpuLimit(pid); err != nil {
			slog.Error("Failed to add cgroup cpu limit.", "err", err)
		}
		if err := memoryLimit(pid); err != nil {
			slog.Error("Failed to add cgroup memory limit.", "err", err)
		}
	} else {
		slog.Warn("The system does not support cgroup resource limits.")
	}
}
