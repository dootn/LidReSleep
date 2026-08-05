//go:build windows

// Package power runs the hidden window and its message loop on a dedicated
// thread, receiving WM_POWERBROADCAST broadcasts and power setting
// notifications and feeding them to the engine.
package power

import (
	"encoding/binary"
	"runtime"
	"sync/atomic"
	"syscall"
	"unsafe"

	"lidresleep/internal/engine"
	"lidresleep/internal/i18n"
	"lidresleep/internal/log"
	"lidresleep/internal/win32"
)

var (
	hwndMain   uintptr
	lidNotify  uintptr     // HPOWERNOTIFY for GUID_LIDSWITCH_STATE_CHANGE
	lastClosed atomic.Bool // last lid state observed from power notifications
	suspended  atomic.Bool // current system suspend state
)

// queryState returns the current system state: whether it is suspended and
// whether the lid is closed. Injected into the engine.
func queryState() (suspendedNow, lidClosed bool) {
	return suspended.Load(), lastClosed.Load()
}

// Run starts the hidden-window message loop on the current thread. It returns
// only when the loop exits (process teardown).
func Run() {
	runtime.LockOSThread()
	hwnd := setupWindow()
	if hwnd == 0 {
		return
	}
	engine.SetStateQuery(queryState)
	engine.Trigger() // evaluate the current state once the query is available
	runMessageLoop(hwnd)
	win32.UnregisterLidNotify(lidNotify)
}

// setupWindow creates the hidden top-level window and registers lid-state
// notifications.
func setupWindow() uintptr {
	clsName, _ := syscall.UTF16PtrFromString(win32.WndClassName)

	wcex := win32.WndClassEx{
		LpfnWndProc:   syscall.NewCallback(wndProc),
		LpszClassName: uintptr(unsafe.Pointer(clsName)),
	}
	wcex.CbSize = uint32(unsafe.Sizeof(wcex))

	if atom, err := win32.RegisterClassExW(&wcex); atom == 0 {
		log.Fatal("RegisterClassExW failed: %v", err)
	}

	wndName, _ := syscall.UTF16PtrFromString("LidReSleep")
	// Hidden top-level window (never shown); still receives system broadcasts
	// and power notifications.
	hwnd, err := win32.CreateWindowExW(clsName, wndName)
	if hwnd == 0 {
		log.Fatal("CreateWindowExW failed: %v", err)
	}

	notify, err := win32.RegisterLidNotify(hwnd)
	if notify == 0 {
		log.Fatal("RegisterPowerSettingNotification(GUID_LIDSWITCH_STATE_CHANGE) failed: %v", err)
	}

	hwndMain = hwnd
	lidNotify = notify
	return hwnd
}

// refreshLidQuery re-registers the power notification. Registering immediately
// fires one PBT_POWERSETTINGCHANGE carrying the current lid state, so the real
// lid state is re-queried after waking (events missed while the process was
// frozen are not lost).
func refreshLidQuery() {
	if hwndMain == 0 {
		return
	}
	old := lidNotify
	notify, err := win32.RegisterLidNotify(hwndMain)
	if notify != 0 {
		if old != 0 {
			win32.UnregisterLidNotify(old)
		}
		lidNotify = notify
		return
	}
	log.Error(i18n.F("log.resumeNotifyFail", err))
}

// runMessageLoop blocks processing window messages until WM_QUIT.
func runMessageLoop(hwnd uintptr) {
	var m win32.MSG
	for {
		ret, err := win32.GetMessageW(&m)
		if ret == 0 { // WM_QUIT received
			break
		}
		if ret == uintptr(0xFFFFFFFF) { // -1: error
			log.Error(i18n.F("log.getmessageFail", err))
			break
		}
		win32.DispatchMessageW(&m)
	}
	_ = hwnd
}

// wndProc handles power broadcast messages.
func wndProc(hwnd uintptr, uMsg uint32, wParam, lParam uintptr) uintptr {
	if uMsg == win32.WMPowerBroadcast {
		return onPowerBroadcast(wParam, lParam)
	}
	return win32.DefWindowProcW(hwnd, uMsg, wParam, lParam)
}

func onPowerBroadcast(wParam, lParam uintptr) uintptr {
	switch uint32(wParam) {
	case win32.PBTAPMSuspend:
		suspended.Store(true)
		log.Event(i18n.T("log.evtSuspend"))
		engine.Trigger()
	case win32.PBTAPMResumeSuspend, win32.PBTAPMResumeAutomatic:
		suspended.Store(false)
		log.Event(i18n.T("log.evtResume"))
		refreshLidQuery() // force a fresh lid-state event after waking
		engine.Trigger()
	case win32.PBTPowerSettingChange:
		onPowerSettingChange(lParam)
	}
	return 1
}

func onPowerSettingChange(lParam uintptr) {
	if lParam == 0 {
		return
	}
	ps := *(*win32.PowerBroadcastSetting)(*(*unsafe.Pointer)(unsafe.Pointer(&lParam)))
	if ps.PowerSetting != win32.LidGuid {
		return
	}
	// Docs (WinNT.h GUID_LIDSWITCH_STATE_CHANGE): Data is a DWORD,
	// 0x0 = lid closed, 0x1 = lid opened.
	val := uint32(0)
	if ps.DataLength >= 4 {
		val = binary.LittleEndian.Uint32(ps.Data[:])
	} else if ps.DataLength >= 1 {
		val = uint32(ps.Data[0])
	}
	lastClosed.Store(val == 0)
	engine.Trigger()
}
