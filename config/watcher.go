package config

import (
	"xbot/logger"

	"github.com/fsnotify/fsnotify"
)

// Watcher 文件监控器
type Watcher struct {
	path     string
	callback func()
	watcher  *fsnotify.Watcher
}

// NewWatcher 创建文件监控器
func NewWatcher(path string, callback func()) (*Watcher, error) {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	return &Watcher{
		path:     path,
		callback: callback,
		watcher:  watcher,
	}, nil
}

// Start 启动监控
func (w *Watcher) Start() error {
	if err := w.watcher.Add(w.path); err != nil {
		return err
	}

	go w.watch()
	return nil
}

// watch 监控文件变化
func (w *Watcher) watch() {
	for {
		select {
		case event, ok := <-w.watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				logger.Info("配置文件已更新", "path", w.path)
				if w.callback != nil {
					w.callback()
				}
			}
		case err, ok := <-w.watcher.Errors:
			if !ok {
				return
			}
			logger.Error("文件监控错误", "error", err)
		}
	}
}

// Stop 停止监控
func (w *Watcher) Stop() {
	w.watcher.Close()
}
