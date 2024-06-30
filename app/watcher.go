package app

import (
	"cine-tool/app/utils/watcher"

	"github.com/fsnotify/fsnotify"
)

var Watcher *fsnotify.Watcher

func InitWatcher() (*fsnotify.Watcher, error) {
	Watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return Watcher, err
	}

	w := watcher.SyncWatcher{Watcher: Watcher}
	w.WatchSyncDir()

	return Watcher, nil
}
