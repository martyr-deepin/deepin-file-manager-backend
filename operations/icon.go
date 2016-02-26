/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package operations

// #cgo pkg-config: gtk+-3.0
// #cgo CFLAGS: -std=c99
// #include <stdlib.h>
// char* icon_name_to_path_with_check_xpm(const char* name, int size);
// char* get_icon_for_app(char* file_path, int size);
// char* get_icon_for_file(char* icons, int size);
import "C"
import "unsafe"
import (
	"gir/gio-2.0"
	"net/url"
	"os"
	"path/filepath"
	. "pkg.deepin.io/service/file-manager-backend/log"
)

func getIconFromGIcon(icon *gio.Icon, size int, fn func(*C.char, C.int) *C.char) string {
	iconStr := icon.ToString()
	if iconStr == "" {
		return ""
	}

	cIconStr := C.CString(iconStr)
	defer C.free(unsafe.Pointer(cIconStr))

	cIcon := fn(cIconStr, C.int(size))
	defer C.free(unsafe.Pointer(cIcon))

	return C.GoString(cIcon)
}

// GIconGetThemeIconForApp returns the icon for application.
func GIconGetThemeIconForApp(icon *gio.Icon, size int) string {
	return getIconFromGIcon(icon, size, func(icon *C.char, size C.int) *C.char {
		return C.get_icon_for_app(icon, size)
	})
}

// GIconGetThemeIconForFile returns the icon for normal files.
func GIconGetThemeIconForFile(icon *gio.Icon, size int) string {
	return getIconFromGIcon(icon, size, func(icon *C.char, size C.int) *C.char {
		return C.get_icon_for_file(icon, size)
	})
}

// GetThemeIconFromIconName returns icon from icon name.
func GetThemeIconFromIconName(iconName string, size int) string {
	if iconName == "" {
		return ""
	}

	cIconName := C.CString(iconName)
	defer C.free(unsafe.Pointer(cIconName))

	cIcon := C.icon_name_to_path_with_check_xpm(cIconName, C.int(size))
	defer C.free(unsafe.Pointer(cIcon))

	return C.GoString(cIcon)
}

const (
	_UserExecutable os.FileMode = 0500
)

func isUserExecutable(perm os.FileMode) bool {
	return perm&_UserExecutable != 0
}

// GetGIconForApp get gio.Icon for application.
// @param filePath app's filepath.
func GetGIconForApp(filePath string) *gio.Icon {
	u, _ := url.Parse(filePath)
	stat, err := os.Stat(u.Path)
	if err != nil {
		Log.Warning("stat", u.Path, "failed:", err)
		return nil
	}

	if isUserExecutable(stat.Mode().Perm()) {
		app := gio.NewDesktopAppInfoFromFilename(u.Path)
		if app != nil {
			defer app.Unref()
			gicon := app.GetIcon() // transfer none, not unref it.
			return gicon
		}
	}

	return nil
}

// GetGIconForFile gets gio.Icon for file.
// @param filePath is file's path
func GetGIconForFile(filePath string) *gio.Icon {
	// gio.FileNewForCommandlineArg never failed, even if the arg is malformed path.
	file := gio.FileNewForCommandlineArg(filePath)
	defer file.Unref()

	info, err := file.QueryInfo(gio.FileAttributeStandardIcon, gio.FileQueryInfoFlagsNone, nil)
	if info == nil {
		Log.Warning("Query file standard icon failed:", err)
		return nil
	}
	defer info.Unref()

	gicon := info.GetIcon() // transfer none, not unref it.
	return gicon
}

func getThemeIconHelper(filePath string, handlerForApp func(*gio.Icon) string, handlerForFile func(*gio.Icon) string) string {
	if filepath.Ext(filePath) == ".desktop" {
		gicon := GetGIconForApp(filePath)
		if gicon != nil {
			return handlerForApp(gicon)
		}
	}

	gicon := GetGIconForFile(filePath)
	if gicon != nil {
		return handlerForFile(gicon)
	}

	return ""
}

func giconToString(gicon *gio.Icon) string {
	return gicon.ToString()
}

// GetIconName gets icon for app or file.
// @param filePath path for app or file.
func GetIconName(filePath string) string {
	return getThemeIconHelper(filePath, giconToString, giconToString)
}

func GetThemeIconForFile(file string, size int) string {
	return getThemeIconHelper(file, func(gicon *gio.Icon) string {
		return GIconGetThemeIconForApp(gicon, size)
	}, func(gicon *gio.Icon) string {
		return GIconGetThemeIconForFile(gicon, size)
	})
}

// GetThemeIcon returns full path for icon.
// @param iconStr can be uri or path of files, or the icon name.
// @param size is the expected size of icon.
func GetThemeIcon(iconStr string, size int) string {
	icon := ""

	// if iconStr is icon name, url.ParseRequestURI returns invalid uri error.
	_, err := url.ParseRequestURI(iconStr)
	if err != nil {
		icon = GetThemeIconFromIconName(iconStr, size)
	}

	if icon == "" {
		icon = GetThemeIconForFile(iconStr, size)
	}

	return icon
}
