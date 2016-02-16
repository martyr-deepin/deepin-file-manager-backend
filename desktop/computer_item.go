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
	. "pkg.deepin.io/lib/gettext"
)

// ComputerItem is computer item on desktop.
type ComputerItem struct {
	*Item
}

// NewComputerItem creates new computer item.
func NewComputerItem(app *Application, uri string) *ComputerItem {
	return &ComputerItem{NewItem(app, []string{uri})}
}

// GenMenu generates json format menu content used in DeepinMenu for ComputerItem.
func (item *ComputerItem) GenMenu() (*Menu, error) {
	item.menu = NewMenu()
	return item.menu.AppendItem(NewMenuItem(Tr("_Open"), func(timestamp uint32) {
		item.app.displayFile("computer://", timestamp)
	}, true)).AddSeparator().AppendItem(NewMenuItem(Tr("_Properties"), func(uint32) {
		showModule("system_info")
	}, true)), nil
}

func (item *ComputerItem) enableExtraItems(enable bool) *ComputerItem {
	item.Item.enableExtraItems(enable)
	return item
}
