//go:build windows

// Package win32 provides raw Win32 bindings used across the app: power events,
// sleep calls, single-instance and a few system calls.
package win32

import (
	"syscall"
	"unsafe"
)

// Win32 constants.
const (
	WMPowerBroadcast         = 0x0218
	WMSysCommand             = 0x0112
	SCMonitorPower           = 0xF170 // SC_MONITORPOWER
	MonitorLowPower          = 2      // display low power (screen off)
	PBTAPMSuspend            = 0x0004
	PBTAPMResumeSuspend      = 0x0007
	PBTAPMResumeAutomatic    = 0x0012
	PBTPowerSettingChange    = 0x8013
	DeviceNotifyWindowHandle = 0x00000000
	HWNDBroadcast            = 0xFFFF // HWND_BROADCAST (screen-off sent to all top-level windows)
)

// WndClassName is the hidden-window class name.
const WndClassName = "LidReSleepHiddenWnd"

// Win32 structs.

type Point struct{ X, Y int32 }

type MSG struct {
	HWnd     uintptr
	Message  uint32
	WParam   uintptr
	LParam   uintptr
	Time     uint32
	Pt       Point
	LPrivate uint32
}

type WndClassEx struct {
	CbSize        uint32
	Style         uint32
	LpfnWndProc   uintptr
	CbClsExtra    int32
	CbWndExtra    int32
	HInstance     uintptr
	HIcon         uintptr
	HCursor       uintptr
	HbrBackground uintptr
	LpszMenuName  uintptr
	LpszClassName uintptr
	HIconSm       uintptr
}

type PowerBroadcastSetting struct {
	PowerSetting syscall.GUID
	DataLength   uint32
	Data         [4]byte // DWORD per docs: lid 0=closed, 1=opened
}

// Win32 proc handles.

var (
	user32 = syscall.NewLazyDLL("user32.dll")

	procRegisterClassExW            = user32.NewProc("RegisterClassExW")
	procCreateWindowExW             = user32.NewProc("CreateWindowExW")
	procDefWindowProcW              = user32.NewProc("DefWindowProcW")
	procGetMessageW                 = user32.NewProc("GetMessageW")
	procDispatchMessageW            = user32.NewProc("DispatchMessageW")
	procRegisterPowerSettingNotif   = user32.NewProc("RegisterPowerSettingNotification")
	procUnregisterPowerSettingNotif = user32.NewProc("UnregisterPowerSettingNotification")
	procPostMessageW                = user32.NewProc("PostMessageW")

	shell32           = syscall.NewLazyDLL("shell32.dll")
	procShellExecuteW = shell32.NewProc("ShellExecuteW")

	kernel32                     = syscall.NewLazyDLL("kernel32.dll")
	procCreateMutexW             = kernel32.NewProc("CreateMutexW")
	procGetModuleHandleW         = kernel32.NewProc("GetModuleHandleW")
	procGetUserDefaultUILanguage = kernel32.NewProc("GetUserDefaultUILanguage")
)

// LidGuid is the GUID_LIDSWITCH_STATE_CHANGE {BA3E0F4D-B817-4094-A2D1-D56379E6A0F3}.
var LidGuid = syscall.GUID{
	Data1: 0xBA3E0F4D,
	Data2: 0xB817,
	Data3: 0x4094,
	Data4: [8]byte{0xA2, 0xD1, 0xD5, 0x63, 0x79, 0xE6, 0xA0, 0xF3},
}

// CheckSingleInstance reports whether another instance is already running.
func CheckSingleInstance() bool {
	name, _ := syscall.UTF16PtrFromString(`Local\LidReSleepSingleton`)
	_, _, err := procCreateMutexW.Call(0, 0, uintptr(unsafe.Pointer(name)))
	return err == syscall.ERROR_ALREADY_EXISTS
}

// SystemUILanguage returns the lowercase ISO 639-1 code of the system UI
// language when it is one of the languages supported by the UI (zh/en/ja/ko/
// fr/de/es/ru/pt/it), otherwise "en". Detection uses the primary language ID.
func SystemUILanguage() string {
	r, _, _ := procGetUserDefaultUILanguage.Call()
	switch uint16(r) & 0x03FF {
	case 0x04: // LANG_CHINESE
		return "zh"
	case 0x11: // LANG_JAPANESE
		return "ja"
	case 0x12: // LANG_KOREAN
		return "ko"
	case 0x0C: // LANG_FRENCH
		return "fr"
	case 0x07: // LANG_GERMAN
		return "de"
	case 0x0A: // LANG_SPANISH
		return "es"
	case 0x19: // LANG_RUSSIAN
		return "ru"
	case 0x16: // LANG_PORTUGUESE
		return "pt"
	case 0x10: // LANG_ITALIAN
		return "it"
	default:
		return "en"
	}
}

// RegisterClassExW registers the given window class.
func RegisterClassExW(wcex *WndClassEx) (uintptr, error) {
	atom, _, err := procRegisterClassExW.Call(uintptr(unsafe.Pointer(wcex)))
	return atom, err
}

// CreateWindowExW creates a hidden top-level window (never shown).
func CreateWindowExW(className, windowName *uint16) (uintptr, error) {
	hinst, _, _ := procGetModuleHandleW.Call(0)
	hwnd, _, err := procCreateWindowExW.Call(
		0,
		uintptr(unsafe.Pointer(className)),
		uintptr(unsafe.Pointer(windowName)),
		0, // dwStyle: WS_OVERLAPPED (0), never shown
		0, 0, 0, 0,
		0, 0, hinst, 0,
	)
	return hwnd, err
}

// GetMessageW pumps one message; returns 0 on WM_QUIT, 0xFFFFFFFF on error.
func GetMessageW(m *MSG) (uintptr, error) {
	ret, _, err := procGetMessageW.Call(uintptr(unsafe.Pointer(m)), 0, 0, 0)
	return ret, err
}

// DispatchMessageW dispatches a message to its window procedure.
func DispatchMessageW(m *MSG) {
	procDispatchMessageW.Call(uintptr(unsafe.Pointer(m)))
}

// DefWindowProcW calls the default window procedure.
func DefWindowProcW(hwnd uintptr, uMsg uint32, wParam, lParam uintptr) uintptr {
	ret, _, _ := procDefWindowProcW.Call(hwnd, uintptr(uMsg), wParam, lParam)
	return ret
}

// RegisterLidNotify registers for GUID_LIDSWITCH_STATE_CHANGE notifications.
func RegisterLidNotify(hwnd uintptr) (uintptr, error) {
	notify, _, err := procRegisterPowerSettingNotif.Call(
		hwnd,
		uintptr(unsafe.Pointer(&LidGuid)),
		DeviceNotifyWindowHandle,
	)
	return notify, err
}

// UnregisterLidNotify releases a power-notification registration.
func UnregisterLidNotify(h uintptr) {
	if h != 0 {
		procUnregisterPowerSettingNotif.Call(h)
	}
}

// SendScreenOff posts SC_MONITORPOWER to turn the display off, which is the
// user-mode way to trigger S0 Modern Standby on modern laptops. Per MSDN
// (WM_SYSCOMMAND), the 4th argument 2 turns the monitor off. PostMessageW is
// used (async broadcast) so a slow/unresponsive top-level window can never
// block the caller (the engine loop); WM_SYSCOMMAND is a broadcast that every
// window handles independently.
func SendScreenOff() {
	procPostMessageW.Call(HWNDBroadcast, WMSysCommand, SCMonitorPower, MonitorLowPower)
}

// OpenURL opens a URL in the default browser via ShellExecute.
func OpenURL(url string) error {
	verb, _ := syscall.UTF16PtrFromString("open")
	urlPtr, _ := syscall.UTF16PtrFromString(url)
	r, _, _ := procShellExecuteW.Call(
		0,
		uintptr(unsafe.Pointer(verb)),
		uintptr(unsafe.Pointer(urlPtr)),
		0,
		0,
		5, // SW_SHOW
	)
	if r <= 32 {
		return syscall.Errno(r)
	}
	return nil
}
