package operations

import (
	"net/url"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"sync"
)

// ListProperty represents the properties of files when ListJob is executed.
type ListProperty struct {
	DisplayName string
	BaseName    string
	URI         string
	MIME        string
	Icon        string
	Size        int64
	FileType    uint16
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

const (
	// TODO: is progress is needed.
	_ListJobSignalProperty = "property"
)

// ListJob lists a directory.
type ListJob struct {
	*CommonJob

	dir           *gio.File
	recusive      bool
	includeHidden bool
	waiter        sync.WaitGroup // make sure the channel data is all handled.
}

// ListenProperty adds observers to property signal.
func (job *ListJob) ListenProperty(fn func(ListProperty)) (func(), error) {
	return job.SignalManager.ListenSignal(_ListJobSignalProperty, fn)
}

func (job *ListJob) emitProperty(property ListProperty) {
	job.Emit(_ListJobSignalProperty, property)
	// func(f interface{}, args ...interface{}) {
	// 	fn := f.(func(ListProperty))
	// 	fn(property)
	// })
}

func (job *ListJob) init() {
	job.progressUnit = AmountUnitSumOfFilesAndDirs
}

func (job *ListJob) finalize() {
	// don't job.dir.Unref() here, because gobject has no Ref method,
	// if the job.dir is used somewhere else, like list directory
	// recusively, then program will be broken.
	job.waiter.Wait()
	job.CommonJob.finalize()
}

func (job *ListJob) appendChild(child *gio.File) {
	info, err := child.QueryInfo(
		gio.FileAttributeStandardName+
			","+gio.FileAttributeStandardType+
			","+gio.FileAttributeStandardDisplayName+
			","+gio.FileAttributeStandardSize+
			","+gio.FileAttributeStandardAllocatedSize+
			","+gio.FileAttributeTimeModified+
			","+gio.FileAttributeTimeAccess+
			","+gio.FileAttributeUnixMode+
			","+gio.FileAttributeUnixUid+
			","+gio.FileAttributeStandardIsHidden+
			","+gio.FileAttributeStandardIsSymlink+
			","+gio.FileAttributeAccessCanExecute+
			","+gio.FileAttributeAccessCanRead+
			","+gio.FileAttributeAccessCanTrash+
			","+gio.FileAttributeAccessCanWrite+
			","+gio.FileAttributeAccessCanDelete+
			","+gio.FileAttributeAccessCanRename+
			","+gio.FileAttributeStandardContentType,
		gio.FileQueryInfoFlagsNofollowSymlinks,
		nil,
	)

	if err != nil {
		return
	}

	if info == nil {
		return
	}
	defer info.Unref()

	fsInfo, _ := child.QueryFilesystemInfo(gio.FileAttributeFilesystemReadonly, job.cancellable)
	if fsInfo == nil {
		return
	}
	defer fsInfo.Unref()

	uri := child.GetUri()
	displayName := info.GetDisplayName()
	if displayName == "" {
		displayName = child.GetBasename()
	}
	basename := child.GetBasename()
	contentType := info.GetContentType()
	canExecute := info.GetAttributeBoolean(gio.FileAttributeAccessCanExecute)
	if contentType == _DesktopMIMEType && canExecute {
		desktopApp := gio.NewDesktopAppInfoFromFilename(child.GetPath())
		if desktopApp != nil {
			displayName = desktopApp.GetDisplayName()
			desktopApp.Unref()
		}
	}
	gIcon := info.GetIcon()
	icon := ""
	if gIcon != nil {
		icon = gIcon.ToString()
	}

	size := info.GetSize()
	property := ListProperty{
		DisplayName: displayName,
		BaseName:    basename,
		URI:         uri,
		MIME:        contentType,
		FileType:    uint16(info.GetFileType()),
		Icon:        icon,
		Size:        size,
		IsHidden:    info.GetIsHidden(),
		IsReadOnly:  info.GetAttributeBoolean(gio.FileAttributeFilesystemReadonly),
		CanDelete:   info.GetAttributeBoolean(gio.FileAttributeAccessCanDelete),
		CanExecute:  canExecute,
		CanRead:     info.GetAttributeBoolean(gio.FileAttributeAccessCanRead),
		CanRename:   info.GetAttributeBoolean(gio.FileAttributeAccessCanRename),
		CanTrash:    info.GetAttributeBoolean(gio.FileAttributeAccessCanTrash),
		CanWrite:    info.GetAttributeBoolean(gio.FileAttributeAccessCanWrite),
	}

	job.emitProperty(property)
	fileType := child.QueryFileType(gio.FileQueryInfoFlagsNofollowSymlinks, job.cancellable)
	unit := AmountUnitFiles
	switch gio.FileType(fileType) {
	case gio.FileTypeDirectory:
		unit = AmountUnitDirectories
	}

	job.setProcessedAmount(job.processedAmount[unit]+1, unit)
	job.setProcessedAmount(job.processedAmount[AmountUnitSumOfFilesAndDirs]+1, AmountUnitSumOfFilesAndDirs)
}

// Execute ListJob.
func (job *ListJob) Execute() {
	defer job.finalize()
	defer job.emitDone()
	enumerator, err := job.dir.EnumerateChildren(
		gio.FileAttributeStandardName+","+gio.FileAttributeStandardIsHidden,
		gio.FileQueryInfoFlagsNofollowSymlinks,
		job.cancellable)
	if err != nil {
		job.setError(err)
		return
	}

	// walk through the dir for progress.
	job.scanSources([]*gio.File{job.dir})

	var info *gio.FileInfo
	for !job.isAborted() {
		info, err = enumerator.NextFile(job.cancellable)
		if info == nil || err != nil {
			break
		}

		name := info.GetName()
		child := job.dir.GetChild(name)

		if job.includeHidden || !info.GetIsHidden() {
			job.appendChild(child)
		}
		info.Unref()

		// 1. if hidden files are included, check file type.
		// 2. if hidden files are not included and file is not hidden, check file type.
		if job.recusive && (job.includeHidden || name[0] != '.') &&
			child.QueryFileType(gio.FileQueryInfoFlagsNofollowSymlinks, job.cancellable) == gio.FileTypeDirectory {
			newJob := newListDir(child, true, job.includeHidden)
			subDirProcessedAmount := map[AmountUnit]int64{
				AmountUnitBytes:       0,
				AmountUnitFiles:       0,
				AmountUnitDirectories: 0,
			}
			newJob.ListenProcessedAmount(func(size int64, unit AmountUnit) {
				// only Files and Directories is used directly, the SumOfFilesAndDirs will be setted automatically.
				if unit == AmountUnitSumOfFilesAndDirs {
					return
				}
				newSize := job.processedAmount[unit] + size - subDirProcessedAmount[unit]
				job.setProcessedAmount(newSize, unit)
				subDirProcessedAmount[unit] = size
			})
			newJob.ListenProperty(func(property ListProperty) {
				job.emitProperty(property)
			})
			newJob.ListenDone(func(e error) {
				job.setError(e)
			})
			newJob.Execute()
			if newJob.HasError() {
				child.Unref()
				break
			}
		}

		child.Unref()
	}

	enumerator.Close(job.cancellable)
	enumerator.Unref()
	return
}

func newListDir(dir *gio.File, recusive bool, includeHidden bool) *ListJob {
	job := &ListJob{
		CommonJob:     newCommon(nil),
		dir:           dir,
		recusive:      recusive,
		includeHidden: includeHidden,
	}
	job.init()

	return job
}

// NewListDirJob creates a new list job to list the contents of a directory.
// if recusive, recusively list the contents of a directory.
// if includeHidden, list hidden files and direcories as well.
func NewListDirJob(dir *url.URL, recusive bool, includeHidden bool) *ListJob {
	dest := uriToGFile(dir)
	return newListDir(dest, recusive, includeHidden)
}
