/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package operations

// #cgo CFLAGS: -std=c99
// #cgo pkg-config: gio-unix-2.0 glib-2.0
// #include "utils.c"
import "C"

import (
	"fmt"
	"net/url"
	"strconv"
	"strings"
	"unicode/utf8"
	"unsafe"

	"gir/gio-2.0"
	"gir/gobject-2.0"
	"pkg.deepin.io/lib/gettext"
)

// Tr is a alias for gettext.Tr, which avoids to use dot import.
var Tr = gettext.Tr
var NTr = gettext.NTr

func dummy(...interface{}) {
}

func uriToGFile(uri *url.URL) *gio.File {
	if uri.Scheme == "" {
		uri.Scheme = "file"

	}
	return gio.FileNewForUri(uri.String())
}

func locationListFromUriList(uris []string) []*gio.File {
	files := make([]*gio.File, len(uris))
	for i, uri := range uris {
		files[i] = gio.FileNewForCommandlineArg(uri)
	}
	return files
}

func locationListFromUrlList(uris []*url.URL) []*gio.File {
	files := []*gio.File{}
	for _, uri := range uris {
		files = append(files, uriToGFile(uri))
	}
	return files
}

func getMaxNameLength(fileDir *gio.File) int {
	return int(C.get_max_name_length((*C.struct__GFile)(fileDir.C)))
}

func fatStrReplace(str string, replacement rune) (string, bool) {
	cstr := C.CString(str)
	ok := C.fat_str_replace(cstr, C.char(replacement)) == 1
	newStr := C.GoString(cstr)
	C.free(unsafe.Pointer(cstr))
	return newStr, ok
}

func makeFileNameValidForDestFs(filename string, fsType string) (string, bool) {
	cname := C.CString(filename)
	defer C.free(unsafe.Pointer(cname))
	cFsType := C.CString(fsType)
	defer C.free(unsafe.Pointer(cFsType))

	ok := C.make_file_name_valid_for_dest_fs(cname, cFsType) == 1
	return C.GoString(cname), ok
}

func queryFsType(file *gio.File, cancellable *gio.Cancellable) string {
	fsinfo, _ := file.QueryFilesystemInfo(gio.FileAttributeFilesystemType, cancellable)
	fsType := ""
	if fsinfo != nil {
		fsType = fsinfo.GetAttributeString(gio.FileAttributeFilesystemType)
		fsinfo.Unref()
	}

	return fsType
}

// FilenameGetExtensionOffset is a C function wrap which return the offset of the extension.
func FilenameGetExtensionOffset(name string) int {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))
	cOffset := C.get_filename_extension_offset(cname)
	offset := int(cOffset)
	return offset
}

// FilenameStripExtension returns the basename without extension name.
func FilenameStripExtension(name string) string {
	filename := name

	offset := FilenameGetExtensionOffset(filename)

	if offset != -1 && filename[offset:] != filename {
		filename = filename[:offset]
	}

	return filename
}

func getUtf8FistValidChar(str []rune) int {
	for i, c := range str {
		if utf8.ValidRune(c) {
			return i
		}
	}

	return -1
}

func shortenUtf8Rune(str []rune, reduceByNumBytes int) []rune {
	if reduceByNumBytes <= 0 {
		return str
	}

	baseLen := len(str)
	baseLen -= reduceByNumBytes

	if baseLen <= 0 {
		return []rune("")
	}

	p := 0
	next := -1
	for baseLen != 0 {
		next = getUtf8FistValidChar(str[p:]) + p
		if next == -1 || next-p > baseLen {
			break
		}

		baseLen -= next + 1 - p
		p = next + 1
	}

	return str[:p]
}

// ShortenUtf8String shortens a utf8 string according to the reduceByNumBytes.
func ShortenUtf8String(str string, reduceByNumBytes int) string {
	if reduceByNumBytes <= 0 {
		return str
	}
	return string(shortenUtf8Rune([]rune(str), reduceByNumBytes))
}

func isReadonlyFileSystem(sourceDir *gio.File) bool {
	if sourceDir != nil {
		defer sourceDir.Unref()
		info, _ := sourceDir.QueryFilesystemInfo(gio.FileAttributeFilesystemReadonly, nil)
		if info != nil {
			defer info.Unref()
			return info.GetAttributeBoolean(gio.FileAttributeFilesystemReadonly)
		}
	}
	return false
}

// HasFsID checks whether the file has the expected filesystem ID.
func HasFsID(file *gio.File, expectedID string) bool {
	info, _ := file.QueryInfo(gio.FileAttributeIdFilesystem, gio.FileQueryInfoFlagsNofollowSymlinks, nil)
	if info != nil {
		defer info.Unref()
		id := info.GetAttributeString(gio.FileAttributeIdFilesystem)
		return id == expectedID
	}
	return false
}

// DirIsParentOf checks whether the root is the parent directory of child.
func DirIsParentOf(root *gio.File, child *gio.File) bool {
	f := child.Dup()
	for f != nil {
		if f.Equal(root) {
			f.Unref()
			return true
		}

		tmp := f
		f = f.GetParent()
		tmp.Unref()
	}

	return false
}

func isDir(file *gio.File) bool {
	info, _ := file.QueryInfo(gio.FileAttributeStandardType, gio.FileQueryInfoFlagsNofollowSymlinks, nil)
	if info != nil {
		defer info.Unref()
		return info.GetFileType() == gio.FileTypeDirectory
	}
	return false
}

func cGetUniqueTargetFile(src *gio.File, destDir *gio.File, sameFs bool, destFsType string, count int) *gio.File {
	cDestType := C.CString(destFsType)
	defer C.free(unsafe.Pointer(cDestType))
	cSameFs := 0
	if sameFs {
		cSameFs = 1
	}
	gfile := C.get_unique_target_file((*C.struct__GFile)(src.ImplementsGFile()), (*C.struct__GFile)(destDir.ImplementsGFile()), C.gboolean(cSameFs), cDestType, C.int(count))
	file := (*gio.File)(gobject.ObjectWrap(unsafe.Pointer(gfile), false))
	return file
}

func getUniqueTargetFile(src *gio.File, destDir *gio.File, sameFs bool, destFsType string, count int) (dest *gio.File) {
	maxLength := getMaxNameLength(destDir)
	info, _ := src.QueryInfo(gio.FileAttributeStandardEditName, 0, nil)
	if info != nil {
		defer info.Unref()

		editname := info.GetEditName()
		if editname != "" {
			newName := getDuplicateName(editname, count, maxLength)
			makeFileNameValidForDestFs(newName, destFsType)
			dest, _ = destDir.GetChildForDisplayName(newName)
			if dest != nil {
				return dest
			}
		}
	}

	basename := src.GetBasename()
	if utf8.ValidString(basename) {
		newName := getDuplicateName(basename, count, maxLength)
		makeFileNameValidForDestFs(newName, destFsType)
		dest, _ = destDir.GetChildForDisplayName(newName)
		if dest != nil {
			return dest
		}
	}

	idx := strings.LastIndex(basename, ".")
	if idx != -1 {
		num, _ := strconv.Atoi(basename[idx+1:])
		count += num
		newName := fmt.Sprintf("%s.%d", basename, count)
		makeFileNameValidForDestFs(newName, destFsType)
		dest = destDir.GetChild(newName)
	}

	return dest
}

func getDuplicateName(name string, countInc int, maxLength int) string {
	namebase, suffix, count := ParsePreviousDuplicateName(name)
	duplicatedName := MakeNextDuplicateName(namebase, suffix, count+countInc, maxLength)
	return duplicatedName
}

func CopyTag() string {
	return Tr(" (Copy)")
}

func CopyTagFmt() string {
	return Tr(" (Copy %d)")
}

func ParsePreviousDuplicateName(name string) (namebase string, suffix string, count int) {
	suffixIdx := FilenameGetExtensionOffset(name)
	if suffixIdx == -1 || suffixIdx+1 == len(name) {
		// no suffix
		suffix = ""
	} else {
		suffix = name[suffixIdx:]
	}

	// case1: xxx (Copy), count = 1
	if idx := strings.Index(name, CopyTag()); idx != -1 {
		namebase = name[:idx]
		count = 1
		return
	}

	// case2: xxx (Copy n), count = n
	idx := strings.Index(name, Tr(" (Copy "))
	if idx != -1 {
		if idx > suffixIdx {
			suffix = ""
		}
		namebase = name[:idx]

		if n, _ := fmt.Sscanf(name[idx:], Tr(" (Copy %d"), &count); n == 1 {
			if count < 1 || count > 1000000 {
				// keep the count within a reasonable range
				count = 0
			}
		}
		return
	}

	if suffix == "" {
		namebase = name[:]
	} else {
		namebase = name[:suffixIdx]
	}
	return
}

func makeDuplicateName(namebase string, suffix string, count int) string {
	switch count {
	case 1:
		return namebase + CopyTag() + suffix
	default:
		return namebase + fmt.Sprintf(CopyTagFmt(), count) + suffix
	}
}

func MakeNextDuplicateName(namebase string, suffix string, count int, maxLength int) (newName string) {
	if count < 1 {
		count = 1
	}

	newName = makeDuplicateName(namebase, suffix, count)

	if maxLength > 0 {
		unshortenLength := len(newName)
		if unshortenLength > maxLength {
			newBase := ShortenUtf8String(newName, unshortenLength-maxLength)
			if newBase != "" {
				newName = makeDuplicateName(newBase, suffix, count)
			}
		}
	}

	return
}
