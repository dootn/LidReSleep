//go:build windows

// LidReSleep: a Windows background utility (GUI, based on walk).
//
// Solves "kept awake and heating up after the lid is closed" on Modern
// Standby (S0) devices:
//  1. Sleep immediately on lid close, independent of Windows' "lid close
//     action" setting (screen off triggers S0 standby).
//  2. When woken with the lid still closed, sleep again within 500ms~2s.
//  3. Cancel on lid open; never fights the user for the machine.
//
// Double-clicking the exe opens the GUI panel. The power listener runs on
// a dedicated thread.
package main

import (
	"lidresleep/internal/engine"
	"lidresleep/internal/gui"
	"lidresleep/internal/i18n"
	"lidresleep/internal/power"
	"lidresleep/internal/win32"
)

func main() {
	i18n.Init()

	if win32.CheckSingleInstance() {
		gui.Notify(i18n.T("ui.title"), i18n.T("notify.running"))
		return
	}

	go power.Run()

	gui.RunUI()
	engine.Stop()
}
