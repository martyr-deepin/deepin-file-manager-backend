/**
 * Copyright (C) 2015 Deepin Technology Co., Ltd.
 *
 * This program is free software; you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation; either version 3 of the License, or
 * (at your option) any later version.
 **/

package main

import (
	"os"
	"os/signal"
	"path/filepath"
	"runtime/pprof"
)

func startCPUProfile(name string) error {
	os.MkdirAll(filepath.Dir(name), 0777)
	f, err := os.Create(name)
	if err != nil {
		return nil
	}

	err = pprof.StartCPUProfile(f)
	if err != nil {
		return err
	}
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			switch sig.String() {
			case "interrupt":
				pprof.StopCPUProfile()
				f.Close()
				os.Exit(0)
			}
		}
	}()
	return nil
}
