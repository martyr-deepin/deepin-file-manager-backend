package monitor

import (
	"fmt"
	"github.com/howeyc/fsnotify"
	"pkg.deepin.io/lib/dbus"
	. "pkg.deepin.io/service/file-manager-backend/log"
)

type WatcherID uint32

type Watcher struct {
	watcher  *fsnotify.Watcher
	dbusInfo dbus.DBusInfo
	isClosed bool
	end      chan struct{}

	ID WatcherID

	Changed func(string, uint32)
}

func (w *Watcher) GetDBusInfo() dbus.DBusInfo {
	return w.dbusInfo
}

const (
	FsNotifyCreated uint32 = iota
	FsNotifyDeleted
	FsNotifyModified
	FsNotifyRename
	FsNotifyAttributeChanged
)

func NewWatcher(id uint32, fileURI string) (*Watcher, error) {
	fsWatcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}

	if err := fsWatcher.Watch(fileURI); err != nil {
		fsWatcher.Close()
		return nil, err
	}

	watcher := &Watcher{
		ID:      WatcherID(id),
		watcher: fsWatcher,
		end:     make(chan struct{}),
	}
	watcher.dbusInfo = dbus.DBusInfo{
		Dest:       "com.deepin.filemanager.Backend.Watcher",
		ObjectPath: fmt.Sprintf("/com/deepin/filemanager/Backend/Watcher/%d", watcher.ID),
		Interface:  "com.deepin.filemanager.Backend.Watcher",
	}

	go func() {
		for {
			select {
			case ev := <-fsWatcher.Event:
				var event uint32
				switch {
				case ev.IsAttrib():
					event = FsNotifyAttributeChanged
				case ev.IsCreate():
					event = FsNotifyCreated
				case ev.IsDelete():
					event = FsNotifyDeleted
				case ev.IsModify():
					event = FsNotifyModified
				case ev.IsRename():
					event = FsNotifyRename
				}
				dbus.Emit(watcher, "Changed", ev.Name, event)
			case err := <-fsWatcher.Error:
				Log.Warning("fsWatcher error:", err)
				return
			case <-watcher.end:
				return
			}
		}
	}()

	return watcher, nil
}

func (w *Watcher) finalize() {
	if w.watcher != nil {
		w.watcher.Close()
		w.watcher = nil
	}

	if !w.isClosed {
		close(w.end)
		w.isClosed = true
	}
}
