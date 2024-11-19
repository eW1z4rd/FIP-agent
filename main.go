package main

import (
	"fip-agent/core"
	"flag"
	"fmt"
	"log/slog"
	"os"
)

func init() {
	var (
		config  bool
		version bool
	)

	flag.StringVar(&core.ConfigPath, "c", "./config.yaml", "set the config file path")
	flag.BoolVar(&config, "g", false, "generate the config file")
	flag.BoolVar(&version, "v", false, "show version info")
	flag.Parse()

	if config {
		if err := os.WriteFile("./config.yaml", []byte(core.ConfigTemplate), 0644); err != nil {
			slog.Error("Failed to generate the config file.", "err", err)
			os.Exit(1)
		}
		slog.Info("The config file is generated successfully.\nPlease run again...")
		os.Exit(0)
	}
	if version {
		fmt.Println("FIP Agent Version", core.Version)
		os.Exit(0)
	}
}

func main() {
	slog.Info("FIP Agent started...")

	core.ParseConfig()
	core.JoinCgroup()
	core.InitWatcher()
}
