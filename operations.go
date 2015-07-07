package main

// #cgo pkg-config: glib-2.0
// #include <glib.h>
import "C"
import (
	d "deepin-file-manager/delegator"
	"errors"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"pkg.linuxdeepin.com/lib/dbus"
	"pkg.linuxdeepin.com/lib/glib-2.0"
	"pkg.linuxdeepin.com/lib/operations"
	"strings"
	"sync"
)

// TODO: change to deepin's log package and remove this init function.
func init() {
	log.SetFlags(log.Flags() | log.Lshortfile)
}

var getConn = (func() func() *dbus.Conn {
	var conn *dbus.Conn
	var once sync.Once
	return func() *dbus.Conn {
		once.Do(func() {
			var err error
			conn, err = dbus.SessionBus()
			if err != nil {
				log.Println(err)
				os.Exit(1)
			}
		})
		return conn
	}
}())

// OperationBackend is the backend to create dbus operations for front end.
type OperationBackend struct {
}

// NewOperationBackend creates a new backend for operations.
func NewOperationBackend() *OperationBackend {
	op := &OperationBackend{}
	return op
}

func (*OperationBackend) newUIDelegate(dest string, objPath string, iface string) d.IUIDelegate {
	uiDelegate, err := d.NewUIDelegate(getConn(), dest, objPath, iface)
	if err != nil {
		log.Println(dest, objPath, iface, err)
	}
	return uiDelegate
}

// GetDBusInfo returns the dbus info which is needed to export dbus.
func (*OperationBackend) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       d.JobDestination,
		ObjectPath: d.JobObjectPath,
		Interface:  d.JobDestination,
	}
}

// Because empty object path is invalid, so JobObjectPath is used as default value.
// So empty interface means install failed.
func installJob(job dbus.DBusObject) (string, dbus.ObjectPath, string) {
	dest := ""
	objPath := dbus.ObjectPath(d.JobObjectPath)
	iface := ""
	if job == nil {
		log.Println("try to install nil on session bus")
		return dest, objPath, iface
	}

	err := dbus.InstallOnSession(job)
	if err != nil {
		log.Println("install dbus on session bus failed", err)
		return dest, objPath, iface
	}

	dbusInfo := job.GetDBusInfo()
	// log.Println(dbusInfo.ObjectPath, dbusInfo.Interface)
	return dbusInfo.Dest, dbus.ObjectPath(dbusInfo.ObjectPath), dbusInfo.Interface
}

// pathToURL transforms a absolute path to URL.
func pathToURL(path string) (*url.URL, error) {
	srcURL, err := url.Parse(path)
	if err != nil {
		return srcURL, err
	}

	if !filepath.IsAbs(srcURL.Path) {
		return srcURL, errors.New("a absolute path is requested")
	}

	if srcURL.Scheme == "" {
		srcURL.Scheme = "file"
	}

	return srcURL, nil
}

// create a new operation from fn and install it to session bus.
// NB: using closure and anonymous function is ok,
// but variable-length argument list is much more safer and much more portable.
func newOperationJob(paths []string, fn func([]string, ...interface{}) dbus.DBusObject, args ...interface{}) (string, dbus.ObjectPath, string) {
	objPath := dbus.ObjectPath(d.JobObjectPath)
	iface := ""

	srcURLs := make([]string, len(paths))
	for i, path := range paths {
		srcURL, err := pathToURL(path)
		if err != nil {
			log.Println(err)
			return "", objPath, iface // maybe continue is a better choice.
		}
		srcURLs[i] = srcURL.String()
	}

	return installJob(fn(srcURLs, args...))
}

// NewListJob creates a new list job for front end.
func (*OperationBackend) NewListJob(path string, flags int32) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{path}, func(uri []string, args ...interface{}) dbus.DBusObject {
		return d.NewListJob(uri[0], operations.ListJobFlag(flags))
	})
}

// NewStatJob creates a new stat job for front end.
func (*OperationBackend) NewStatJob(path string) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{path}, func(uri []string, args ...interface{}) dbus.DBusObject {
		return d.NewStatJob(uri[0])
	})
}

// NewDeleteJob creates a new delete job for front end.
func (backend *OperationBackend) NewDeleteJob(path []string, shouldConfirm bool, dest string, objPath string, iface string) (string, dbus.ObjectPath, string) {
	return newOperationJob(path, func(uri []string, args ...interface{}) dbus.DBusObject {
		uiDelegate := backend.newUIDelegate(dest, objPath, iface)
		return d.NewDeleteJob(uri, shouldConfirm, uiDelegate)
	})
}

// NewTrashJob creates a new TrashJob for front end.
func (backend *OperationBackend) NewTrashJob(path []string, shouldConfirm bool, dest string, objPath string, iface string) (string, dbus.ObjectPath, string) {
	return newOperationJob(path, func(uri []string, args ...interface{}) dbus.DBusObject {
		uiDelegate := backend.newUIDelegate(dest, objPath, iface)
		return d.NewTrashJob(uri, shouldConfirm, uiDelegate)
	})
}

// NewEmptyTrashJob creates a new EmptyTrashJob for front end.
func (backend *OperationBackend) NewEmptyTrashJob(shouldConfirm bool, dest string, objPath string, iface string) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{}, func(uris []string, args ...interface{}) dbus.DBusObject {
		uiDelegate := backend.newUIDelegate(dest, objPath, iface)
		return d.NewEmptyTrashJob(shouldConfirm, uiDelegate)
	})
}

// NewChmodJob creates a new change mode job for dbus.
func (*OperationBackend) NewChmodJob(path string, permission uint32) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{path}, func(uris []string, args ...interface{}) dbus.DBusObject {
		return d.NewChmodJob(uris[0], permission)
	})
}

// NewChownJob creates a new change owner job for dbus.
func (*OperationBackend) NewChownJob(path string, newOwner string, newGroup string) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{path}, func(uris []string, args ...interface{}) dbus.DBusObject {
		return d.NewChownJob(uris[0], newOwner, newGroup)
	})
}

// NewCreateFileJob creates a new create file job for dbus.
func (backend *OperationBackend) NewCreateFileJob(destDir string, filename string, initContent string, dest string, objPath string, iface string) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{destDir}, func(uris []string, args ...interface{}) dbus.DBusObject {
		uiDelegate := backend.newUIDelegate(dest, objPath, iface)
		return d.NewCreateFileJob(uris[0], filename, []byte(initContent), uiDelegate)
	})
}

// NewCreateDirectoryJob creates a new create directory job for dbus.
func (backend *OperationBackend) NewCreateDirectoryJob(destDir string, dirname string, dest string, objPath string, iface string) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{destDir}, func(uris []string, args ...interface{}) dbus.DBusObject {
		uiDelegate := backend.newUIDelegate(dest, objPath, iface)
		return d.NewCreateDirectoryJob(uris[0], dirname, uiDelegate)
	})
}

// NewCreateFileFromTemplateJob creates a new create file job fro dbus.
func (backend *OperationBackend) NewCreateFileFromTemplateJob(destDir string, templatePath string, dest string, objPath string, iface string) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{destDir}, func(uris []string, args ...interface{}) dbus.DBusObject {
		templateURL, err := pathToURL(templatePath)
		if err != nil {
			log.Println(err)
			return nil
		}

		uiDelegate := backend.newUIDelegate(dest, objPath, iface)
		return d.NewCreateFileFromTemplateJob(uris[0], templateURL.String(), uiDelegate)
	})
}

// NewLinkJob creates a new link job for dbus.
func (backend *OperationBackend) NewLinkJob(src string, destDir string, dest string, objPath string, iface string) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{src}, func(uris []string, args ...interface{}) dbus.DBusObject {
		destDirURL, err := pathToURL(destDir)
		if err != nil {
			log.Println(err)
			return nil
		}

		uiDelegate := backend.newUIDelegate(dest, objPath, iface)
		return d.NewLinkJob(uris[0], destDirURL.String(), uiDelegate)
	})
}

// NewGetLaunchAppJob creates a new get launch app job for dbus.
func (*OperationBackend) NewGetLaunchAppJob(path string, mustSupportURI bool) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{path}, func(uris []string, args ...interface{}) dbus.DBusObject {
		return d.NewGetDefaultLaunchAppJob(uris[0], mustSupportURI)
	})
}

func (*OperationBackend) NewGetRecommendedLaunchAppsJob(uri string) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{uri}, func(uris []string, args ...interface{}) dbus.DBusObject {
		return d.NewGetRecommendedLaunchAppsJob(uris[0])
	})
}

func (*OperationBackend) NewGetAllLaunchAppsJob() (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{}, func(uris []string, args ...interface{}) dbus.DBusObject {
		return d.NewGetAllLaunchAppsJob()
	})
}

// NewSetLaunchAppJob creates a new set default launch app job for dbus.
func (*OperationBackend) NewSetLaunchAppJob(id string, mimeType string) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{}, func(uris []string, args ...interface{}) dbus.DBusObject {
		if !strings.HasSuffix(id, ".desktop") {
			log.Println("wrong desktop id")
			return nil
		}
		return d.NewSetDefaultLaunchAppJob(id, mimeType)
	})
}

// NewCopyJob creates a new copy job for dbus.
func (backend *OperationBackend) NewCopyJob(srcs []string, destDir string, targetName string, flags uint32, dest string, objPath string, iface string) (string, dbus.ObjectPath, string) {
	return newOperationJob(srcs, func(uris []string, args ...interface{}) dbus.DBusObject {
		destDirURL, err := pathToURL(destDir)
		if err != nil {
			log.Println(err)
			return nil
		}

		uiDelegate := backend.newUIDelegate(dest, objPath, iface)
		return d.NewCopyJob(uris, destDirURL.String(), targetName, flags, uiDelegate)
	})
}

// NewMoveJob creates a new move job for dbus.
func (backend *OperationBackend) NewMoveJob(paths []string, destDir string, targetName string, flags uint32, dest string, objPath string, iface string) (string, dbus.ObjectPath, string) {
	return newOperationJob(paths, func(uris []string, args ...interface{}) dbus.DBusObject {
		destDirURL, err := pathToURL(destDir)
		if err != nil {
			log.Println(err)
			return nil
		}

		uiDelegate := backend.newUIDelegate(dest, objPath, iface)
		return d.NewMoveJob(uris, destDirURL.String(), targetName, flags, uiDelegate)
	})
}

func (backend *OperationBackend) NewRenameJob(fileURL string, newName string) (string, dbus.ObjectPath, string) {
	return newOperationJob([]string{fileURL}, func(uris []string, args ...interface{}) dbus.DBusObject {
		fileURL := uris[0]
		return d.NewRenameJob(fileURL, newName)
	})
}

func (backend *OperationBackend) NewGetTemplateJob() (string, dbus.ObjectPath, string) {
	C.g_reload_user_special_dirs_cache()
	templateDirPath := glib.GetUserSpecialDir(glib.UserDirectoryDirectoryTemplates)

	return newOperationJob([]string{templateDirPath}, func(uris []string, args ...interface{}) dbus.DBusObject {
		templateDirURI := uris[0]
		return d.NewGetTemplateJob(templateDirURI)
	})
}
