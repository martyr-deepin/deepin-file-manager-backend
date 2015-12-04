package main

import (
	"pkg.deepin.io/lib/gettext"
)

func InitI18n() {
	gettext.InitI18n()
	gettext.Textdomain("DFMB")
}
