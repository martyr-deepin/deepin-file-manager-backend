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

	"gir/gio-2.0"
	. "pkg.deepin.io/lib/gettext"
)

// TrashItem is TrashItem.
type TrashItem struct {
	*Item
}

// NewTrashItem creates new trash item.
func NewTrashItem(app *Application, uri string) *TrashItem {
	return &TrashItem{NewItem(app, []string{uri})}
}

// GenMenu generates json format menu content used in DeepinMenu for TrashItem.
func (item *TrashItem) GenMenu() (*Menu, error) {
	item.menu = NewMenu()
	clearMenuItemText := Tr("_Clear")

	trash := gio.FileNewForUri("trash://")
	info, err := trash.QueryInfo(gio.FileAttributeTrashItemCount, gio.FileQueryInfoFlagsNone, nil)
	if err != nil {
		return nil, err
	}
	defer info.Unref()

	trashedItemCount := info.GetAttributeUint32(gio.FileAttributeTrashItemCount)
	if item.app.settings.ShowTrashedItemCountIsEnable() {
		clearMenuItemText = fmt.Sprintf(NTr("_Clear %d Item", "_Clear %d Items", int(trashedItemCount)), trashedItemCount)
	}

	return item.menu.AppendItem(NewMenuItem(Tr("_Open"), func(timestamp uint32) {
		item.app.displayFile("trash://", timestamp)
	}, true)).AddSeparator().AppendItem(NewMenuItem(clearMenuItemText, func(uint32) {
		item.emitRequestEmptyTrash()
	}, trashedItemCount != 0)), nil
}

func (item *TrashItem) enableExtraItems(enable bool) *TrashItem {
	item.Item.enableExtraItems(enable)
	return item
}
