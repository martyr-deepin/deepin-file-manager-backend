package main

import (
	"log"
	"pkg.linuxdeepin.com/lib"
	"pkg.linuxdeepin.com/lib/dbus"
)

func main() {
	operationBackend := NewOperationBackend()
	info := operationBackend.GetDBusInfo()
	if !lib.UniqueOnSession(info.Dest) {
		log.Println("already exists a session bus named", info.Dest)
		return
	}

	err := dbus.InstallOnSession(operationBackend)
	if err != nil {
		log.Fatal(err)
	}

	dbus.DealWithUnhandledMessage()
	if err = dbus.Wait(); err != nil {
		log.Fatal(err)
	}
}
