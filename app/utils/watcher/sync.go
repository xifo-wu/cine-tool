package watcher

import (
	"cine-tool/app/model"
	"cine-tool/app/utils/symlink"
	"cine-tool/core"
	"log"
	"os"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

type SyncWatcher struct {
	Watcher *fsnotify.Watcher
}

func (w *SyncWatcher) WatchSyncDir() {
	var dirs []model.CloudSymlinkSync

	core.DB.Where("cloud_path <> ?", "").Find(&dirs)

	for _, item := range dirs {
		w.AddSyncDirs(item.CloudPath)
	}
}

func (w *SyncWatcher) AddSyncDirs(syncPath string) {
	err := filepath.Walk(syncPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = w.Watcher.Add(path)
			if err != nil {
				log.Println("无法添加路径到监视器:", err)
			} else {
				log.Println("添加路径到监视器:", path)
			}
		}
		return nil
	})

	if err != nil {
		log.Println("无法遍历目录:", err)
	}
}

func (w *SyncWatcher) RemoveSyncDirs(syncPath string) {
	err := filepath.Walk(syncPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			err = w.Watcher.Remove(path)
			if err != nil {
				log.Println("无法移除监控目录:", err)
			} else {
				log.Println("移除监控目录:", path)
			}
		}
		return nil
	})

	if err != nil {
		log.Println("无法遍历目录:", err)
	}
}

func (w *SyncWatcher) HandleEvent(event fsnotify.Event) {
	// 只监听 CREATE 和 REMOVE 事件，不监听其他事件
	log.Println("监听到事件:", event)
	if !(event.Op == fsnotify.Create || event.Op == fsnotify.Remove) {
		// 条件取反后的逻辑
		log.Println("Event is not CREATE or REMOVE")
		return
	}

	fi, err := os.Stat(event.Name)

	// 如果创建了目录也添加到监控中。
	if event.Op == fsnotify.Create {
		if err == nil && fi.IsDir() {
			w.AddSyncDirs(event.Name)
		}
	}

	if event.Op != fsnotify.Remove && !fi.IsDir() && !symlink.IsExtensionMatch(event.Name) {
		return
	}

	var dirs []model.CloudSymlinkSync

	core.DB.Where("cloud_path <> ?", "").Find(&dirs)
	var cloudSymlinkSync model.CloudSymlinkSync
	for _, item := range dirs {
		if symlink.IsSubdirectory(item.CloudPath, event.Name) {
			cloudSymlinkSync = item
			break
		}
	}

	if cloudSymlinkSync.ID == 0 {
		return
	}

	relPath, err := filepath.Rel(cloudSymlinkSync.CloudPath, event.Name)
	if err != nil {
		log.Println("无法获取相对路径:", err)
		return
	}

	targetPath := filepath.Join(cloudSymlinkSync.LocalPath, relPath)

	if event.Op == fsnotify.Remove {
		err := os.Remove(targetPath)
		if err != nil {
			log.Println("删除软链接错误：", err)
		}
		return
	}

	log.Println(targetPath, event.Op, "event.Op")
	if event.Op == fsnotify.Create {
		err = symlink.CreateSymlink(event.Name, targetPath)
		if err != nil {
			log.Println("无法创建软链接:", err)
		}
	}
}
