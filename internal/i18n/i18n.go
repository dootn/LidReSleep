//go:build windows

// Package i18n provides minimal UI localization. The language comes from
// config.json ("lang"); when unset it is auto-detected from the system UI
// language. Each map falls back to English, then to Chinese. T/F are used by
// the UI and the log.
package i18n

import (
	"fmt"
	"sync/atomic"

	"lidresleep/internal/config"
	"lidresleep/internal/win32"
)

// codes lists the supported UI languages (ISO 639-1), in the order they appear
// in the language menu.
var codes = []string{"zh", "en", "ja", "ko", "fr", "de", "es", "ru", "pt", "it"}

// nativeNames holds each language's name in its own language, shown in the
// language menu and in the switch confirmation message.
var nativeNames = map[string]string{
	"zh": "中文",
	"en": "English",
	"ja": "日本語",
	"ko": "한국어",
	"fr": "Français",
	"de": "Deutsch",
	"es": "Español",
	"ru": "Русский",
	"pt": "Português",
	"it": "Italiano",
}

// Lang is the active UI language code. Read/written from different goroutines
// (engine logs and the GUI thread), so it is stored atomically.
var Lang atomic.Value // string

// Codes returns the supported language codes in menu order.
func Codes() []string {
	return append([]string(nil), codes...)
}

// LangName returns the native name of a language code (e.g. "日本語" for "ja").
func LangName(code string) string {
	if n, ok := nativeNames[code]; ok {
		return n
	}
	return code
}

// IsSupported reports whether code is one of the supported languages.
func IsSupported(code string) bool {
	_, ok := nativeNames[code]
	return ok
}

// GetLang returns the active UI language code; defaults to Chinese.
func GetLang() string {
	v, _ := Lang.Load().(string)
	if v == "" {
		return "zh"
	}
	return v
}

// SetLang sets the active UI language code.
func SetLang(code string) { Lang.Store(code) }

// Init resolves the UI language early at startup.
func Init() {
	pc := config.Load()
	if IsSupported(pc.Lang) {
		SetLang(pc.Lang)
		return
	}
	SetLang(win32.SystemUILanguage())
}

// T returns the text for the current language; falls back to English, then to
// Chinese if missing.
func T(key string) string {
	if m, ok := langMap[GetLang()]; ok {
		if s, ok := m[key]; ok {
			return s
		}
	}
	if s, ok := en[key]; ok {
		return s
	}
	return zh[key]
}

// F formats the current-language text.
func F(key string, args ...interface{}) string {
	return fmt.Sprintf(T(key), args...)
}

var langMap = map[string]map[string]string{
	"zh": zh,
	"en": en,
	"ja": ja,
	"ko": ko,
	"fr": fr,
	"de": de,
	"es": es,
	"ru": ru,
	"pt": pt,
	"it": it,
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
	"log.langChanged":         "语言已切换为 %s。",
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
	"log.langChanged":         "Language switched to %s.",
}

var ja = map[string]string{
	"ui.title":             "LidReSleep",
	"menu.file":            "ファイル",
	"menu.exit":            "終了",
	"menu.language":        "言語",
	"menu.tools":           "ツール",
	"menu.testSleep":       "スリープをテスト",
	"menu.help":            "ヘルプ",
	"menu.about":           "LidReSleep について",
	"menu.github":          "プロジェクトページ",
	"menu.checkUpdate":     "更新を確認",
	"log.openURLFail":      "プロジェクトページを開けませんでした\n理由: %v",
	"update.title":         "更新を確認",
	"update.available":     "新しいバージョンがあります: %s → %s\n\nGitHub のリリースページを開いてダウンロードしますか？",
	"update.upToDate":      "最新バージョン (v%s) を使用しています。",
	"update.checkFail":     "更新の確認に失敗しました\n理由: %v",
	"log.updateChecking":   "更新を確認しています…",
	"log.updateAvailable":  "新しいバージョンがあります: %s",
	"log.updateUpToDate":   "最新版です。",
	"log.updateCheckFail":  "更新の確認に失敗しました\n理由: %v",
	"status.group":         "状態",
	"status.stopped":       "● 停止中",
	"status.guarding":      "● 監視中",
	"btn.startGuard":       "▶ 監視を開始",
	"btn.stopGuard":        "■ 監視を停止",
	"config.group":         "設定",
	"param.delay":          "スリープ遅延(ms)",
	"section.startup":      "起動",
	"section.window":       "ウィンドウ",
	"cb.autoboot":          "起動時に実行",
	"cb.autoguard":         "ログイン後に自動で監視",
	"cb.mintotray":         "最小化でトレイに格納",
	"cb.closetotray":       "閉じるボタンでトレイに格納",
	"log.group":            "ログ",
	"tt.delay":             "蓋が閉じた状態で予期せず起動した場合、この時間待ってから再びスリープします",
	"tt.autoboot":          "Windows ログイン時に自動実行します（レジストリの Run キー）",
	"tt.autoguard":         "起動後すぐに監視を開始し、トレイに最小化します",
	"tt.mintotray":         "最小化ボタンをクリックしたときにシステムトレイへ非表示にします",
	"tt.closetotray":       "ウィンドウを閉じたときに終了せずトレイへ非表示にします",
	"tray.tooltip":         "LidReSleep — スリープを維持",
	"tray.balloon":         "システムトレイに最小化しました。トレイアイコンをクリックして復元できます。",
	"menu.showWindow":      "メインウィンドウを表示",
	"notify.running":       "LidReSleep はすでに実行中です。",
	"notify.createWinFail": "ウィンドウの作成に失敗しました: %s",
	"fatal.prefix":         "致命的エラー: ",

	"about.title": "LidReSleep について",
	"about.body": `LidReSleep v%s
──────────────────────────

蓋を閉じたらスリープ維持 · 起動したら再スリープ

蓋が閉じたまま Modern Standby から予期せず起動した場合は、自動的に再びスリープします。

機能
· 蓋を閉じたらすぐにスリープ（Modern Standby / S0）
· 蓋が閉じたまま起動された場合は再スリープ
· 蓋を開けたらキャンセル

Modern Standby とは？
Modern Standby 対応デバイスは蓋を閉じた後も低消費電力でネットワークに接続したままになり、ネットワークリクエストやバックグラウンドタスク、ドライバなどによって起動され、深いスリープに入れず発熱やバッテリー消費が発生します。このツールは蓋がまだ閉じていることを検出して、蓋を開けるまでシステムを再びスリープさせます。

プロジェクトページ: https://github.com/dootn/LidReSleep`,

	"log.started":             "LidReSleep v%s を起動しました。",
	"log.engineAlready":       "監視はすでに実行中です。",
	"log.engineStarted":       "監視を開始しました (delay=%dms)",
	"log.engineStopped":       "監視を停止しました",
	"log.sepTest":             "--- スリープテスト ---",
	"log.evtLidClosed":        "蓋を閉じました",
	"log.evtLidOpened":        "蓋を開けました",
	"log.evtSuspend":          "システムがスリープしました",
	"log.evtResume":           "システムが復帰しました",
	"log.actSchedule":         "蓋が閉じて起動状態、%dms 後にスリープします",
	"log.actReschedule":       "蓋が閉じて起動状態、保留中のスリープをキャンセルし、%dms 後にスリープします",
	"log.actCancel":           "蓋を開けたため、保留中のスリープをキャンセルしました",
	"log.actSleep":            "スリープを実行（画面オフで待機）",
	"log.actSleepRetry":       "スリープが有効にならなかったため、再試行します（%d/%d 回目）",
	"log.actSleepRetryGiveUp": "スリープが何度も失敗したため、自動再試行を停止しました",
	"log.actSleepConfirmed":   "スリープが有効になりました（サスペンドでタイマーが延期）、通常処理を再開します",
	"log.resumeNotifyFail":    "電源通知の再登録に失敗しました\n理由: %v",
	"log.getmessageFail":      "GetMessageW が失敗しました\n理由: %v",
	"log.autobootOn":          "起動時に実行を有効にしました",
	"log.autobootOff":         "起動時に実行を無効にしました",
	"log.autobootFail":        "起動時に実行の設定に失敗しました\n理由: %v",
	"log.saveConfigFail":      "設定の保存に失敗しました\n理由: %v",
	"log.trayFail":            "トレイアイコンの作成に失敗しました\n理由: %v",
	"log.autoguard":           "自動監視を有効にし、トレイに最小化しました。",
	"log.langChanged":         "言語を %s に変更しました。",
}

var ko = map[string]string{
	"ui.title":             "LidReSleep",
	"menu.file":            "파일",
	"menu.exit":            "종료",
	"menu.language":        "언어",
	"menu.tools":           "도구",
	"menu.testSleep":       "절전 테스트",
	"menu.help":            "도움말",
	"menu.about":           "LidReSleep 정보",
	"menu.github":          "프로젝트 페이지",
	"menu.checkUpdate":     "업데이트 확인",
	"log.openURLFail":      "프로젝트 페이지를 열지 못했습니다\n원인: %v",
	"update.title":         "업데이트 확인",
	"update.available":     "새 버전이 있습니다: %s → %s\n\nGitHub 릴리스 페이지를 열어 다운로드할까요?",
	"update.upToDate":      "최신 버전(v%s)을 사용 중입니다.",
	"update.checkFail":     "업데이트 확인에 실패했습니다\n원인: %v",
	"log.updateChecking":   "업데이트를 확인하는 중…",
	"log.updateAvailable":  "새 버전이 있습니다: %s",
	"log.updateUpToDate":   "최신 버전입니다.",
	"log.updateCheckFail":  "업데이트 확인에 실패했습니다\n원인: %v",
	"status.group":         "상태",
	"status.stopped":       "● 중지됨",
	"status.guarding":      "● 보호 중",
	"btn.startGuard":       "▶ 보호 시작",
	"btn.stopGuard":        "■ 보호 중지",
	"config.group":         "설정",
	"param.delay":          "절전 지연(ms)",
	"section.startup":      "시작",
	"section.window":       "창",
	"cb.autoboot":          "시작 시 실행",
	"cb.autoguard":         "로그인 후 자동 보호",
	"cb.mintotray":         "최소화 시 트레이로",
	"cb.closetotray":       "닫을 때 트레이로",
	"log.group":            "로그",
	"tt.delay":             "덮개가 닫힌 상태에서 예기치 않게 깨어나면 이 시간을 기다린 후 다시 절전 모드로 전환합니다",
	"tt.autoboot":          "Windows 로그인 시 자동 실행(레지스트리 Run 키)",
	"tt.autoguard":         "시작 직후 보호를 시작하고 트레이로 최소화합니다",
	"tt.mintotray":         "최소화 버튼을 클릭하면 시스템 트레이로 숨깁니다",
	"tt.closetotray":       "창을 닫을 때 종료하지 않고 트레이로 숨깁니다",
	"tray.tooltip":         "LidReSleep — 절전 유지",
	"tray.balloon":         "시스템 트레이로 최소화되었습니다. 트레이 아이콘을 클릭하면 복원됩니다.",
	"menu.showWindow":      "기본 창 표시",
	"notify.running":       "LidReSleep이 이미 실행 중입니다.",
	"notify.createWinFail": "창을 만들지 못했습니다: %s",
	"fatal.prefix":         "치명적 오류: ",

	"about.title": "LidReSleep 정보",
	"about.body": `LidReSleep v%s
──────────────────────────

덮개 닫기 시 절전 유지 · 깨어나면 자동 재절전

덮개가 닫힌 상태에서 Modern Standby에서 예기치 않게 깨어나면 자동으로 다시 절전 모드로 전환합니다.

기능
· 덮개를 닫으면 즉시 절전(Modern Standby / S0)
· 덮개가 닫힌 상태에서 깨어나면 재절전
· 덮개를 열면 취소

Modern Standby란?
Modern Standby 지원 기기는 덮개를 닫아도 낮은 전력으로 네트워크에 연결된 상태가 유지되며, 네트워크 요청, 백그라운드 작업, 드라이버 등에 의해 깨어나 깊은 절전에 들어가지 못해 발열과 배터리 소모가 발생합니다. 이 도구는 덮개가 여전히 닫혀 있는지 감지하여 덮개를 열 때까지 시스템을 다시 절전 상태로 만듭니다.

프로젝트 페이지: https://github.com/dootn/LidReSleep`,

	"log.started":             "LidReSleep v%s을(를) 시작했습니다.",
	"log.engineAlready":       "보호가 이미 실행 중입니다.",
	"log.engineStarted":       "보호를 시작했습니다 (delay=%dms)",
	"log.engineStopped":       "보호를 중지했습니다",
	"log.sepTest":             "--- 절전 테스트 ---",
	"log.evtLidClosed":        "덮개가 닫힘",
	"log.evtLidOpened":        "덮개가 열림",
	"log.evtSuspend":          "시스템이 절전됨",
	"log.evtResume":           "시스템이 깨어남",
	"log.actSchedule":         "덮개가 닫혀 있고 깨어 있는 상태, %dms 후 절전 모드로 전환합니다",
	"log.actReschedule":       "덮개가 닫혀 있고 깨어 있는 상태, 대기 중인 절전을 취소하고 %dms 후 절전 모드로 전환합니다",
	"log.actCancel":           "덮개가 열려 대기 중인 절전을 취소했습니다",
	"log.actSleep":            "절전 실행(화면 끄고 대기)",
	"log.actSleepRetry":       "절전이 적용되지 않아 다시 시도합니다(%d/%d회)",
	"log.actSleepRetryGiveUp": "절전이 여러 번 실패하여 자동 재시도를 중지했습니다",
	"log.actSleepConfirmed":   "절전이 적용되었습니다(절전으로 타이머 지연), 정상 흐름 재개",
	"log.resumeNotifyFail":    "전원 알림 재등록에 실패했습니다\n원인: %v",
	"log.getmessageFail":      "GetMessageW 실패\n원인: %v",
	"log.autobootOn":          "시작 시 실행이 활성화되었습니다",
	"log.autobootOff":         "시작 시 실행이 비활성화되었습니다",
	"log.autobootFail":        "시작 시 실행 설정에 실패했습니다\n원인: %v",
	"log.saveConfigFail":      "설정 저장에 실패했습니다\n원인: %v",
	"log.trayFail":            "트레이 아이콘 생성에 실패했습니다\n원인: %v",
	"log.autoguard":           "자동 보호가 활성화되어 트레이로 최소화했습니다.",
	"log.langChanged":         "언어가 %s(으)로 변경되었습니다.",
}

var fr = map[string]string{
	"ui.title":             "LidReSleep",
	"menu.file":            "Fichier",
	"menu.exit":            "Quitter",
	"menu.language":        "Langue",
	"menu.tools":           "Outils",
	"menu.testSleep":       "Tester le sommeil",
	"menu.help":            "Aide",
	"menu.about":           "À propos de LidReSleep",
	"menu.github":          "Page du projet",
	"menu.checkUpdate":     "Rechercher des mises à jour",
	"log.openURLFail":      "Impossible d'ouvrir la page du projet\nRaison : %v",
	"update.title":         "Rechercher des mises à jour",
	"update.available":     "Une nouvelle version est disponible : %s → %s\n\nOuvrir la page des versions GitHub pour télécharger ?",
	"update.upToDate":      "Vous utilisez la dernière version (v%s).",
	"update.checkFail":     "Échec de la recherche de mises à jour\nRaison : %v",
	"log.updateChecking":   "Recherche de mises à jour…",
	"log.updateAvailable":  "Mise à jour disponible : %s",
	"log.updateUpToDate":   "Vous êtes à jour.",
	"log.updateCheckFail":  "Échec de la recherche de mises à jour\nRaison : %v",
	"status.group":         "État",
	"status.stopped":       "● Arrêté",
	"status.guarding":      "● Garde active",
	"btn.startGuard":       "▶ Démarrer la garde",
	"btn.stopGuard":        "■ Arrêter la garde",
	"config.group":         "Paramètres",
	"param.delay":          "Délai de sommeil (ms)",
	"section.startup":      "Démarrage",
	"section.window":       "Fenêtre",
	"cb.autoboot":          "Lancer au démarrage",
	"cb.autoguard":         "Garde auto après connexion",
	"cb.mintotray":         "Réduire dans la zone de notification",
	"cb.closetotray":       "Fermer dans la zone de notification",
	"log.group":            "Journal",
	"tt.delay":             "Lorsqu'il est réveillé de façon inattendue avec l'écran fermé, attendre ce délai avant de se rendormir",
	"tt.autoboot":          "Se lancer automatiquement à la connexion Windows (clé Run du registre)",
	"tt.autoguard":         "Démarrer la garde immédiatement et réduire dans la zone de notification",
	"tt.mintotray":         "Masquer dans la zone de notification lors de la réduction",
	"tt.closetotray":       "Masquer dans la zone de notification au lieu de quitter à la fermeture",
	"tray.tooltip":         "LidReSleep — Maintien du sommeil",
	"tray.balloon":         "Réduit dans la zone de notification. Cliquez sur l'icône pour restaurer.",
	"menu.showWindow":      "Afficher la fenêtre principale",
	"notify.running":       "LidReSleep est déjà en cours d'exécution.",
	"notify.createWinFail": "Échec de création de la fenêtre : %s",
	"fatal.prefix":         "Erreur fatale : ",

	"about.title": "À propos de LidReSleep",
	"about.body": `LidReSleep v%s
──────────────────────────

Gardez le sommeil à la fermeture · se rendormir au réveil

Si le système est réveillé de manière inattendue depuis Modern Standby alors que l'écran est toujours fermé, il se rendormira automatiquement.

Fonctionnalités
· Sommeil immédiat à la fermeture de l'écran (Modern Standby / S0)
· Re-sommeil en cas de réveil avec l'écran encore fermé
· Annulation à l'ouverture de l'écran

Qu'est-ce que Modern Standby ?
Les appareils Modern Standby restent connectés à faible consommation après la fermeture de l'écran et peuvent être réveillés par des requêtes réseau, des tâches en arrière-plan, des pilotes, etc. Cet outil détecte que l'écran est encore fermé et remet le système en sommeil jusqu'à ce que vous ouvriez l'écran.

Page du projet : https://github.com/dootn/LidReSleep`,

	"log.started":             "LidReSleep v%s démarré.",
	"log.engineAlready":       "La garde est déjà en cours d'exécution.",
	"log.engineStarted":       "Garde démarrée (délai=%dms)",
	"log.engineStopped":       "Garde arrêtée",
	"log.sepTest":             "--- Test de sommeil ---",
	"log.evtLidClosed":        "Écran fermé",
	"log.evtLidOpened":        "Écran ouvert",
	"log.evtSuspend":          "Système suspendu",
	"log.evtResume":           "Système réveillé",
	"log.actSchedule":         "Écran fermé et éveillé, sommeil dans %dms",
	"log.actReschedule":       "Écran fermé et éveillé, sommeil en attente annulé, sommeil dans %dms",
	"log.actCancel":           "Écran ouvert, sommeil en attente annulé",
	"log.actSleep":            "Exécution du sommeil (veille à écran éteint)",
	"log.actSleepRetry":       "Le sommeil n'a pas pris effet, nouvelle tentative (essai %d/%d)",
	"log.actSleepRetryGiveUp": "Le sommeil a échoué à plusieurs reprises, arrêt des nouvelles tentatives automatiques",
	"log.actSleepConfirmed":   "Sommeil confirmé (minuteur retardé par la suspension), reprise du flux normal",
	"log.resumeNotifyFail":    "Échec de la réinscription de la notification d'alimentation\nRaison : %v",
	"log.getmessageFail":      "Échec de GetMessageW\nRaison : %v",
	"log.autobootOn":          "Lancement au démarrage activé",
	"log.autobootOff":         "Lancement au démarrage désactivé",
	"log.autobootFail":        "Échec de la configuration du lancement au démarrage\nRaison : %v",
	"log.saveConfigFail":      "Échec de l'enregistrement de la configuration\nRaison : %v",
	"log.trayFail":            "Échec de création de l'icône de la zone de notification\nRaison : %v",
	"log.autoguard":           "Garde auto activée, réduit dans la zone de notification.",
	"log.langChanged":         "Langue changée en %s.",
}

var de = map[string]string{
	"ui.title":             "LidReSleep",
	"menu.file":            "Datei",
	"menu.exit":            "Beenden",
	"menu.language":        "Sprache",
	"menu.tools":           "Extras",
	"menu.testSleep":       "Ruhezustand testen",
	"menu.help":            "Hilfe",
	"menu.about":           "Über LidReSleep",
	"menu.github":          "Projektseite",
	"menu.checkUpdate":     "Nach Updates suchen",
	"log.openURLFail":      "Projektseite konnte nicht geöffnet werden\nGrund: %v",
	"update.title":         "Nach Updates suchen",
	"update.available":     "Eine neue Version ist verfügbar: %s → %s\n\nGitHub-Veröffentlichungsseite öffnen, um herunterzuladen?",
	"update.upToDate":      "Sie verwenden die neueste Version (v%s).",
	"update.checkFail":     "Suche nach Updates fehlgeschlagen\nGrund: %v",
	"log.updateChecking":   "Suche nach Updates…",
	"log.updateAvailable":  "Update verfügbar: %s",
	"log.updateUpToDate":   "Sie sind auf dem neuesten Stand.",
	"log.updateCheckFail":  "Suche nach Updates fehlgeschlagen\nGrund: %v",
	"status.group":         "Status",
	"status.stopped":       "● Gestoppt",
	"status.guarding":      "● Wache aktiv",
	"btn.startGuard":       "▶ Wache starten",
	"btn.stopGuard":        "■ Wache stoppen",
	"config.group":         "Einstellungen",
	"param.delay":          "Ruheverzögerung (ms)",
	"section.startup":      "Start",
	"section.window":       "Fenster",
	"cb.autoboot":          "Beim Start ausführen",
	"cb.autoguard":         "Nach Anmeldung automatisch bewachen",
	"cb.mintotray":         "Beim Minimieren in das Infobereich",
	"cb.closetotray":       "Beim Schließen in das Infobereich",
	"log.group":            "Protokoll",
	"tt.delay":             "Wenn bei geschlossenem Deckel unerwartet aufgeweckt, so lange warten, bevor erneut geschlafen wird",
	"tt.autoboot":          "Automatisch bei Windows-Anmeldung ausführen (Registrierungs-Run-Schlüssel)",
	"tt.autoguard":         "Sofort mit der Bewachung beginnen und in das Infobereich minimieren",
	"tt.mintotray":         "Beim Minimieren im Infobereich ausblenden",
	"tt.closetotray":       "Beim Schließen des Fensters statt Beenden im Infobereich ausblenden",
	"tray.tooltip":         "LidReSleep — Schlaf erhalten",
	"tray.balloon":         "Im Infobereich minimiert. Klicken Sie auf das Symbol, um wiederherzustellen.",
	"menu.showWindow":      "Hauptfenster anzeigen",
	"notify.running":       "LidReSleep läuft bereits.",
	"notify.createWinFail": "Fenster konnte nicht erstellt werden: %s",
	"fatal.prefix":         "Schwerer Fehler: ",

	"about.title": "Über LidReSleep",
	"about.body": `LidReSleep v%s
──────────────────────────

Schlaf beim Schließen des Deckels bewahren · bei Aufwecken erneut schlafen

Wenn das System bei geschlossenem Deckel unerwartet aus dem Modern Standby aufgeweckt wird, schläft es automatisch erneut ein.

Funktionen
· Sofortiger Schlaf beim Schließen des Deckels (Modern Standby / S0)
· Erneuter Schlaf bei Aufwecken mit weiterhin geschlossenem Deckel
· Abbrechen beim Öffnen des Deckels

Was ist Modern Standby?
Modern-Standby-Geräte bleiben nach dem Schließen des Deckels mit geringem Stromverbrauch verbunden und können durch Netzwerkanfragen, Hintergrundaufgaben, Treiber usw. aufgeweckt werden. Dieses Tool erkennt, dass der Deckel noch geschlossen ist, und versetzt das System wieder in den Schlaf, bis Sie den Deckel öffnen.

Projektseite: https://github.com/dootn/LidReSleep`,

	"log.started":             "LidReSleep v%s gestartet.",
	"log.engineAlready":       "Die Wache läuft bereits.",
	"log.engineStarted":       "Wache gestartet (Verzögerung=%dms)",
	"log.engineStopped":       "Wache gestoppt",
	"log.sepTest":             "--- Ruhetest ---",
	"log.evtLidClosed":        "Deckel geschlossen",
	"log.evtLidOpened":        "Deckel geöffnet",
	"log.evtSuspend":          "System angehalten",
	"log.evtResume":           "System aufgeweckt",
	"log.actSchedule":         "Deckel zu & wach, Schlaf in %dms",
	"log.actReschedule":       "Deckel zu & wach, ausstehender Schlaf abgebrochen, Schlaf in %dms",
	"log.actCancel":           "Deckel geöffnet, ausstehender Schlaf abgebrochen",
	"log.actSleep":            "Schlaf ausführen (Standby mit ausgeschaltetem Bildschirm)",
	"log.actSleepRetry":       "Schlaf hat nicht gewirkt, neuer Versuch (Versuch %d/%d)",
	"log.actSleepRetryGiveUp": "Schlaf wiederholt fehlgeschlagen, automatische Wiederholung gestoppt",
	"log.actSleepConfirmed":   "Schlaf bestätigt (Timer durch Anhalten verzögert), normaler Ablauf fortgesetzt",
	"log.resumeNotifyFail":    "Neu-Registrierung der Energiebenachrichtigung fehlgeschlagen\nGrund: %v",
	"log.getmessageFail":      "GetMessageW fehlgeschlagen\nGrund: %v",
	"log.autobootOn":          "Autostart aktiviert",
	"log.autobootOff":         "Autostart deaktiviert",
	"log.autobootFail":        "Autostart konnte nicht gesetzt werden\nGrund: %v",
	"log.saveConfigFail":      "Konfiguration konnte nicht gespeichert werden\nGrund: %v",
	"log.trayFail":            "Infobereich-Symbol konnte nicht erstellt werden\nGrund: %v",
	"log.autoguard":           "Auto-Wache aktiviert, in das Infobereich minimiert.",
	"log.langChanged":         "Sprache auf %s geändert.",
}

var es = map[string]string{
	"ui.title":             "LidReSleep",
	"menu.file":            "Archivo",
	"menu.exit":            "Salir",
	"menu.language":        "Idioma",
	"menu.tools":           "Herramientas",
	"menu.testSleep":       "Probar suspensión",
	"menu.help":            "Ayuda",
	"menu.about":           "Acerca de LidReSleep",
	"menu.github":          "Página del proyecto",
	"menu.checkUpdate":     "Buscar actualizaciones",
	"log.openURLFail":      "No se pudo abrir la página del proyecto\nMotivo: %v",
	"update.title":         "Buscar actualizaciones",
	"update.available":     "Hay una nueva versión disponible: %s → %s\n\n¿Abrir la página de versiones de GitHub para descargarla?",
	"update.upToDate":      "Está usando la última versión (v%s).",
	"update.checkFail":     "No se pudieron buscar actualizaciones\nMotivo: %v",
	"log.updateChecking":   "Buscando actualizaciones…",
	"log.updateAvailable":  "Actualización disponible: %s",
	"log.updateUpToDate":   "Ya está actualizado.",
	"log.updateCheckFail":  "No se pudieron buscar actualizaciones\nMotivo: %v",
	"status.group":         "Estado",
	"status.stopped":       "● Detenido",
	"status.guarding":      "● Protegiendo",
	"btn.startGuard":       "▶ Iniciar protección",
	"btn.stopGuard":        "■ Detener protección",
	"config.group":         "Configuración",
	"param.delay":          "Retardo de suspensión (ms)",
	"section.startup":      "Inicio",
	"section.window":       "Ventana",
	"cb.autoboot":          "Ejecutar al inicio",
	"cb.autoguard":         "Protección automática al iniciar sesión",
	"cb.mintotray":         "Minimizar a la bandeja",
	"cb.closetotray":       "Cerrar en la bandeja",
	"log.group":            "Registro",
	"tt.delay":             "Si se despierta inesperadamente con la tapa cerrada, esperar este tiempo antes de volver a dormir",
	"tt.autoboot":          "Ejecutar automáticamente al iniciar Windows (clave Run del registro)",
	"tt.autoguard":         "Empezar a proteger inmediatamente y minimizar a la bandeja",
	"tt.mintotray":         "Ocultar en la bandeja al minimizar",
	"tt.closetotray":       "Ocultar en la bandeja en lugar de salir al cerrar",
	"tray.tooltip":         "LidReSleep — Mantener suspensión",
	"tray.balloon":         "Minimizado a la bandeja del sistema. Haga clic en el icono para restaurar.",
	"menu.showWindow":      "Mostrar ventana principal",
	"notify.running":       "LidReSleep ya está en ejecución.",
	"notify.createWinFail": "No se pudo crear la ventana: %s",
	"fatal.prefix":         "Error fatal: ",

	"about.title": "Acerca de LidReSleep",
	"about.body": `LidReSleep v%s
──────────────────────────

Mantener la suspensión al cerrar la tapa · volver a dormir al despertar

Si el sistema se despierta inesperadamente desde Modern Standby con la tapa aún cerrada, volverá a dormirse automáticamente.

Funciones
· Suspensión inmediata al cerrar la tapa (Modern Standby / S0)
· Volver a dormir si se despierta con la tapa aún cerrada
· Cancelar al abrir la tapa

¿Qué es Modern Standby?
Los dispositivos con Modern Standby permanecen conectados a bajo consumo tras cerrar la tapa y pueden ser despertados por solicitudes de red, tareas en segundo plano, controladores, etc. Esta herramienta detecta que la tapa sigue cerrada y vuelve a poner el sistema en suspensión hasta que abra la tapa.

Página del proyecto: https://github.com/dootn/LidReSleep`,

	"log.started":             "LidReSleep v%s iniciado.",
	"log.engineAlready":       "La protección ya está en ejecución.",
	"log.engineStarted":       "Protección iniciada (retardo=%dms)",
	"log.engineStopped":       "Protección detenida",
	"log.sepTest":             "--- Prueba de suspensión ---",
	"log.evtLidClosed":        "Tapa cerrada",
	"log.evtLidOpened":        "Tapa abierta",
	"log.evtSuspend":          "Sistema suspendido",
	"log.evtResume":           "Sistema despertado",
	"log.actSchedule":         "Tapa cerrada y despierto, suspensión en %dms",
	"log.actReschedule":       "Tapa cerrada y despierto, suspensión pendiente cancelada, suspensión en %dms",
	"log.actCancel":           "Tapa abierta, suspensión pendiente cancelada",
	"log.actSleep":            "Ejecutando suspensión (standby con pantalla apagada)",
	"log.actSleepRetry":       "La suspensión no surtió efecto, reintentando (intento %d/%d)",
	"log.actSleepRetryGiveUp": "La suspensión falló repetidamente, se detuvieron los reintentos automáticos",
	"log.actSleepConfirmed":   "Suspensión confirmada (temporizador retrasado por la suspensión), reanudando el flujo normal",
	"log.resumeNotifyFail":    "No se pudo volver a registrar la notificación de energía\nMotivo: %v",
	"log.getmessageFail":      "GetMessageW falló\nMotivo: %v",
	"log.autobootOn":          "Ejecución al inicio activada",
	"log.autobootOff":         "Ejecución al inicio desactivada",
	"log.autobootFail":        "No se pudo configurar la ejecución al inicio\nMotivo: %v",
	"log.saveConfigFail":      "No se pudo guardar la configuración\nMotivo: %v",
	"log.trayFail":            "No se pudo crear el icono de la bandeja\nMotivo: %v",
	"log.autoguard":           "Protección automática activada, minimizado a la bandeja.",
	"log.langChanged":         "Idioma cambiado a %s.",
}

var ru = map[string]string{
	"ui.title":             "LidReSleep",
	"menu.file":            "Файл",
	"menu.exit":            "Выход",
	"menu.language":        "Язык",
	"menu.tools":           "Инструменты",
	"menu.testSleep":       "Тест сна",
	"menu.help":            "Справка",
	"menu.about":           "О LidReSleep",
	"menu.github":          "Страница проекта",
	"menu.checkUpdate":     "Проверить обновления",
	"log.openURLFail":      "Не удалось открыть страницу проекта\nПричина: %v",
	"update.title":         "Проверить обновления",
	"update.available":     "Доступна новая версия: %s → %s\n\nОткрыть страницу релизов GitHub для загрузки?",
	"update.upToDate":      "У вас установлена последняя версия (v%s).",
	"update.checkFail":     "Не удалось проверить обновления\nПричина: %v",
	"log.updateChecking":   "Проверка обновлений…",
	"log.updateAvailable":  "Доступно обновление: %s",
	"log.updateUpToDate":   "У вас актуальная версия.",
	"log.updateCheckFail":  "Не удалось проверить обновления\nПричина: %v",
	"status.group":         "Статус",
	"status.stopped":       "● Остановлено",
	"status.guarding":      "● Охрана активна",
	"btn.startGuard":       "▶ Запустить охрану",
	"btn.stopGuard":        "■ Остановить охрану",
	"config.group":         "Настройки",
	"param.delay":          "Задержка сна (мс)",
	"section.startup":      "Запуск",
	"section.window":       "Окно",
	"cb.autoboot":          "Запуск при входе",
	"cb.autoguard":         "Автоохрана после входа",
	"cb.mintotray":         "Сворачивать в трей",
	"cb.closetotray":       "Закрывать в трей",
	"log.group":            "Журнал",
	"tt.delay":             "Если система проснулась неожиданно при закрытой крышке, подождать это время перед повторным сном",
	"tt.autoboot":          "Автоматический запуск при входе в Windows (ключ Run в реестре)",
	"tt.autoguard":         "Сразу начать охрану и свернуть в трей",
	"tt.mintotray":         "Скрывать в трей при сворачивании",
	"tt.closetotray":       "Скрывать в трей вместо выхода при закрытии окна",
	"tray.tooltip":         "LidReSleep — поддержание сна",
	"tray.balloon":         "Свёрнуто в системный трей. Щёлкните значок в трее, чтобы восстановить.",
	"menu.showWindow":      "Показать главное окно",
	"notify.running":       "LidReSleep уже запущен.",
	"notify.createWinFail": "Не удалось создать окно: %s",
	"fatal.prefix":         "Критическая ошибка: ",

	"about.title": "О LidReSleep",
	"about.body": `LidReSleep v%s
──────────────────────────

Поддержание сна при закрытой крышке · повторный сон при пробуждении

Если система неожиданно проснулась из Modern Standby, а крышка всё ещё закрыта, она снова уснёт автоматически.

Возможности
· Немедленный сон при закрытии крышки (Modern Standby / S0)
· Повторный сон при пробуждении с закрытой крышкой
· Отмена при открытии крышки

Что такое Modern Standby?
Устройства Modern Standby после закрытия крышки остаются подключёнными с низким энергопотреблением и могут быть разбужены сетевыми запросами, фоновыми задачами, драйверами и т. п. Этот инструмент определяет, что крышка всё ещё закрыта, и снова переводит систему в сон, пока вы не откроете крышку.

Страница проекта: https://github.com/dootn/LidReSleep`,

	"log.started":             "LidReSleep v%s запущен.",
	"log.engineAlready":       "Охрана уже запущена.",
	"log.engineStarted":       "Охрана запущена (задержка=%dмс)",
	"log.engineStopped":       "Охрана остановлена",
	"log.sepTest":             "--- Тест сна ---",
	"log.evtLidClosed":        "Крышка закрыта",
	"log.evtLidOpened":        "Крышка открыта",
	"log.evtSuspend":          "Система приостановлена",
	"log.evtResume":           "Система пробудилась",
	"log.actSchedule":         "Крышка закрыта и система активна, сон через %dмс",
	"log.actReschedule":       "Крышка закрыта и система активна, ожидавшийся сон отменён, сон через %dмс",
	"log.actCancel":           "Крышка открыта, ожидавшийся сон отменён",
	"log.actSleep":            "Выполнение сна (резервный режим с выключенным экраном)",
	"log.actSleepRetry":       "Сон не сработал, повторная попытка (попытка %d/%d)",
	"log.actSleepRetryGiveUp": "Сон многократно не срабатывал, автоматические повторы остановлены",
	"log.actSleepConfirmed":   "Сон подтверждён (таймер отложен приостановкой), возобновление обычного режима",
	"log.resumeNotifyFail":    "Не удалось перерегистрировать уведомление о питании\nПричина: %v",
	"log.getmessageFail":      "Ошибка GetMessageW\nПричина: %v",
	"log.autobootOn":          "Запуск при входе включён",
	"log.autobootOff":         "Запуск при входе отключён",
	"log.autobootFail":        "Не удалось настроить запуск при входе\nПричина: %v",
	"log.saveConfigFail":      "Не удалось сохранить конфигурацию\nПричина: %v",
	"log.trayFail":            "Не удалось создать значок в трее\nПричина: %v",
	"log.autoguard":           "Автоохрана включена, свёрнуто в трей.",
	"log.langChanged":         "Язык изменён на %s.",
}

var pt = map[string]string{
	"ui.title":             "LidReSleep",
	"menu.file":            "Arquivo",
	"menu.exit":            "Sair",
	"menu.language":        "Idioma",
	"menu.tools":           "Ferramentas",
	"menu.testSleep":       "Testar suspensão",
	"menu.help":            "Ajuda",
	"menu.about":           "Sobre o LidReSleep",
	"menu.github":          "Página do projeto",
	"menu.checkUpdate":     "Verificar atualizações",
	"log.openURLFail":      "Não foi possível abrir a página do projeto\nMotivo: %v",
	"update.title":         "Verificar atualizações",
	"update.available":     "Uma nova versão está disponível: %s → %s\n\nAbrir a página de lançamentos do GitHub para baixar?",
	"update.upToDate":      "Você está usando a versão mais recente (v%s).",
	"update.checkFail":     "Falha ao verificar atualizações\nMotivo: %v",
	"log.updateChecking":   "Verificando atualizações…",
	"log.updateAvailable":  "Atualização disponível: %s",
	"log.updateUpToDate":   "Você está atualizado.",
	"log.updateCheckFail":  "Falha ao verificar atualizações\nMotivo: %v",
	"status.group":         "Estado",
	"status.stopped":       "● Parado",
	"status.guarding":      "● Protegendo",
	"btn.startGuard":       "▶ Iniciar proteção",
	"btn.stopGuard":        "■ Parar proteção",
	"config.group":         "Configurações",
	"param.delay":          "Atraso de suspensão (ms)",
	"section.startup":      "Inicialização",
	"section.window":       "Janela",
	"cb.autoboot":          "Executar na inicialização",
	"cb.autoguard":         "Proteção automática após login",
	"cb.mintotray":         "Minimizar para a bandeja",
	"cb.closetotray":       "Fechar para a bandeja",
	"log.group":            "Registro",
	"tt.delay":             "Se for ativado inesperadamente com a tampa fechada, aguardar este tempo antes de suspender novamente",
	"tt.autoboot":          "Executar automaticamente no login do Windows (chave Run do registro)",
	"tt.autoguard":         "Iniciar a proteção imediatamente e minimizar para a bandeja",
	"tt.mintotray":         "Ocultar na bandeja ao minimizar",
	"tt.closetotray":       "Ocultar na bandeja em vez de sair ao fechar a janela",
	"tray.tooltip":         "LidReSleep — Manter suspensão",
	"tray.balloon":         "Minimizado para a bandeja do sistema. Clique no ícone para restaurar.",
	"menu.showWindow":      "Mostrar janela principal",
	"notify.running":       "O LidReSleep já está em execução.",
	"notify.createWinFail": "Falha ao criar a janela: %s",
	"fatal.prefix":         "Erro fatal: ",

	"about.title": "Sobre o LidReSleep",
	"about.body": `LidReSleep v%s
──────────────────────────

Manter a suspensão ao fechar a tampa · voltar a suspender ao ser ativado

Se o sistema for ativado inesperadamente do Modern Standby enquanto a tampa ainda estiver fechada, ele voltará a suspender automaticamente.

Recursos
· Suspensão imediata ao fechar a tampa (Modern Standby / S0)
· Voltar a suspender se for ativado com a tampa ainda fechada
· Cancelar ao abrir a tampa

O que é Modern Standby?
Dispositivos com Modern Standby permanecem conectados em baixo consumo após o fechamento da tampa e podem ser ativados por solicitações de rede, tarefas em segundo plano, drivers etc. Esta ferramenta detecta que a tampa ainda está fechada e coloca o sistema novamente em suspensão até você abrir a tampa.

Página do projeto: https://github.com/dootn/LidReSleep`,

	"log.started":             "LidReSleep v%s iniciado.",
	"log.engineAlready":       "A proteção já está em execução.",
	"log.engineStarted":       "Proteção iniciada (atraso=%dms)",
	"log.engineStopped":       "Proteção interrompida",
	"log.sepTest":             "--- Teste de suspensão ---",
	"log.evtLidClosed":        "Tampa fechada",
	"log.evtLidOpened":        "Tampa aberta",
	"log.evtSuspend":          "Sistema suspenso",
	"log.evtResume":           "Sistema ativado",
	"log.actSchedule":         "Tampa fechada e ativo, suspensão em %dms",
	"log.actReschedule":       "Tampa fechada e ativo, suspensão pendente cancelada, suspensão em %dms",
	"log.actCancel":           "Tampa aberta, suspensão pendente cancelada",
	"log.actSleep":            "Executando suspensão (standby com tela desligada)",
	"log.actSleepRetry":       "A suspensão não surtiu efeito, tentando novamente (tentativa %d/%d)",
	"log.actSleepRetryGiveUp": "A suspensão falhou repetidamente, interrompidas as novas tentativas automáticas",
	"log.actSleepConfirmed":   "Suspensão confirmada (temporizador atrasado pela suspensão), retomando o fluxo normal",
	"log.resumeNotifyFail":    "Falha ao registrar novamente a notificação de energia\nMotivo: %v",
	"log.getmessageFail":      "GetMessageW falhou\nMotivo: %v",
	"log.autobootOn":          "Execução na inicialização ativada",
	"log.autobootOff":         "Execução na inicialização desativada",
	"log.autobootFail":        "Falha ao definir execução na inicialização\nMotivo: %v",
	"log.saveConfigFail":      "Falha ao salvar a configuração\nMotivo: %v",
	"log.trayFail":            "Falha ao criar o ícone da bandeja\nMotivo: %v",
	"log.autoguard":           "Proteção automática ativada, minimizado para a bandeja.",
	"log.langChanged":         "Idioma alterado para %s.",
}

var it = map[string]string{
	"ui.title":             "LidReSleep",
	"menu.file":            "File",
	"menu.exit":            "Esci",
	"menu.language":        "Lingua",
	"menu.tools":           "Strumenti",
	"menu.testSleep":       "Prova sospensione",
	"menu.help":            "Aiuto",
	"menu.about":           "Informazioni su LidReSleep",
	"menu.github":          "Pagina del progetto",
	"menu.checkUpdate":     "Verifica aggiornamenti",
	"log.openURLFail":      "Impossibile aprire la pagina del progetto\nMotivo: %v",
	"update.title":         "Verifica aggiornamenti",
	"update.available":     "È disponibile una nuova versione: %s → %s\n\nAprire la pagina delle release di GitHub per scaricarla?",
	"update.upToDate":      "Si sta utilizzando l'ultima versione (v%s).",
	"update.checkFail":     "Verifica degli aggiornamenti non riuscita\nMotivo: %v",
	"log.updateChecking":   "Verifica degli aggiornamenti in corso…",
	"log.updateAvailable":  "Aggiornamento disponibile: %s",
	"log.updateUpToDate":   "Si è aggiornati.",
	"log.updateCheckFail":  "Verifica degli aggiornamenti non riuscita\nMotivo: %v",
	"status.group":         "Stato",
	"status.stopped":       "● Arrestato",
	"status.guarding":      "● Protezione attiva",
	"btn.startGuard":       "▶ Avvia protezione",
	"btn.stopGuard":        "■ Arresta protezione",
	"config.group":         "Impostazioni",
	"param.delay":          "Ritardo sospensione (ms)",
	"section.startup":      "Avvio",
	"section.window":       "Finestra",
	"cb.autoboot":          "Esegui all'avvio",
	"cb.autoguard":         "Protezione automatica dopo l'accesso",
	"cb.mintotray":         "Riduci a icona nella barra delle applicazioni",
	"cb.closetotray":       "Chiudi nella barra delle applicazioni",
	"log.group":            "Registro",
	"tt.delay":             "Se riattivato inaspettatamente con lo schermo chiuso, attendere questo tempo prima di rientrare in sospensione",
	"tt.autoboot":          "Esegui automaticamente all'accesso di Windows (chiave Run del registro)",
	"tt.autoguard":         "Avvia subito la protezione e riduci a icona nella barra delle applicazioni",
	"tt.mintotray":         "Nascondi nella barra delle applicazioni durante la riduzione a icona",
	"tt.closetotray":       "Nascondi nella barra delle applicazioni invece di uscire alla chiusura",
	"tray.tooltip":         "LidReSleep — Mantieni sospensione",
	"tray.balloon":         "Ridotto a icona nella barra delle applicazioni. Fare clic sull'icona per ripristinare.",
	"menu.showWindow":      "Mostra finestra principale",
	"notify.running":       "LidReSleep è già in esecuzione.",
	"notify.createWinFail": "Impossibile creare la finestra: %s",
	"fatal.prefix":         "Errore fatale: ",

	"about.title": "Informazioni su LidReSleep",
	"about.body": `LidReSleep v%s
──────────────────────────

Mantieni la sospensione alla chiusura · ri-sospensione al risveglio

Se il sistema viene riattivato inaspettatamente da Modern Standby mentre lo schermo è ancora chiuso, si riaddormenterà automaticamente.

Funzionalità
· Sospensione immediata alla chiusura dello schermo (Modern Standby / S0)
· Ri-sospensione se riattivato con lo schermo ancora chiuso
· Annulla all'apertura dello schermo

Che cos'è Modern Standby?
I dispositivi Modern Standby restano connessi a basso consumo dopo la chiusura dello schermo e possono essere riattivati da richieste di rete, attività in background, driver, ecc. Questo strumento rileva che lo schermo è ancora chiuso e riporta il sistema in sospensione fino all'apertura dello schermo.

Pagina del progetto: https://github.com/dootn/LidReSleep`,

	"log.started":             "LidReSleep v%s avviato.",
	"log.engineAlready":       "La protezione è già in esecuzione.",
	"log.engineStarted":       "Protezione avviata (ritardo=%dms)",
	"log.engineStopped":       "Protezione arrestata",
	"log.sepTest":             "--- Prova sospensione ---",
	"log.evtLidClosed":        "Schermo chiuso",
	"log.evtLidOpened":        "Schermo aperto",
	"log.evtSuspend":          "Sistema sospeso",
	"log.evtResume":           "Sistema riattivato",
	"log.actSchedule":         "Schermo chiuso e attivo, sospensione tra %dms",
	"log.actReschedule":       "Schermo chiuso e attivo, sospensione in sospeso annullata, sospensione tra %dms",
	"log.actCancel":           "Schermo aperto, sospensione in sospeso annullata",
	"log.actSleep":            "Esecuzione sospensione (standby a schermo spento)",
	"log.actSleepRetry":       "La sospensione non ha avuto effetto, nuovo tentativo (tentativo %d/%d)",
	"log.actSleepRetryGiveUp": "La sospensione è fallita ripetutamente, tentativi automatici interrotti",
	"log.actSleepConfirmed":   "Sospensione confermata (timer ritardato dalla sospensione), ripresa del flusso normale",
	"log.resumeNotifyFail":    "Re-registrazione della notifica di alimentazione non riuscita\nMotivo: %v",
	"log.getmessageFail":      "GetMessageW non riuscito\nMotivo: %v",
	"log.autobootOn":          "Esecuzione all'avvio attivata",
	"log.autobootOff":         "Esecuzione all'avvio disattivata",
	"log.autobootFail":        "Impostazione dell'esecuzione all'avvio non riuscita\nMotivo: %v",
	"log.saveConfigFail":      "Salvataggio della configurazione non riuscito\nMotivo: %v",
	"log.trayFail":            "Creazione dell'icona nella barra delle applicazioni non riuscita\nMotivo: %v",
	"log.autoguard":           "Protezione automatica attivata, ridotto a icona nella barra delle applicazioni.",
	"log.langChanged":         "Lingua cambiata in %s.",
}
