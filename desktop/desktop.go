package desktop

import (
	"os/exec"
	"path/filepath"
	. "pkg.deepin.io/lib/gettext"
	"pkg.deepin.io/lib/glib-2.0"
	"pkg.deepin.io/lib/utils"
	"pkg.deepin.io/service/file-manager-backend/operations"
	"sort"
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

	menu.AppendItem(NewMenuItem(Tr("New _folder"), func() {
		desktop.app.emitRequestCreateDirectory()
	}, true))

	newSubMenu := NewMenu().SetIDGenerator(menu.genID)
	newSubMenu.AppendItem(NewMenuItem(Tr("_Text document"), func() {
		desktop.app.emitRequestCreateFile()
	}, true))

	templatePath := GetUserSpecialDir(glib.UserDirectoryDirectoryTemplates)
	job := operations.NewGetTemplateJob(templatePath)
	templates := job.Execute()
	if len(templates) != 0 {
		newSubMenu.AddSeparator()
		sort.Sort(byName(templates))
		for _, template := range templates {
			templateURI := template
			newSubMenu.AppendItem(NewMenuItem(getBaseName(templateURI), func() {
				desktop.app.emitRequestCreateFileFromTemplate(templateURI)
			}, true))
		}
	}

	newMenuItem := NewMenuItem(Tr("New _document"), func() {}, true)
	newMenuItem.subMenu = newSubMenu
	menu.AppendItem(newMenuItem)

	sortSubMenu := NewMenu().SetIDGenerator(desktop.menu.genID)
	sortPolicies := desktop.app.settings.getSortPolicies()
	for _, sortPolicy := range sortPolicies {
		if _, ok := sortPoliciesName[sortPolicy]; !ok {
			continue
		}

		sortSubMenu.AppendItem(NewMenuItem(Tr(sortPoliciesName[sortPolicy]), func(sortPolicy string) func() {
			return func() {
				desktop.app.emitRequestSort(sortPolicy)
			}
		}(sortPolicy), true))
	}

	// TODO: not handle clean up for now.
	// sortSubMenu.AddSeparator().AppendItem(NewMenuItem(Tr("Clean up"), func() {
	// 	desktop.app.emitRequestCleanup()
	// }, true))

	sortMenuItem := NewMenuItem(Tr("_Sort by"), func() {}, true)
	sortMenuItem.subMenu = sortSubMenu

	menu.AppendItem(sortMenuItem).AppendItem(NewMenuItem(Tr("_Paste"), func() {
		desktop.app.emitRequestPaste(GetDesktopDir())
	}, operations.CanPaste(GetDesktopDir())))

	// TODO: plugin
	if true {
		menu.AddSeparator().AppendItem(NewMenuItem(Tr("Display settings(_M)"), func() {
			showModule("display")
		}, true)).AppendItem(NewMenuItem(Tr("_Corner navigation"), func() {
			exec.Command("/usr/lib/deepin-daemon/dde-zone").Start()
		}, true)).AppendItem(NewMenuItem(Tr("_Personalize"), func() {
			showModule("personalization")
		}, true))
	}

	return menu, nil

}
