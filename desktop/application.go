/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package desktop

// #cgo pkg-config: gio-2.0
// #include <glib.h>
// #include <stdlib.h>
// int content_type_can_be_executable(char* type);
// int content_type_is(char* type, char* expected_type);
import "C"
import "unsafe"
import (
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"path/filepath"
	"strings"

	"gir/gio-2.0"
	"gir/glib-2.0"
	"pkg.deepin.io/dde/api/thumbnails"
	"pkg.deepin.io/lib/dbus"
	"pkg.deepin.io/lib/utils"
	. "pkg.deepin.io/service/file-manager-backend/log"
	"pkg.deepin.io/service/file-manager-backend/operations"
	"unicode/utf8"
)

// GetUserSpecialDir returns user special dir, like music directory.
func GetUserSpecialDir(dir glib.UserDirectory) string {
	return glib.GetUserSpecialDir(dir)
}

// GetDesktopDir returns desktop's path.
func GetDesktopDir() string {
	return GetUserSpecialDir(glib.UserDirectoryDirectoryDesktop)
}

// IMenuable is the interface for something can generate menu.
type IMenuable interface {
	GenMenu() (*Menu, error)
	destroy()
}

// used by RequestOpen signal
const (
	// OpOpen indicates opening files
	OpOpen int32 = iota
	// OpSelect indicates selecting open programming.
	OpSelect
)

// Application for desktop daemon.
type Application struct {
	desktop  *Desktop
	settings *Settings
	menuable IMenuable
	menu     *Menu

	ActivateFlagDisplay       int32
	ActivateFlagRunInTerminal int32
	ActivateFlagRun           int32

	RequestOpenPolicyOpen   int32
	RequestOpenPolicySelect int32

	RequestRename                 func(string)
	RequestDelete                 func([]string)
	RequestEmptyTrash             func()
	RequestSort                   func(string)
	RequestCleanup                func()
	ReqeustAutoArrange            func()
	RequestCreateFile             func()
	RequestCreateFileFromTemplate func(string)
	RequestCreateDirectory        func()
	ItemCut                       func([]string)
	ItemCopied                    func([]string)

	RequestOpen func([]string, []int32)

	AppGroupCreated func(string, []string)
	AppGroupDeleted func(string)
	AppGroupMerged  func(string, []string)

	ItemDeleted  func(string)
	ItemCreated  func(string)
	ItemModified func(string)
}

// GetDBusInfo returns dbus info of Application.
func (app *Application) GetDBusInfo() dbus.DBusInfo {
	return dbus.DBusInfo{
		Dest:       "com.deepin.dde.daemon.Desktop",
		ObjectPath: "/com/deepin/dde/daemon/Desktop",
		Interface:  "com.deepin.dde.daemon.Desktop",
	}
}

const (
	// ActivateFlagNone do nothing.
	ActivateFlagNone int32 = iota
	// ActivateFlagRun run file directly.
	ActivateFlagRun
	// ActivateFlagRunInTerminal run files in terminal.
	ActivateFlagRunInTerminal
	// ActivateFlagDisplay display files.
	ActivateFlagDisplay
)

// NewApplication creates a application, settings must not be nil.
func NewApplication(s *Settings) *Application {
	app := &Application{
		ActivateFlagRun:           ActivateFlagRun,
		ActivateFlagRunInTerminal: ActivateFlagRunInTerminal,
		ActivateFlagDisplay:       ActivateFlagDisplay,
		RequestOpenPolicyOpen:     OpOpen,
		RequestOpenPolicySelect:   OpSelect,
	}
	app.desktop = NewDesktop(app)
	app.settings = s
	return app
}

func (app *Application) setSettings(s *Settings) {
	app.settings = s
}

func (app *Application) emitRequestRename(uri string) error {
	return dbus.Emit(app, "RequestRename", uri)
}

func (app *Application) emitRequestDelete(uris []string) error {
	return dbus.Emit(app, "RequestDelete", uris)
}

func (app *Application) emitRequestEmptyTrash() error {
	return dbus.Emit(app, "RequestEmptyTrash")
}

func (app *Application) emitRequestSort(sortPolicy string) error {
	return dbus.Emit(app, "RequestSort", sortPolicy)
}

func (app *Application) emitRequestCleanup() error {
	return dbus.Emit(app, "RequestCleanup")
}

func (app *Application) emitRequestAutoArrange() error {
	return dbus.Emit(app, "RequestAutoArrange")
}

func (app *Application) emitRequestCreateFile() error {
	return dbus.Emit(app, "RequestCreateFile")
}

func (app *Application) emitRequestCreateFileFromTemplate(template string) error {
	return dbus.Emit(app, "RequestCreateFileFromTemplate", template)
}

func (app *Application) emitRequestCreateDirectory() error {
	return dbus.Emit(app, "RequestCreateDirectory")
}

func (app *Application) emitItemCut(uris []string) error {
	return dbus.Emit(app, "ItemCut", uris)
}

func (app *Application) emitItemCopied(uris []string) error {
	return dbus.Emit(app, "ItemCopied", uris)
}

func (app *Application) emitRequestOpen(uris []string, op []int32) {
	dbus.Emit(app, "RequestOpen", uris, op)
}

func (app *Application) emitAppGroupCreated(group string, files []string) error {
	return dbus.Emit(app, "AppGroupCreated", group, files)
}

func (app *Application) emitAppGroupMerged(group string, files []string) error {
	return dbus.Emit(app, "AppGroupMerged", group, files)
}

func (app *Application) emitRequestPaste(dest string) error {
	conn, err := dbus.SessionBus()
	if err != nil {
		return err
	}

	obj := conn.Object("com.deepin.filemanager.Backend.Clipboard", "/com/deepin/filemanager/Backend/Clipboard")
	if obj != nil {
		return obj.Call("com.deepin.filemanager.Backend.Clipboard.EmitPaste", 0, dest).Store()
	}

	return nil
}

func (app *Application) showProperties(uris []string) {
	go exec.Command("deepin-nautilus-properties", uris...).Run()
}

func isAppGroup(uri string) bool {
	return filepath.HasPrefix(filepath.Base(uri), AppGroupPrefix)
}

func isAllAppGroup(uris []string) bool {
	for _, uri := range uris {
		if !isAppGroup(uri) {
			return false
		}
	}
	return true
}

// IsAppGroup returns whether uri is a AppGroup
func (app *Application) IsAppGroup(uri string) bool {
	return isAppGroup(uri)
}

func isTrash(uri string) bool {
	return strings.HasPrefix(uri, "trash:")
}

func isComputer(uri string) bool {
	return strings.HasPrefix(uri, "computer:")
}

func isDesktopFile(uri string) bool {
	return strings.HasSuffix(uri, ".desktop")
}

func (app *Application) getMenuableWithExtraItems(uris []string, withExtraItems bool) IMenuable {
	if len(uris) == 0 {
		return app.desktop.enableExtraItems(withExtraItems)
	}

	if len(uris) == 1 {
		uri := uris[0]
		if isTrash(uri) {
			return NewTrashItem(app, uri).enableExtraItems(withExtraItems)
		} else if isComputer(uri) {
			return NewComputerItem(app, uri).enableExtraItems(withExtraItems)
		}
	}

	// multiple selection.
	// 1. open each file.
	// 2. if files whose open programming are unknown exist, notify front-end to select open programming,
	//    and open others files with default open programming.
	// 3. if files which should ask for behaviour exist, notify front-end to ask for behaviour one by one(desktop files shouldn't be asked),
	//    and open others files with default open programming.
	// 4. if all files are archived files, just 'extract here' for archived files.
	// 5. if non-archived files exist, 'extract here' shouldn't be shown.
	// 6. if specific item exists, just 'open' menu item exists.

	// disable app group right-button menu
	// if isAllAppGroup(uris) {
	// 	return NewAppGroup(app, uris).enableExtraItems(withExtraItems)
	// }

	return NewItem(app, uris).enableExtraItems(withExtraItems)
}

// GenMenuContent returns the menu content in json format used in DeepinMenu.
func (app *Application) GenMenuContent(uris []string) string {
	return app.GenMenuContentWithExtraItems(uris, app.settings.DisplayExtraItems())
}

func (app *Application) GenMenuContentWithExtraItems(uris []string, withExtraItems bool) string {
	if app.menuable != nil {
		app.menuable.destroy()
	}
	app.menuable = app.getMenuableWithExtraItems(uris, withExtraItems)
	if app.menuable == nil {
		Log.Error("get menuable item failed")
		return ""
	}
	menu, err := app.menuable.GenMenu()
	if err != nil {
		Log.Error("gen menu failed:", err)
		return ""
	}

	app.menu = menu
	return menu.ToJSON()
}

// HandleSelectedMenuItemWithTimestamp will handle selected menu item according to passed id and timestamp.
func (app *Application) HandleSelectedMenuItemWithTimestamp(id string, timestamp uint32) {
	if app.menu == nil {
		return
	}

	defer func() { app.menu = nil }()
	app.menu.HandleAction(id, timestamp)
}

// HandleSelectedMenuItem will handle selected menu item according to passed id.
func (app *Application) HandleSelectedMenuItem(id string) {
	app.HandleSelectedMenuItemWithTimestamp(id, 0)
}

// DestroyMenu destroys the useless menu.
func (app *Application) DestroyMenu() {
}

func filterDesktop(files []string) []string {
	availableFiles := []string{}
	for _, file := range files {
		if strings.HasSuffix(file, ".desktop") {
			availableFiles = append(availableFiles, file)
		}
	}
	return availableFiles
}

// RequestCreatingAppGroup creates app group according to the files, and emits AppGroupCreated signal when it's done.
func (app *Application) RequestCreatingAppGroup(files []string) error {
	C.g_reload_user_special_dirs_cache()
	desktopDir := GetDesktopDir()

	availableFiles := filterDesktop(files)

	// get group name
	groupName := getGroupName(availableFiles)

	dirName := AppGroupPrefix + groupName

	// create app group.
	createJob := operations.NewCreateDirectoryJob(desktopDir, dirName, nil)
	createJob.Execute()

	if err := createJob.GetError(); err != nil {
		Log.Error("create appgroup failed:", err)
		return err
	}

	appGroupURI := createJob.Result().(string)

	// move files into app group.
	moveJob := operations.NewMoveJob(availableFiles, appGroupURI, "", 0, nil)
	moveJob.Execute()
	if err := moveJob.GetError(); err != nil {
		Log.Error("move apps to appgroup failed:", err)
		return err
	}

	return app.emitAppGroupCreated(appGroupURI, availableFiles)
}

// RequestMergeIntoAppGroup will merge files into existed AppGroup, and emits AppGroupMerged signal when it's done.
func (app *Application) RequestMergeIntoAppGroup(files []string, appGroup string) error {
	availableFiles := filterDesktop(files)

	moveJob := operations.NewMoveJob(availableFiles, appGroup, "", 0, nil)
	moveJob.Execute()
	if err := moveJob.GetError(); err != nil {
		return err
	}

	return app.emitAppGroupMerged(appGroup, availableFiles)
}

func (app *Application) doDisplayFile(file *gio.File, contentType string, timestamp uint32) error {
	defaultApp := gio.AppInfoGetDefaultForType(contentType, false)
	if defaultApp == nil {
		return errors.New("unknown default application")
	}
	defer defaultApp.Unref()

	_, err := defaultApp.Launch([]*gio.File{file}, gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))

	return err
}

// displayFile will display file using default app.
func (app *Application) displayFile(file string, timestamp uint32) error {
	f := gio.FileNewForCommandlineArg(file)
	if f == nil {
		return errors.New("invalid file")
	}
	defer f.Unref()

	info, err := f.QueryInfo(gio.FileAttributeStandardContentType, gio.FileQueryInfoFlagsNone, nil)
	if err != nil {
		return err
	}
	defer info.Unref()

	contentType := info.GetContentType()
	return app.doDisplayFile(f, contentType, timestamp)
}

// ActivateFile will activate file.
// NB: **deprecated**, compatible interface, use ActivateFileWithTimestamp instead.
func (app *Application) ActivateFile(file string, args []string, isExecutable bool, flag int32) error {
	return app.ActivateFileWithTimestamp(file, args, isExecutable, 0, flag)
}

// ActivateFileWithTimestamp will activate file.
func (app *Application) ActivateFileWithTimestamp(file string, args []string, isExecutable bool, timestamp uint32, flag int32) error {
	if isDesktopFile(file) && isExecutable {
		return app.activateDesktopFile(file, args)
	}

	return app.activateFile(file, args, isExecutable, timestamp, flag)
}

func (app *Application) activateDesktopFile(file string, args []string) error {
	uri, err := url.Parse(file)
	if err != nil {
		return err
	}

	// NewDesktopAppInfoFromFilename cannot use uri.
	a := gio.NewDesktopAppInfoFromFilename(uri.Path)
	if a == nil {
		return errors.New("invalid desktop file")
	}
	defer a.Unref()

	_, err = a.LaunchUris(args, gio.GetGdkAppLaunchContext())
	return err
}

func contentTypeCanBeExecutable(contentType string) bool {
	cContentType := C.CString(contentType)
	defer C.free(unsafe.Pointer(cContentType))

	return C.int(C.content_type_can_be_executable(cContentType)) == 1
}

func contentTypeIs(contentType, t string) bool {
	cContentType := C.CString(contentType)
	defer C.free(unsafe.Pointer(cContentType))

	cT := C.CString(t)
	defer C.free(unsafe.Pointer(cT))

	return C.int(C.content_type_is(cContentType, cT)) == 1
}

func isExecutableScript() bool {
	return false
}

func (app *Application) doActivateFile(f *gio.File, args []string, isExecutable bool, contentType string, timestamp uint32, flag int32) error {
	plainType := "text/plain"
	file := utils.DecodeURI(f.GetUri())

	if isExecutable && (contentTypeCanBeExecutable(contentType) || strings.HasSuffix(file, ".bin")) {
		if contentTypeIs(contentType, plainType) { // runable file
			switch flag {
			case app.ActivateFlagRunInTerminal:
				runInTerminal("", file)
				return nil
			case app.ActivateFlagRun:
				handler := exec.Command(file)
				err := handler.Start()
				go handler.Wait()
				return err
			case app.ActivateFlagDisplay:
				return app.doDisplayFile(f, contentType, timestamp)
			}
		} else { // binary file
			// FIXME: strange logic from dde-workspace, why args is not used on the other places.
			handler := exec.Command(file, args...)
			err := handler.Start()
			go handler.Wait()
			return err
		}
	}

	return app.doDisplayFile(f, contentType, timestamp)
}

func (app *Application) activateFile(file string, args []string, isExecutable bool, timestamp uint32, flag int32) error {
	f := gio.FileNewForCommandlineArg(file)
	if f == nil {
		return errors.New("invalid file")
	}
	defer f.Unref()

	info, err := f.QueryInfo(gio.FileAttributeStandardContentType, gio.FileQueryInfoFlagsNone, nil)
	if err != nil {
		return err
	}
	defer info.Unref()

	contentType := info.GetContentType()
	return app.doActivateFile(f, args, isExecutable, contentType, timestamp, flag)
}

// TODO: move to filemanager.

// ItemInfo includes some simple informations.
type ItemInfo struct {
	DisplayName string
	BaseName    string
	URI         string
	MIME        string
	Icon        string
	IconName    string
	Thumbnail   string
	Size        int64
	FileType    uint16
	IsBackup    bool
	IsHidden    bool
	IsReadOnly  bool
	IsSymlink   bool
	CanDelete   bool
	CanExecute  bool
	CanRead     bool
	CanRename   bool
	CanTrash    bool
	CanWrite    bool
}

func toItemInfo(p operations.ListProperty) ItemInfo {
	return ItemInfo{
		DisplayName: p.DisplayName,
		BaseName:    p.BaseName,
		URI:         p.URI,
		MIME:        p.MIME,
		Size:        p.Size,
		FileType:    p.FileType,
		IsBackup:    p.IsBackup,
		IsHidden:    p.IsHidden,
		IsReadOnly:  p.IsReadOnly,
		IsSymlink:   p.IsSymlink,
		CanDelete:   p.CanDelete,
		CanExecute:  p.CanExecute,
		CanRead:     p.CanRead,
		CanRename:   p.CanRename,
		CanTrash:    p.CanTrash,
		CanWrite:    p.CanWrite,
	}
}

func (app *Application) getItemInfo(p operations.ListProperty) ItemInfo {
	info := toItemInfo(p)

	// TODO: filter MIME type
	if app.settings.thumbnailSizeLimitation >= uint64(info.Size) {
		thumbnail, err := thumbnails.GenThumbnailWithMime(p.URI, p.MIME, app.settings.iconSize)
		info.Thumbnail = thumbnail

		if err != nil {
			Log.Warningf("Get thumbnail for %q failed: %s\n", p.URI, err)
		}
	}

	info.Icon = operations.GetThemeIcon(p.URI, app.settings.iconSize)
	info.IconName = operations.GetIconName(p.URI)

	return info
}

func shouldNotShow(p operations.ListProperty) bool {
	return p.IsBackup || (p.IsHidden && !isAppGroup(p.URI))
}

func (app *Application) listDir(dir string, flag operations.ListJobFlag) (map[string]ItemInfo, error) {
	infos := map[string]ItemInfo{}

	job := operations.NewListDirJob(dir, flag)

	job.ListenProperty(func(p operations.ListProperty) {
		if !shouldNotShow(p) {
			infos[p.URI] = app.getItemInfo(p)
		}
	})

	job.Execute()

	if job.HasError() {
		return map[string]ItemInfo{}, job.GetError()
	}

	return infos, nil
}

func fixDisplayNameForAppGroup(info ItemInfo) ItemInfo {
	if info.FileType == uint16(gio.FileTypeDirectory) && strings.HasPrefix(info.BaseName, AppGroupPrefix) {
		info.DisplayName = strings.TrimPrefix(info.DisplayName, AppGroupPrefix)
	}

	// fix bug: name contains invalid coding
	if !utf8.ValidString(info.DisplayName) {
		info.DisplayName = fmt.Sprintf("%q", info.DisplayName)
	}
	if !utf8.ValidString(info.BaseName) {
		info.BaseName = fmt.Sprintf("%q", info.BaseName)
	}

	return info
}

// GetDesktopItems returns all desktop files.
func (app *Application) GetDesktopItems() (map[string]ItemInfo, error) {
	infos, err := app.listDir(GetDesktopDir(), operations.ListJobFlagIncludeHidden)
	for k := range infos {
		infos[k] = fixDisplayNameForAppGroup(infos[k])
	}
	return infos, err
}

// GetItemInfo gets ItemInfo for file.
func (app *Application) GetItemInfo(file string) (ItemInfo, error) {
	info := ItemInfo{}
	f := gio.FileNewForCommandlineArg(file)
	if f == nil {
		return info, fmt.Errorf("Invalid file: %q", file)
	}
	defer f.Unref()

	listProperty, err := operations.GetListProperty(f, nil)
	if err != nil {
		return info, err
	}

	info = fixDisplayNameForAppGroup(app.getItemInfo(listProperty))
	return info, nil
}

func (app *Application) GetAppGroupItems(appGroup string) (map[string]ItemInfo, error) {
	return app.listDir(appGroup, operations.ListJobFlagNone)
}
