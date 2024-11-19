package core

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"gopkg.in/yaml.v3"
)

var (
	cfg        *FipConfig
	ConfigPath string
	IsLinux    bool
)

func init() {
	//goland:noinspection GoBoolExpressions
	IsLinux = runtime.GOOS == "linux"
}

// ParseConfig 解析config.yaml配置
func ParseConfig() {
	cfg = new(FipConfig)
	f, err := os.ReadFile(ConfigPath)
	if err != nil {
		LogErrorAndExit("The config.yaml file is not exist.", err)
	}
	if err = yaml.Unmarshal(f, cfg); err != nil {
		LogErrorAndExit("Failed to parse config.yaml file.", err)
	}

	slog.Info("The config.yaml file parsing complete.", "cfg", *cfg)
}

// LogErrorAndExit 打印错误并退出
func LogErrorAndExit(v ...any) {
	slog.Error(fmt.Sprintf("%s\n[Trace] %v\nexit...", v...))
	os.Exit(1)
}

// 创建目录（忽略已存在报错）
func checkMakeDir(name string) error {
	if err := os.Mkdir(name, 0755); err != nil && !os.IsExist(err) {
		return err
	}

	return nil
}

// 遍历目录（递归）
func listDir(path string, m map[string]bool) {
	dirs, _ := os.ReadDir(path)
	for _, dir := range dirs {
		if dir.IsDir() {
			dirPath := filepath.Join(path, dir.Name())
			if !trimDir(dirPath) {
				m[dirPath] = true
				listDir(dirPath, m)
			}
		}
	}
}

// 修剪目录
func trimDir(path string) bool {
	for _, dir := range cfg.Watch.Exclude {
		if filepath.Clean(path) == filepath.Clean(dir) {
			return true
		}
	}

	return false
}

// 判断目录是否存在
func isDirExist(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		} else {
			slog.Error("Failed to check the directory.", "err", err)
			return false
		}
	}

	return fileInfo.IsDir()
}

// 判断元素是否在列表中
func isElementInList(element string, list []string, wildcard bool) bool {
	for _, v := range list {
		if wildcard && v == ".*" {
			return true
		} else if v == element {
			return true
		}
	}

	return false
}

// 判断是否处于放行时间
func isReleaseTime() bool {
	now := time.Now().Local()
	today := now.Format(time.DateOnly)
	start, _ := time.ParseInLocation(time.DateTime, today+" "+cfg.Watch.Release["start"], time.Local)
	end, _ := time.ParseInLocation(time.DateTime, today+" "+cfg.Watch.Release["end"], time.Local)

	if now.After(start) && now.Before(end) {
		return true
	}

	return false
}
