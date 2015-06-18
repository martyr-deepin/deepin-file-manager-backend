package delegator_test

import (
	. "deepin-file-manager/delegator"
	"pkg.linuxdeepin.com/lib/gio-2.0"
	"testing"
)

func TestQueryFileInfo(t *testing.T) {
	job := new(QueryFileInfoJob)
	job.QueryFileInfo("/tmp", gio.FileAttributeAccessCanDelete, uint32(gio.FileQueryInfoFlagsNofollowSymlinks))
}
