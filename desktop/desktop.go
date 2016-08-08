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
	"os/exec"
	"path/filepath"
	"sort"

	"gir/glib-2.0"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/utils"
	"pkg.deepin.io/service/file-manager-backend/operations"
)

func getBaseName(uri string) string {
	return filepath.Base(utils.DecodeURI(uri))
}

type byName []string

func (s byName) Less(i, j int) bool {
	return getBaseName(s[i]) < getBaseName(s[j])
}

func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byName) Len() int {
	return len(s)
}

// Desktop is desktop itself.
type Desktop struct {
	app  *Application
	menu *Menu

	displayExtraItems bool
}

// NewDesktop creates new desktop.
func NewDesktop(app *Application) *Desktop {
	return &Desktop{
		app: app,
	}
}

func (desktop *Desktop) destroy() {
}

// GenMenu generates json format menu content used in DeepinMenu for Desktop itself.
func (desktop *Desktop) GenMenu() (*Menu, error) {
	desktop.menu = NewMenu()
	menu := desktop.menu

	menu.AppendItem(NewMenuItem(Tr("New _folder"), func(uint32) {
		desktop.app.emitRequestCreateDirectory()
	}, true))

	// NB: remove new document item for now, revert it when deepin file manager is out.
	newSubMenu := NewMenu().SetIDGenerator(menu.genID)
	// newSubMenu.AppendItem(NewMenuItem(Tr("_Text document"), func(uint32) {
	// 	desktop.app.emitRequestCreateFile()
	// }, true))

	templatePath := GetUserSpecialDir(glib.UserDirectoryDirectoryTemplates)
	job := operations.NewGetTemplateJob(templatePath)
	templates := job.Execute()
	hasTemplates := len(templates) != 0
	if hasTemplates {
		// newSubMenu.AddSeparator()
		sort.Sort(byName(templates))
		for _, template := range templates {
			templateURI := template
			newSubMenu.AppendItem(NewMenuItem(getBaseName(templateURI), func(uint32) {
				desktop.app.emitRequestCreateFileFromTemplate(templateURI)
			}, true))
		}
	}

	newMenuItem := NewMenuItem(Tr("New _document"), func(uint32) {}, hasTemplates)
	newMenuItem.subMenu = newSubMenu
	menu.AppendItem(newMenuItem)

	sortSubMenu := NewMenu().SetIDGenerator(desktop.menu.genID)
	sortPolicies := desktop.app.settings.getSortPolicies()
	for _, sortPolicy := range sortPolicies {
		if _, ok := sortPoliciesName[sortPolicy]; !ok {
			continue
		}

		sortSubMenu.AppendItem(NewMenuItem(Tr(sortPoliciesName[sortPolicy]), func(sortPolicy string) func(uint32) {
			return func(uint32) {
				desktop.app.emitRequestSort(sortPolicy)
			}
		}(sortPolicy), true))
	}

	// TODO: not handle clean up for now.
	// sortSubMenu.AddSeparator().AppendItem(NewMenuItem(Tr("Clean up"), func(uint32) {
	// 	desktop.app.emitRequestCleanup()
	// }, true))

	sortMenuItem := NewMenuItem(Tr("_Sort by"), func(uint32) {}, true)
	sortMenuItem.subMenu = sortSubMenu

	menu.AppendItem(sortMenuItem).AppendItem(NewMenuItem(Tr("_Paste"), func(uint32) {
		desktop.app.emitRequestPaste(GetDesktopDir())
	}, operations.CanPaste(GetDesktopDir())))

	// TODO: plugin
	if true {
		menu.AddSeparator().AppendItem(NewMenuItem(Tr("Display settings(_M)"), func(uint32) {
			showModule("display")
		}, true)).AppendItem(NewMenuItem(Tr("_Corner navigation"), func(uint32) {
			go exec.Command("/usr/lib/deepin-daemon/dde-zone").Run()
		}, true)).AppendItem(NewMenuItem(Tr("Set _wallpaper"), func(uint32) {
			go exec.Command("/usr/lib/deepin-daemon/dde-wallpaper-chooser").Run()
		}, true))
	}

	if desktop.displayExtraItems {
		menu.AppendItem(NewMenuItem(Tr("Open in _terminal"), func(uint32) {
			runInTerminal(GetDesktopDir(), "")
		}, true))
	}

	return menu, nil

}

func (item *Desktop) enableExtraItems(enable bool) *Desktop {
	item.displayExtraItems = enable
	return item
}
