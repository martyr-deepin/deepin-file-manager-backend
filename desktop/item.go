/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package desktop

import (
	"fmt"
	"os/exec"
	"sort"
	"strings"

	"gir/gio-2.0"
	. "pkg.deepin.io/lib/gettext"
	. "pkg.deepin.io/service/file-manager-backend/log"
	"pkg.deepin.io/service/file-manager-backend/operations"
)

type byDisplayName []*gio.AppInfo

func (s byDisplayName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byDisplayName) Less(i, j int) bool {
	return s[i].GetDisplayName() < s[j].GetDisplayName()
}

func (s byDisplayName) Len() int {
	return len(s)
}

func containsSpecificItem(uris []string) bool {
	for _, uri := range uris {
		// disable app group right-button menu
		// if isTrash(uri) || isComputer(uri) || isAppGroup(uri) {
		if isTrash(uri) || isComputer(uri) {
			return true
		}
	}

	return false
}

func getDefaultOpenApp(uri string) (*gio.AppInfo, error) {
	job := operations.NewGetDefaultLaunchAppJob(uri, false)
	job.Execute()
	if job.HasError() {
		return nil, job.GetError()
	}

	return job.Result().(*gio.AppInfo), nil
}

// ArchiveMimeTypes is a list of MIMEType for archive files.
var ArchiveMimeTypes = []string{
	"application/x-gtar",
	"application/x-zip",
	"application/x-zip-compressed",
	"application/zip",
	"application/x-zip",
	"application/x-tar",
	"application/x-7z-compressed",
	"application/x-rar",
	"application/x-rar-compressed",
	"application/x-jar",
	"application/x-java-archive",
	"application/x-war",
	"application/x-ear",
	"application/x-arj",
	"application/x-gzip",
	"application/gzip",
	"application/x-bzip-compressed-tar",
	"application/x-compressed-tar",
	"application/x-archive",
	"application/x-xz-compressed-tar",
	"application/x-bzip",
	"application/x-cbz",
	"application/x-xz",
	"application/x-lzma-compressed-tar",
	"application/x-ms-dos-executable",
	"application/x-lzma",
	"application/x-cd-image",
	"application/x-deb",
	"application/x-rpm",
	"application/x-stuffit",
	"application/x-tzo",
	"application/x-tarz",
	"application/x-tzo",
	"application/x-msdownload",
	"application/x-lha",
	"application/x-zoo",
}

func isArchived(f *gio.File) bool {
	info, err := f.QueryInfo(gio.FileAttributeStandardContentType, gio.FileQueryInfoFlagsNone, nil)
	if err != nil {
		Log.Error("Query file ContentType failed:", err)
		return false
	}
	defer info.Unref()
	contentType := info.GetContentType()
	for _, MIMEType := range ArchiveMimeTypes {
		if contentType == MIMEType {
			return true
		}
	}

	return false
}

func contains(a *gio.AppInfo, b []*gio.AppInfo) bool {
	for _, app := range b {
		if app.GetId() == a.GetId() {
			return true
		}
	}

	return false
}

func getIntersection(a []*gio.AppInfo, b []*gio.AppInfo) []*gio.AppInfo {
	intersection := []*gio.AppInfo{}
	for _, app := range a {
		if contains(app, b) {
			intersection = append(intersection, app)
		}
	}

	return intersection
}

func getPossibleOpenProgramming(uris []string) []*gio.AppInfo {
	openProgrammings := make([][]*gio.AppInfo, len(uris))
	for i, uri := range uris {
		job := operations.NewGetRecommendedLaunchAppsJob(uri)
		job.Execute()
		if job.HasError() {
			break
		}

		openProgrammings[i] = job.Result().([]*gio.AppInfo)
	}

	intersection := openProgrammings[0]
	for i, l := 1, len(openProgrammings); len(intersection) > 0 && i < l; i++ {
		intersection = getIntersection(intersection, openProgrammings[i])
	}

	possibleOpenProgrammings := make([]*gio.AppInfo, len(intersection))
	for i, app := range intersection {
		possibleOpenProgrammings[i] = app.Dup()
	}

	// destroy all apps.
	for _, apps := range openProgrammings {
		for _, app := range apps {
			app.Unref()
		}
	}

	sort.Sort(byDisplayName(possibleOpenProgrammings))

	return possibleOpenProgrammings
}

// Item is Normal Item, like file/directory/link.
type Item struct {
	uri      string
	uris     []string
	files    []*gio.File
	multiple bool
	app      *Application
	menu     *Menu

	displayExtraItems bool
}

// NewItem creates new item.
func NewItem(app *Application, uris []string) *Item {
	return &Item{
		app:      app,
		uri:      uris[0],
		uris:     uris,
		multiple: len(uris) > 1,
	}
}

func (item *Item) emitRequestDelete() error {
	return item.app.emitRequestDelete(item.uris)
}

func (item *Item) emitRequestRename() error {
	return item.app.emitRequestRename(item.uri)
}

func (item *Item) emitRequestEmptyTrash() error {
	return item.app.emitRequestEmptyTrash()
}

func (item *Item) emitRequestCreateFile() error {
	return item.app.emitRequestCreateFile()
}

func (item *Item) emitRequestCreateFileFromTemplate(template string) error {
	return item.app.emitRequestCreateFileFromTemplate(template)
}

func (item *Item) emitRequestCreateDirectory() error {
	return item.app.emitRequestCreateDirectory()
}

func (item *Item) showProperties() {
	item.app.showProperties(item.uris)
}

func (item *Item) destroy() {
	for _, file := range item.files {
		file.Unref()
	}
}

func (item *Item) addOpenWithMenu(possibleOpenProgrammings []*gio.AppInfo) {
	openWithMenuItem := NewMenuItem(Tr("Open with(_A)"), func(uint32) {}, true)
	item.menu.AppendItem(openWithMenuItem)

	openWithSubMenu := NewMenu()
	openWithSubMenu.SetIDGenerator(item.menu.genID)
	openWithMenuItem.subMenu = openWithSubMenu

	for _, app := range possibleOpenProgrammings {
		openWithSubMenu.AppendItem(NewMenuItem(app.GetDisplayName(), func(id string) func(uint32) {
			return func(timestamp uint32) {
				Log.Debug("open with", id)
				app := gio.NewDesktopAppInfo(id)
				if app == nil {
					Log.Error("get app failed when open files:", id)
					return
				}
				defer app.Unref()

				app.Launch(item.files, gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
			}
		}(app.GetId()), true))
		app.Unref()
	}

	if len(possibleOpenProgrammings) > 0 {
		openWithSubMenu.AddSeparator()
	}

	openWithSubMenu.AppendItem(NewMenuItem(Tr("_Others"), func(uint32) {
		exec.Command("deepin-open-chooser", item.uris...).Start()
	}, true))
}

func (item *Item) shouldHasOpenWith() bool {
	for _, uri := range item.uris {
		if isDesktopFile(uri) {
			return false
		}
	}
	return true
}

// GenMenu generates json format menu content used in DeepinMenu for normal itself.
func (item *Item) GenMenu() (*Menu, error) {
	item.menu = NewMenu()
	item.files = make([]*gio.File, len(item.uris))
	for i, uri := range item.uris {
		item.files[i] = gio.FileNewForCommandlineArg(uri)
		if item.files[i] == nil {
			return nil, fmt.Errorf("No such a file or directory: %s", item.uri)
		}
	}

	menu := item.menu.AppendItem(NewMenuItem(Tr("_Open"), func(timestamp uint32) {
		activationPolicy := item.app.settings.ActivationPolicy()

		askingFiles := []string{}
		ops := []int32{}

		for _, itemURI := range item.uris {
			// FIXME: how to handle these errors.
			f := gio.FileNewForCommandlineArg(itemURI)
			if f == nil {
				continue
			}

			info, err := f.QueryInfo(strings.Join([]string{
				gio.FileAttributeAccessCanExecute,
				gio.FileAttributeStandardContentType,
			}, ","), gio.FileQueryInfoFlagsNone, nil)
			if err != nil {
				f.Unref()
				continue
			}

			isExecutable := info.GetAttributeBoolean(gio.FileAttributeAccessCanExecute)
			contentType := info.GetAttributeString(gio.FileAttributeStandardContentType)
			info.Unref()

			if isExecutable && isDesktopFile(itemURI) {
				item.app.activateDesktopFile(itemURI, []string{})
				continue
			}

			if activationPolicy == ActivationPolicyAsk && isExecutable && (contentTypeCanBeExecutable(contentType) || strings.HasSuffix(itemURI, ".bin")) {
				askingFiles = append(askingFiles, itemURI)
				ops = append(ops, OpOpen)
				f.Unref()
				continue
			}

			defaultApp, _ := getDefaultOpenApp(itemURI)
			if defaultApp == nil {
				askingFiles = append(askingFiles, itemURI)
				ops = append(ops, OpSelect)
				f.Unref()
				continue
			}
			defaultApp.Unref()

			item.app.doActivateFile(f, []string{}, isExecutable, contentType, timestamp, ActivateFlagRun)

			f.Unref()
		}

		if len(askingFiles) > 0 {
			item.app.emitRequestOpen(askingFiles, ops)
		}
	}, true))

	if containsSpecificItem(item.uris) {
		return menu, nil
	}

	// 1. multiple selection: not show "open with" if no possible open programmings.
	// 2. signle selection: show "open with" with "chose".
	// 3. desktop file should not has 'open with'.
	if item.shouldHasOpenWith() {
		possibleOpenProgrammings := getPossibleOpenProgramming(item.uris)
		if len(possibleOpenProgrammings) > 0 || !item.multiple {
			item.addOpenWithMenu(possibleOpenProgrammings)
		}
	}

	menu.AddSeparator()

	// TODO: use plugin, remove useless function.
	// disable app group right-button menu
	// isAppGroupItem := isAppGroup(filepath.Dir(item.uri))
	// if !isAppGroupItem {
	runFileRoller := func(cmd string, files []*gio.File, timestamp uint32) error {
		app, err := gio.AppInfoCreateFromCommandline(cmd, "", gio.AppInfoCreateFlagsSupportsStartupNotification)
		if err != nil {
			return err
		}
		defer app.Unref()
		_, err = app.Launch(files, gio.GetGdkAppLaunchContext().SetTimestamp(timestamp))
		return err
	}

	menu.AppendItem(NewMenuItem(Tr("Co_mpress"), func(timestamp uint32) {
		err := runFileRoller("file-roller -d %U", item.files, timestamp)
		if err != nil {
			Log.Error("run file-roller failed:", err)
		}
	}, true))

	allIsArchived := true
	for _, file := range item.files {
		if !isArchived(file) {
			allIsArchived = false
			break
		}
	}

	if allIsArchived {
		menu.AppendItem(NewMenuItem(Tr("_Extract Here"), func(timestamp uint32) {
			err := runFileRoller("file-roller -h", item.files, timestamp)
			if err != nil {
				Log.Error("run file-roller failed:", err)
			}
		}, true)).AddSeparator()
	}
	// }

	menu.AppendItem(NewMenuItem(Tr("Cu_t"), func(uint32) {
		operations.CutToClipboard(item.uris)
		item.app.emitItemCut(item.uris)
	}, true)).AppendItem(NewMenuItem(Tr("_Copy"), func(uint32) {
		operations.CopyToClipboard(item.uris)
		item.app.emitItemCopied(item.uris)
	}, true))

	// TODO: no design
	if !item.multiple {
		fileType := item.files[0].QueryFileType(gio.FileQueryInfoFlagsNone, nil)
		if fileType == gio.FileTypeDirectory {
			// menu.AppendItem(NewMenuItem(Tr("Paste _Into"), func(uint32) {
			// 	item.app.emitRequestPaste(item.uri)
			// }, operations.CanPaste(item.uri))).AddSeparator()
			if item.displayExtraItems {
				menu.AppendItem(NewMenuItem(Tr("Open in _terminal"), func(uint32) {
					runInTerminal(item.uri, "")
				}, true))
			}
		}
	}

	menu.AddSeparator()

	// if !isAppGroupItem {
	menu.AppendItem(NewMenuItem(Tr("_Rename"), func(uint32) {
		item.emitRequestRename()
	}, !item.multiple))
	// }

	menu.AppendItem(NewMenuItem(Tr("_Delete"), func(uint32) {
		item.emitRequestDelete()
	}, true))

	menu.AddSeparator()

	return item.menu.AppendItem(NewMenuItem(Tr("_Properties"), func(uint32) {
		item.showProperties()
	}, true)), nil
}

func (item *Item) enableExtraItems(enable bool) *Item {
	item.displayExtraItems = enable
	return item
}
