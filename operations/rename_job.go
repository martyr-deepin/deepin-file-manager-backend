/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package operations

import (
	"bufio"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"gir/gio-2.0"
	"gir/glib-2.0"
)

const (
	ErrorInvalidFileName = -1 - iota
	ErrorSameFileName
	ErrorStatFileFailed
	ErrorSaveFileFailed
	ErrorReadFileFailed
)
const _FullNameKey = "X-GNOME-FullName"

type RenameError struct {
	Code int
	Msg  string
}

func (e *RenameError) Error() string {
	j, _ := json.Marshal(e)
	return string(j)
}

func newRenameError(code int, msg string) error {
	return &RenameError{Code: code, Msg: msg}
}

const (
	_RenameJobSignalOldName string = "old-name"
	_RenameJobSignalNewFile string = "new-file"
)

// Rename is just working for current directory.
type RenameJob struct {
	*CommonJob

	file    *gio.File
	newName string

	userDir string // like XDG_VIDEOS_DIR, store on $XDG_CONFIG_HOME/user-dirs.dirs.
}

func (job *RenameJob) emitOldName(oldName string) error {
	return job.Emit(_RenameJobSignalOldName, oldName)
}

func (job *RenameJob) emitNewFile(newFileURL string) error {
	return job.Emit(_RenameJobSignalNewFile, newFileURL)
}

func (job *RenameJob) ListenOldName(fn func(string)) {
	job.ListenSignal(_RenameJobSignalOldName, fn)
}

func (job *RenameJob) ListenNewFile(fn func(string)) {
	job.ListenSignal(_RenameJobSignalNewFile, fn)
}

func (job *RenameJob) checkLocale(keyFile *glib.KeyFile) string {
	locale := ""
	names := GetLanguageNames()
	for _, localeName := range names {
		name, err := keyFile.GetLocaleString(glib.KeyFileDesktopGroup, glib.KeyFileDesktopKeyName, localeName)
		if name != "" && err == nil {
			locale = localeName
			break
		}
	}
	return locale
}

func getUserConfigDir() string {
	h := os.Getenv("XDG_CONFIG_HOME")
	if h == "" {
		home := os.Getenv("HOME")
		h = filepath.Join(home, ".config")
	}
	return h
}

func getUserDirsPath() string {
	// userConfigDir := glib.GetUserConfigDir() // cannot reload
	userConfigDir := getUserConfigDir()
	return filepath.Join(userConfigDir, "user-dirs.dirs")
}

func (job *RenameJob) checkUserDirs() {
	userDirs := getUserDirsPath()

	fileContent, err := ioutil.ReadFile(userDirs)
	if err != nil || string(fileContent) == "" {
		return
	}

	fileReader := bytes.NewReader(fileContent)
	scanner := bufio.NewScanner(fileReader)
	for scanner.Scan() {
		lineText := strings.TrimSpace(scanner.Text())
		if lineText == "" || lineText[0] == '#' {
			continue
		}

		values := strings.SplitN(lineText, "=", 2)
		if len(values) != 2 {
			continue
		}
		userDir := strings.TrimSpace(values[0])
		value := strings.TrimSpace(values[1])

		filePath := job.file.GetPath()
		value = os.ExpandEnv(value)
		if value == filePath {
			job.userDir = userDir
			break
		}
	}
}

func (job *RenameJob) changeUserDir() error {
	buffer := bytes.NewBuffer([]byte{})
	writer := bufio.NewWriter(buffer)
	userDirs := getUserDirsPath()
	fileContent, err := ioutil.ReadFile(userDirs)
	if err != nil {
		return newRenameError(ErrorReadFileFailed, err.Error())
	}

	fileReader := bytes.NewReader(fileContent)
	scanner := bufio.NewScanner(fileReader)
	for scanner.Scan() {
		originLineText := scanner.Text()

		lineText := strings.TrimSpace(originLineText)
		if len(lineText) == 0 || lineText[0] == '#' {
			writer.WriteString(originLineText)
			writer.WriteString("\n")
			continue
		}

		values := strings.SplitN(lineText, "=", 2)
		userDir := strings.TrimSpace(values[0])

		if userDir == job.userDir {
			writer.WriteString(userDir)
			writer.WriteString("=\"$HOME")
			writer.WriteRune(filepath.Separator)
			writer.WriteString(job.newName)
			writer.WriteString("\"\n")
			break
		}
	}
	writer.Flush()
	stat, err := os.Stat(userDirs)
	if err != nil {
		return newRenameError(ErrorStatFileFailed, err.Error())
	}

	e := ioutil.WriteFile(userDirs, buffer.Bytes(), stat.Mode())
	if e != nil {
		return newRenameError(ErrorSaveFileFailed, e.Error())
	}
	return nil
}

func (job *RenameJob) setDesktopName() (string, error) {
	var oldDisplayName string

	keyFile := glib.NewKeyFile()
	defer keyFile.Free()
	filePath := job.file.GetPath()
	_, err := keyFile.LoadFromFile(filePath, glib.KeyFileFlagsKeepComments|glib.KeyFileFlagsKeepTranslations)
	if err != nil {
		e := err.(gio.GError)
		return oldDisplayName, newRenameError(int(e.Code), e.Message)
	}

	appInfo := gio.NewDesktopAppInfoFromKeyfile(keyFile)
	if appInfo != nil {
		oldDisplayName = appInfo.GetDisplayName()
		appInfo.Unref()
	}

	locale := job.checkLocale(keyFile)
	if locale != "" {
		keyFile.SetLocaleString(glib.KeyFileDesktopGroup, glib.KeyFileDesktopKeyName, locale, job.newName)
	} else {
		keyFile.SetString(glib.KeyFileDesktopGroup, glib.KeyFileDesktopKeyName, job.newName)
	}

	_, keys, _ := keyFile.GetKeys(glib.KeyFileDesktopGroup)
	for _, key := range keys {
		if key == _FullNameKey {
			if locale != "" {
				keyFile.SetLocaleString(glib.KeyFileDesktopGroup, _FullNameKey, locale, job.newName)
			} else {
				keyFile.SetString(glib.KeyFileDesktopGroup, _FullNameKey, job.newName)
			}
			break
		}
	}

	_, content, err := keyFile.ToData()
	if err != nil {
		e := err.(gio.GError)
		return oldDisplayName, newRenameError(int(e.Code), e.Message)
	}

	stat, err := os.Stat(filePath)
	if err != nil {
		return oldDisplayName, newRenameError(ErrorStatFileFailed, err.Error())
	}
	e := ioutil.WriteFile(filePath, []byte(content), stat.Mode().Perm())
	if e != nil {
		return oldDisplayName, newRenameError(ErrorSaveFileFailed, e.Error())
	}
	return oldDisplayName, nil
}

func (job *RenameJob) init() {
	job.RegisterMonitor(_RenameJobSignalNewFile)
	job.RegisterMonitor(_RenameJobSignalOldName)

	job.checkUserDirs()
}

func (job *RenameJob) finalize() {
	defer job.CommonJob.finalize()
	job.file.Unref()
}

func (job *RenameJob) isValidName(name string) bool {
	return name == "." || name == ".." || strings.ContainsRune(name, filepath.Separator)
}

func (job *RenameJob) Execute() {
	defer finishJob(job)

	if job.isValidName(job.newName) {
		job.setError(newRenameError(ErrorInvalidFileName, job.newName))
		return
	}

	info, err := job.file.QueryInfo(strings.Join(
		[]string{
			gio.FileAttributeStandardContentType,
			gio.FileAttributeStandardDisplayName,
			gio.FileAttributeAccessCanExecute,
		}, ","),
		gio.FileQueryInfoFlagsNofollowSymlinks, nil)
	if err != nil {
		e := err.(gio.GError)
		job.setError(newRenameError(int(e.Code), e.Message))
		return
	}
	defer info.Unref()

	oldDisplayName := info.GetDisplayName()
	if oldDisplayName == "" {
		oldDisplayName = job.file.GetBasename()
	}

	if oldDisplayName == job.newName {
		job.setError(newRenameError(ErrorSameFileName, job.newName))
		return
	}

	mimeType := info.GetContentType()
	if mimeType == _DesktopMIMEType {
		if info.GetAttributeBoolean(gio.FileAttributeAccessCanExecute) {
			oldDisplayName, err = job.setDesktopName()
			if err != nil {
				job.setError(err)
				return
			}
			job.newName = job.newName + ".desktop"
		}

	}

	job.emitOldName(oldDisplayName)
	newFile, err := job.file.SetDisplayName(job.newName, job.cancellable)
	if newFile != nil {
		job.emitNewFile(newFile.GetUri())
		newFile.Unref()
	}
	if err != nil {
		e := err.(gio.GError)
		job.setError(newRenameError(int(e.Code), e.Message))
		return
	}

	if job.userDir != "" {
		err = job.changeUserDir()
		if err != nil {
			job.setError(err)
		}
	}
}

func newRenameJob(file *gio.File, newName string) *RenameJob {
	job := &RenameJob{
		CommonJob: newCommon(nil),
		file:      file,
		newName:   newName,
	}
	job.init()
	return job
}

func NewRenameJob(file string, newName string) *RenameJob {
	gfile := gio.FileNewForCommandlineArg(file)
	return newRenameJob(gfile, newName)
}
