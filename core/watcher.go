package core

import (
	"log/slog"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/fsnotify/fsnotify"
)

var eventMapStr = map[fsnotify.Op]string{
	fsnotify.Write:  "write",
	fsnotify.Rename: "rename",
	fsnotify.Remove: "remove",
	fsnotify.Create: "create",
	fsnotify.Chmod:  "chmod",
}

func InitWatcher() {
	checkInotify()

	// Create new watcher.
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		LogErrorAndExit("Failed to create new watcher.", err)
	}
	defer watcher.Close()

	// Start listening for events.
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				handleEvent(event)
				addWatcherDynamic(event, watcher)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				slog.Error("Failed to process watching event.", "err", err)
			}
		}
	}()

	// Add paths.
	addWatcher(watcher)
	slog.Info("FIP is ready. Watching started...", "Mode", cfg.Watch.Mode)

	// Block main goroutine forever.
	<-make(chan bool)
}

func handleEvent(event fsnotify.Event) {
	ext := path.Ext(event.Name)

	if !isElementInList(ext, cfg.Watch.Type, true) {
		return
	}

	msg := "Detected file change:"
	if isReleaseTime() {
		releaseMode := cfg.Watch.Release["mode"]
		switch releaseMode {
		case "quiet":
			return
		case "tag":
			msg = "[release]Detected file change:"
		default:
			msg = "[release]Detected file change:"
		}
	}

	if isElementInList(eventMapStr[event.Op], cfg.Watch.Mode, false) {
		slog.Info(msg, "filename", event.Name, "operation", event.Op)
	}
}

func addWatcher(watcher *fsnotify.Watcher) {
	dirsMap := map[string]bool{}
	slog.Info("Retrieving the directory, waiting...")
	for _, dir := range cfg.Watch.Include {
		dirsMap[dir] = true
		listDir(dir, dirsMap)
	}

	for dir := range dirsMap {
		if _, err := os.Stat(dir); err == nil {
			slog.Info("Add watch", "dir", dir)
			err := watcher.Add(dir)
			if err != nil {
				LogErrorAndExit("Failed to add watching path.", err)
			}
		} else {
			delete(dirsMap, dir)
		}
	}

	slog.Info("Total watched dirs: " + strconv.Itoa(len(dirsMap)))
}

func addWatcherDynamic(event fsnotify.Event, watcher *fsnotify.Watcher) {
	fileInfo, err := os.Stat(event.Name)
	if err != nil {
		return
	}

	if fileInfo.IsDir() && event.Op == fsnotify.Create {
		dirsMap := map[string]bool{}
		dirsMap[event.Name] = true
		listDir(event.Name, dirsMap)

		for dir := range dirsMap {
			if _, err := os.Stat(dir); err == nil {
				slog.Info("Add watch", "dir", dir)
				err := watcher.Add(dir)
				if err != nil {
					slog.Error("Failed to add watching path.", "err", err)
				}
			}
		}
	}
}

func checkInotify() {
	if IsLinux {
		content, err := os.ReadFile("/proc/sys/fs/inotify/max_user_watches")
		if err != nil {
			LogErrorAndExit("The max_user_watches file is not exist.", err)
		}

		size, _ := strconv.Atoi(strings.Trim(string(content), "\n"))
		if size < 124983 {
			slog.Error("The value of Inotify's max_user_watches is too small.")
			os.Exit(1)
		}
	}
}
