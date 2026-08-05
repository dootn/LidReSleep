//go:build windows

// Package i18n provides minimal Chinese/English localization. The language comes
// from config.json ("lang"); when unset it is auto-detected from the system UI
// language. T/F are used by the UI and the log.
package i18n

import (
	"fmt"
	"sync/atomic"

	"lidresleep/internal/config"
	"lidresleep/internal/win32"
)

// Lang is the active UI language ("zh" or "en"). Read/written from different
// goroutines (engine logs and the GUI thread), so it is stored atomically.
var Lang atomic.Value // string

// GetLang returns the active UI language ("zh" or "en").
func GetLang() string {
	v, _ := Lang.Load().(string)
	if v == "" {
		return "zh"
	}
	return v
}

// SetLang sets the active UI language.
func SetLang(code string) { Lang.Store(code) }

// Init resolves the UI language early at startup.
func Init() {
	pc := config.Load()
	if pc.Lang == "zh" || pc.Lang == "en" {
		SetLang(pc.Lang)
		return
	}
	if !win32.SystemIsChinese() {
		SetLang("en")
	} else {
		SetLang("zh")
	}
}

// T returns the text for the current language; falls back to Chinese if missing.
func T(key string) string {
	if GetLang() == "en" {
		if s, ok := en[key]; ok {
			return s
		}
	}
	return zh[key]
}

// F formats the current-language text.
func F(key string, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}

var zh = map[string]string{
	"ui.title":             "LidReSleep",
	"menu.file":            "文件",
	"menu.exit":            "退出",
	"menu.language":        "语言",
	"menu.tools":           "工具",
	"menu.testSleep":       "测试睡眠",
	"menu.help":            "帮助",
	"menu.about":           "关于 LidReSleep",
	"menu.github":          "项目主页",
	"menu.checkUpdate":     "检查更新",
	"log.openURLFail":      "打开项目主页失败\n原因：%v",
	"update.title":         "检查更新",
	"update.available":     "发现新版本 %s → %s\n\n是否打开 GitHub 发布页下载最新版？",
	"update.upToDate":      "当前已是最新版本 (v%s)。",
	"update.checkFail":     "检查更新失败\n原因：%v",
	"log.updateChecking":   "正在检查更新…",
	"log.updateAvailable":  "发现新版本: %s",
	"log.updateUpToDate":   "当前已是最新版本。",
	"log.updateCheckFail":  "检查更新失败\n原因：%v",
	"status.group":         "状态",
	"status.stopped":       "● 已停止",
	"status.guarding":      "● 正在守护",
	"btn.startGuard":       "▶ 启动守护",
	"btn.stopGuard":        "■ 停止守护",
	"config.group":         "配置",
	"param.delay":          "入睡延迟(ms)",
	"section.startup":      "启动",
	"section.window":       "窗口",
	"cb.autoboot":          "开机启动",
	"cb.autoguard":         "登录后自动守护",
	"cb.mintotray":         "最小化到托盘",
	"cb.closetotray":       "关闭窗口隐藏到托盘",
	"log.group":            "运行日志",
	"tt.delay":             "合盖状态下被意外唤醒后，等待这么久再次进入睡眠",
	"tt.autoboot":          "随 Windows 登录自动运行（注册表 Run 键）",
	"tt.autoguard":         "启动后直接开始守护并最小化到托盘",
	"tt.mintotray":         "点击最小化按钮时隐藏到系统托盘",
	"tt.closetotray":       "点击窗口关闭按钮时隐藏到托盘而非退出",
	"tray.tooltip":         "LidReSleep",
	"tray.balloon":         "已最小化到系统托盘，点击托盘图标可恢复。",
	"menu.showWindow":      "显示主窗口",
	"notify.running":       "LidReSleep 已在运行。",
	"notify.createWinFail": "创建窗口失败: %s",
	"fatal.prefix":         "致命错误: ",

	"about.title": "关于 LidReSleep",
	"about.body": `LidReSleep v%s
──────────────────────────

保持合盖睡眠 · 唤醒即自动再次入睡

当系统意外从 Modern Standby 唤醒时，若盖子仍关闭，
则自动再次进入睡眠。

功能
· 合盖立即进入睡眠（Modern Standby / S0）
· 被意外唤醒且盖子仍闭合时，再次入睡
· 开盖取消

什么是 Modern Standby？
现代待机设备合盖后系统仍保持低功耗联网，易被网络请求、后台任务、驱动等唤醒而无法深睡，导致发热耗电；本工具检测到盖子仍关闭时会自动让系统再次入睡。

项目主页：https://github.com/dootn/LidReSleep`,

	"log.started":             "LidReSleep v%s 已启动。",
	"log.engineAlready":       "引擎已在运行。",
	"log.engineStarted":       "已开始守护 (delay=%dms)",
	"log.engineStopped":       "已停止守护",
	"log.sepTest":             "--- 睡眠测试 ---",
	"log.evtLidClosed":        "盖子关闭",
	"log.evtLidOpened":        "盖子打开",
	"log.evtSuspend":          "系统挂起",
	"log.evtResume":           "系统被唤醒",
	"log.actSchedule":         "当前盖子关闭，未睡眠，将在 %dms 后执行睡眠操作",
	"log.actReschedule":       "当前盖子关闭，未睡眠，取消未执行的睡眠操作，将在 %dms 后执行睡眠操作",
	"log.actCancel":           "开盖，取消未执行的睡眠操作",
	"log.actSleep":            "执行睡眠操作（关屏触发待机）",
	"log.actSleepRetry":       "睡眠未成功，将重新尝试入睡（第 %d/%d 次）",
	"log.actSleepRetryGiveUp": "多次尝试睡眠仍未成功，已停止自动重试",
	"log.actSleepConfirmed":   "睡眠已生效（计时器被睡眠延迟），按正常流程处理",
	"log.resumeNotifyFail":    "重新注册电源通知失败\n原因：%v",
	"log.getmessageFail":      "GetMessageW 失败\n原因：%v",
	"log.autobootOn":          "已开启开机启动",
	"log.autobootOff":         "已关闭开机启动",
	"log.autobootFail":        "设置开机启动失败\n原因：%v",
	"log.saveConfigFail":      "保存配置失败\n原因：%v",
	"log.trayFail":            "创建系统托盘失败\n原因：%v",
	"log.autoguard":           "已启用登录后自动守护，最小化到系统托盘。",
	"log.langChanged":         "语言已切换为 %s，重启后生效。",
	"lang.zh":                 "中文",
	"lang.en":                 "English",
}

var en = map[string]string{
	"ui.title":             "LidReSleep",
	"menu.file":            "File",
	"menu.exit":            "Exit",
	"menu.language":        "Language",
	"menu.tools":           "Tools",
	"menu.testSleep":       "Test Sleep",
	"menu.help":            "Help",
	"menu.about":           "About LidReSleep",
	"menu.github":          "Project Page",
	"menu.checkUpdate":     "Check for Updates",
	"log.openURLFail":      "Failed to open the project page\nReason: %v",
	"update.title":         "Check for Updates",
	"update.available":     "A new version is available: %s → %s\n\nOpen the GitHub release page to download?",
	"update.upToDate":      "You are running the latest version (v%s).",
	"update.checkFail":     "Failed to check for updates\nReason: %v",
	"log.updateChecking":   "Checking for updates…",
	"log.updateAvailable":  "Update available: %s",
	"log.updateUpToDate":   "You are up to date.",
	"log.updateCheckFail":  "Failed to check for updates\nReason: %v",
	"status.group":         "Status",
	"status.stopped":       "● Stopped",
	"status.guarding":      "● Guarding",
	"btn.startGuard":       "▶ Start Guard",
	"btn.stopGuard":        "■ Stop Guard",
	"config.group":         "Settings",
	"param.delay":          "Sleep delay (ms)",
	"section.startup":      "Startup",
	"section.window":       "Window",
	"cb.autoboot":          "Run at startup",
	"cb.autoguard":         "Auto-guard after login",
	"cb.mintotray":         "Minimize to tray",
	"cb.closetotray":       "Close window to tray",
	"log.group":            "Log",
	"tt.delay":             "When woken unexpectedly with the lid closed, wait this long before sleeping again",
	"tt.autoboot":          "Run automatically at Windows login (registry Run key)",
	"tt.autoguard":         "Start guarding immediately and minimize to tray",
	"tt.mintotray":         "Hide to tray when minimizing",
	"tt.closetotray":       "Hide to tray instead of exiting when closing the window",
	"tray.tooltip":         "LidReSleep — Keep Sleep",
	"tray.balloon":         "Minimized to the system tray. Click the tray icon to restore.",
	"menu.showWindow":      "Show Main Window",
	"notify.running":       "LidReSleep is already running.",
	"notify.createWinFail": "Failed to create window: %s",
	"fatal.prefix":         "Fatal error: ",

	"about.title": "About LidReSleep",
	"about.body": `LidReSleep v%s
──────────────────────────

Keep sleep on lid close · re-sleep when woken

If the system is unexpectedly woken from Modern Standby while the lid is still closed, it will automatically sleep again.

Features
· Sleep immediately when the lid is closed (Modern Standby / S0)
· Re-sleep when woken with the lid still closed
· Cancel when the lid is opened

What is Modern Standby?
Modern-standby devices stay connected at low power after the lid is closed and can be woken by network requests, background tasks, drivers, etc. This tool detects that the lid is still closed and puts the system back to sleep until you open the lid.

Project page: https://github.com/dootn/LidReSleep`,

	"log.started":             "LidReSleep v%s started.",
	"log.engineAlready":       "Guard is already running.",
	"log.engineStarted":       "Guard started (delay=%dms)",
	"log.engineStopped":       "Guard stopped",
	"log.sepTest":             "--- Sleep test ---",
	"log.evtLidClosed":        "Lid closed",
	"log.evtLidOpened":        "Lid opened",
	"log.evtSuspend":          "System suspended",
	"log.evtResume":           "System woken",
	"log.actSchedule":         "Lid closed & awake, sleep in %dms",
	"log.actReschedule":       "Lid closed & awake, cancelled pending sleep, sleep in %dms",
	"log.actCancel":           "Lid opened, cancelled pending sleep",
	"log.actSleep":            "Executing sleep (display-off standby)",
	"log.actSleepRetry":       "Sleep did not take effect, retrying (attempt %d/%d)",
	"log.actSleepRetryGiveUp": "Sleep failed repeatedly, stopped auto-retrying",
	"log.actSleepConfirmed":   "Sleep confirmed (timer delayed by suspend), resuming normal flow",
	"log.resumeNotifyFail":    "Failed to re-register power notification\nReason: %v",
	"log.getmessageFail":      "GetMessageW failed\nReason: %v",
	"log.autobootOn":          "Run at startup enabled",
	"log.autobootOff":         "Run at startup disabled",
	"log.autobootFail":        "Failed to set run-at-startup\nReason: %v",
	"log.saveConfigFail":      "Failed to save config\nReason: %v",
	"log.trayFail":            "Failed to create tray icon\nReason: %v",
	"log.autoguard":           "Auto-guard enabled, minimized to tray.",
	"log.langChanged":         "Language switched to %s. Restart to take effect.",
	"lang.zh":                 "中文",
	"lang.en":                 "English",
}
