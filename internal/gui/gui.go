//go:build windows

// Package gui implements the walk panel: menus, status, settings, tray and the
// log box, wiring config/engine/i18n/log/registry together.
package gui

import (
	"bytes"
	"embed"
	"encoding/binary"
	"errors"
	"image"
	"image/png"
	"os"
	"runtime"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	"github.com/lxn/walk"
	. "github.com/lxn/walk/declarative"
	"github.com/lxn/win"

	"lidresleep/internal/config"
	"lidresleep/internal/engine"
	"lidresleep/internal/i18n"
	"lidresleep/internal/log"
	"lidresleep/internal/registry"
	"lidresleep/internal/update"
	"lidresleep/internal/win32"
)

const appVersion = "1.0.0"

//go:embed icon.ico
var iconFS embed.FS

// loadIcon extracts the largest PNG frame from the embedded .ico and builds a
// walk icon from it. This path avoids the unreliable Windows resource-loading
// APIs in lxn/walk.
func loadIcon() (*walk.Icon, error) {
	b, err := iconFS.ReadFile("icon.ico")
	if err != nil {
		return nil, err
	}
	img, err := iconFrameFromICO(b)
	if err != nil {
		return nil, err
	}
	return walk.NewIconFromImage(img)
}

// iconFrameFromICO extracts the largest PNG-compressed frame from an .ico blob.
func iconFrameFromICO(b []byte) (image.Image, error) {
	if len(b) < 6 || b[0] != 0 || b[1] != 0 || binary.LittleEndian.Uint16(b[2:4]) != 1 {
		return nil, errors.New("invalid ico header")
	}
	count := int(binary.LittleEndian.Uint16(b[4:6]))
	bestSize, bestOff, bestLen := 0, -1, 0
	for i := 0; i < count; i++ {
		e := 6 + i*16
		if e+16 > len(b) {
			break
		}
		w, h := int(b[e]), int(b[e+1])
		if w == 0 {
			w = 256
		}
		if h == 0 {
			h = 256
		}
		size := int(binary.LittleEndian.Uint32(b[e+8 : e+12]))
		off := int(binary.LittleEndian.Uint32(b[e+12 : e+16]))
		if w*h > bestSize && off+size <= len(b) {
			bestSize, bestOff, bestLen = w*h, off, size
		}
	}
	if bestOff < 0 {
		return nil, errors.New("no usable icon frame")
	}
	return png.Decode(bytes.NewReader(b[bestOff : bestOff+bestLen]))
}

// maxLogLen caps the log box text length (bytes); oldest lines are dropped
// beyond it.
const maxLogLen = 200 * 1024

var (
	uiMu      sync.Mutex
	mainWin   *walk.MainWindow
	logBox    *walk.TextEdit
	delayEdit *walk.NumberEdit

	startStopBtn *walk.PushButton
	statusLbl    *walk.Label

	autoBootCb  *walk.CheckBox
	autoGuardCb *walk.CheckBox
	minTrayCb   *walk.CheckBox
	closeTrayCb *walk.CheckBox

	// Text-bearing widgets, kept so the language can be re-applied live.
	fileMenuAct  *walk.Action
	exitAct      *walk.Action
	toolsMenuAct *walk.Action
	testSleepAct *walk.Action
	langMenuAct  *walk.Action
	helpMenuAct  *walk.Action
	checkUpdAct  *walk.Action
	githubAct    *walk.Action
	aboutAct     *walk.Action

	statusGroup *walk.GroupBox
	configGroup *walk.GroupBox
	logGroup    *walk.GroupBox
	delayLbl    *walk.Label
	startupLbl  *walk.Label
	windowLbl   *walk.Label

	langActs []*walk.Action

	tray             *walk.NotifyIcon
	trayOpenAct      *walk.Action
	trayExitAct      *walk.Action
	trayBalloonShown bool
	quitting         bool // user-initiated exit (tray/menu); allows the window to actually close

	// savedConfig is the last configuration persisted to disk; applyUiConfig
	// compares against it and skips the file write when nothing changed.
	savedConfig config.PersistConfig
)

var (
	statusRun     = walk.RGB(0, 150, 60)
	statusStopRed = walk.RGB(220, 60, 50)
	descColor     = walk.RGB(130, 130, 130)
)

// uiLogWriter appends log text to the walk log box (callable from any goroutine,
// thread-safe). Windows edit controls need \r\n line endings; converted here,
// and it auto-scrolls to the bottom.
type uiLogWriter struct{}

func (uiLogWriter) Write(p []byte) (int, error) {
	uiMu.Lock()
	w, lb := mainWin, logBox
	uiMu.Unlock()
	if w == nil || lb == nil {
		return len(p), nil
	}
	s := strings.ReplaceAll(string(p), "\r\n", "\n")
	s = strings.ReplaceAll(s, "\n", "\r\n")
	w.Synchronize(func() {
		text := lb.Text()
		if len(text)+len(s) > maxLogLen {
			// Drop whole oldest lines until under the byte cap. Walk back to a
			// UTF-8 rune boundary so a multi-byte character is never sliced in
			// half (TextLength is UTF-16 chars, so it must not be mixed with
			// len() byte math here).
			trim := len(text) + len(s) - maxLogLen
			for trim > 0 && !utf8.RuneStart(text[trim]) {
				trim--
			}
			if idx := strings.IndexByte(text[trim:], '\n'); idx >= 0 {
				trim += idx + 1
			}
			lb.SetText(text[trim:] + s)
		} else {
			lb.AppendText(s)
		}
		lb.SetTextSelection(lb.TextLength(), lb.TextLength())
	})
	return len(p), nil
}

// Notify shows a message box.
func Notify(title, text string) {
	uiMu.Lock()
	owner := mainWin
	uiMu.Unlock()
	style := walk.MsgBoxOK | walk.MsgBoxIconInformation
	walk.MsgBox(owner, title, text, style)
}

// showAbout shows the About dialog.
func showAbout() {
	Notify(i18n.T("about.title"), i18n.F("about.body", appVersion))
}

// openProjectPage opens the project repository in the default browser.
func openProjectPage() {
	if err := win32.OpenURL("https://github.com/dootn/LidReSleep"); err != nil {
		log.Error(i18n.F("log.openURLFail", err))
	}
}

// openReleasePage opens the latest release page in the default browser.
func openReleasePage() {
	if err := win32.OpenURL("https://github.com/dootn/LidReSleep/releases/latest"); err != nil {
		log.Error(i18n.F("log.openURLFail", err))
	}
}

// checkForUpdates queries GitHub for the latest release in the background (the
// network call must never block the UI thread) and reports the result in a
// dialog on the UI thread.
func checkForUpdates() {
	uiMu.Lock()
	w := mainWin
	uiMu.Unlock()
	if w == nil {
		return
	}
	log.Print(i18n.T("log.updateChecking"))
	go func() {
		res, err := update.Check(appVersion, 10*time.Second)
		w.Synchronize(func() {
			if err != nil {
				log.Error(i18n.F("log.updateCheckFail", err))
				Notify(i18n.T("update.title"), i18n.F("update.checkFail", err))
				return
			}
			if res.Update {
				log.Info(i18n.F("log.updateAvailable", res.Latest))
				if walk.MsgBox(w, i18n.T("update.title"),
					i18n.F("update.available", appVersion, res.Latest),
					walk.MsgBoxYesNo|walk.MsgBoxIconInformation) == walk.DlgCmdYes {
					openReleasePage()
				}
			} else {
				log.Info(i18n.T("log.updateUpToDate"))
				Notify(i18n.T("update.title"), i18n.F("update.upToDate", appVersion))
			}
		})
	}()
}

// createTray creates the tray icon and its context menu; returns an error on
// failure (caller degrades gracefully).
func createTray() error {
	icon, err := loadIcon()
	if err != nil {
		return err
	}
	mainWin.SetIcon(icon)

	t, err := walk.NewNotifyIcon(mainWin)
	if err != nil {
		return err
	}
	if err := t.SetIcon(icon); err != nil {
		t.Dispose()
		return err
	}
	if err := t.SetToolTip(i18n.T("tray.tooltip")); err != nil {
		t.Dispose()
		return err
	}

	menu := t.ContextMenu()
	trayOpenAct = walk.NewAction()
	trayOpenAct.SetText(i18n.T("menu.showWindow"))
	trayOpenAct.Triggered().Attach(showMainWindow)
	menu.Actions().Add(trayOpenAct)
	menu.Actions().Add(walk.NewSeparatorAction())
	trayExitAct = walk.NewAction()
	trayExitAct.SetText(i18n.T("menu.exit"))
	trayExitAct.Triggered().Attach(quitApp)
	menu.Actions().Add(trayExitAct)

	t.MouseDown().Attach(func(x, y int, button walk.MouseButton) {
		if button == walk.LeftButton {
			showMainWindow()
		}
	})

	if err := t.SetVisible(true); err != nil {
		t.Dispose()
		return err
	}
	tray = t
	return nil
}

// showMainWindow restores and shows the main window from the tray.
func showMainWindow() {
	uiMu.Lock()
	w := mainWin
	uiMu.Unlock()
	if w == nil {
		return
	}
	w.Show()
	win.ShowWindow(w.Handle(), win.SW_RESTORE)
}

// minimizeToTray hides the main window to the tray (balloon shown once).
func minimizeToTray() {
	uiMu.Lock()
	w := mainWin
	uiMu.Unlock()
	if w == nil {
		return
	}
	w.Hide()
	if !trayBalloonShown && tray != nil {
		trayBalloonShown = true
		tray.ShowInfo(i18n.T("ui.title"), i18n.T("tray.balloon"))
	}
}

// RunUI builds and runs the walk panel; returns before the process exits.
func RunUI() {
	runtime.LockOSThread()

	log.SetSink(uiLogWriter{})
	log.SetFatalHandler(func(text string) {
		Notify(i18n.T("ui.title"), i18n.T("fatal.prefix")+text)
	})

	pc := config.Load()
	savedConfig = pc

	runErr := MainWindow{
		AssignTo: &mainWin,
		Title:    i18n.T("ui.title"),
		MinSize:  Size{Width: 560, Height: 640},
		Size:     Size{Width: 600, Height: 740},
		Layout:   VBox{Margins: Margins{Left: 14, Top: 10, Right: 14, Bottom: 12}, Spacing: 8},
		MenuItems: []MenuItem{
			Menu{
				AssignActionTo: &fileMenuAct,
				Text:           i18n.T("menu.file"),
				Items: []MenuItem{
					Action{AssignTo: &exitAct, Text: i18n.T("menu.exit"), OnTriggered: quitApp},
				},
			},
			Menu{
				AssignActionTo: &toolsMenuAct,
				Text:           i18n.T("menu.tools"),
				Items: []MenuItem{
					Action{AssignTo: &testSleepAct, Text: i18n.T("menu.testSleep"), OnTriggered: engine.SleepNow},
				},
			},
			Menu{
				AssignActionTo: &langMenuAct,
				Text:           i18n.T("menu.language"),
				Items:          langMenuItems(),
			},
			Menu{
				AssignActionTo: &helpMenuAct,
				Text:           i18n.T("menu.help"),
				Items: []MenuItem{
					Action{AssignTo: &checkUpdAct, Text: i18n.T("menu.checkUpdate"), OnTriggered: checkForUpdates},
					Action{AssignTo: &githubAct, Text: i18n.T("menu.github"), OnTriggered: openProjectPage},
					Action{AssignTo: &aboutAct, Text: i18n.T("menu.about"), OnTriggered: showAbout},
				},
			},
		},
		Children: []Widget{
			// Status
			GroupBox{
				AssignTo: &statusGroup,
				Title:    i18n.T("status.group"),
				Layout:   VBox{Margins: Margins{Left: 10, Top: 6, Right: 10, Bottom: 8}, Spacing: 6},
				Children: []Widget{
					Composite{
						Layout: HBox{Spacing: 8},
						Children: []Widget{
							Label{AssignTo: &statusLbl, Text: i18n.T("status.stopped"), Font: Font{Family: "Microsoft YaHei UI", PointSize: 12, Bold: true}},
							HSpacer{},
							PushButton{AssignTo: &startStopBtn, Text: i18n.T("btn.startGuard"), MinSize: Size{Width: 130, Height: 32}, OnClicked: toggleEngine},
						},
					},
				},
			},
			// Settings
			GroupBox{
				AssignTo: &configGroup,
				Title:    i18n.T("config.group"),
				Layout:   VBox{Margins: Margins{Left: 10, Top: 6, Right: 10, Bottom: 8}, Spacing: 4},
				Children: []Widget{
					Label{AssignTo: &delayLbl, Text: i18n.T("param.delay")},
					NumberEdit{AssignTo: &delayEdit, ToolTipText: i18n.T("tt.delay"), MinSize: Size{Width: 200, Height: 24}, MaxSize: Size{Width: 200, Height: 24}},
					Label{AssignTo: &startupLbl, Text: i18n.T("section.startup"), Font: Font{Family: "Microsoft YaHei UI", PointSize: 9, Bold: true}},
					CheckBox{AssignTo: &autoBootCb, Text: i18n.T("cb.autoboot"), ToolTipText: i18n.T("tt.autoboot"), Alignment: AlignHNearVNear},
					CheckBox{AssignTo: &autoGuardCb, Text: i18n.T("cb.autoguard"), ToolTipText: i18n.T("tt.autoguard"), Alignment: AlignHNearVNear},
					Label{AssignTo: &windowLbl, Text: i18n.T("section.window"), Font: Font{Family: "Microsoft YaHei UI", PointSize: 9, Bold: true}},
					CheckBox{AssignTo: &minTrayCb, Text: i18n.T("cb.mintotray"), ToolTipText: i18n.T("tt.mintotray"), Alignment: AlignHNearVNear},
					CheckBox{AssignTo: &closeTrayCb, Text: i18n.T("cb.closetotray"), ToolTipText: i18n.T("tt.closetotray"), Alignment: AlignHNearVNear},
				},
			},
			// Log
			GroupBox{
				AssignTo: &logGroup,
				Title:    i18n.T("log.group"),
				Layout:   VBox{Margins: Margins{Left: 8, Top: 6, Right: 8, Bottom: 8}, Spacing: 4},
				Children: []Widget{
					TextEdit{
						AssignTo:      &logBox,
						ReadOnly:      true,
						VScroll:       true,
						StretchFactor: 1,
						MinSize:       Size{Width: 0, Height: 160},
					},
				},
			},
		},
	}.Create()

	if runErr != nil {
		Notify(i18n.T("ui.title"), i18n.F("notify.createWinFail", runErr.Error()))
		return
	}

	// Parameter defaults: from persisted config (or defaults if none)
	delayEdit.SetValue(float64(pc.DelayMS))
	delayEdit.SetIncrement(100)
	delayEdit.SetDecimals(0)

	// Checkbox initial values (run-at-startup reads the registry: system-level,
	// not in config.json)
	autoBootCb.SetChecked(registry.Enabled())
	autoGuardCb.SetChecked(pc.AutoStartGuard)
	minTrayCb.SetChecked(pc.MinToTray)
	closeTrayCb.SetChecked(pc.CloseToTray)
	updateLangMenu()

	// Toggling run-at-startup writes to the registry immediately
	autoBootCb.CheckedChanged().Attach(func() {
		if err := registry.Set(autoBootCb.Checked()); err != nil {
			log.Error(i18n.F("log.autobootFail", err))
			return
		}
		if autoBootCb.Checked() {
			log.Info(i18n.T("log.autobootOn"))
		} else {
			log.Info(i18n.T("log.autobootOff"))
		}
	})

	// Auto-save on parameter change (debounced 250ms): persisted after Enter,
	// clicking away, or switching windows
	delayEdit.ValueChanged().Attach(scheduleConfigSave)
	mainWin.Deactivating().Attach(func() { applyUiConfig() })

	updateStatus("stopped")
	updateStartStopBtn()

	if err := createTray(); err != nil {
		log.Error(i18n.F("log.trayFail", err))
	}

	// Minimize → hide to tray (controlled by the "Minimize to tray" checkbox)
	mainWin.SizeChanged().Attach(func() {
		if win.IsIconic(mainWin.Handle()) && minTrayCb.Checked() {
			minimizeToTray()
		}
	})

	// Close window: hide to tray if "Close to tray" is checked, otherwise exit
	mainWin.Closing().Attach(func(canceled *bool, reason walk.CloseReason) {
		if quitting {
			applyUiConfig()
			shutdown()
			return
		}
		applyUiConfig()
		if closeTrayCb.Checked() {
			*canceled = true
			minimizeToTray()
			return
		}
		shutdown()
	})

	log.Print(i18n.F("log.started", appVersion))

	// Auto-guard after login
	if pc.AutoStartGuard {
		log.Print(i18n.T("log.autoguard"))
		startEngine()
		minimizeToTray()
	}

	mainWin.Run()
}

// shutdown releases resources on exit (tray / power notifications).
func shutdown() {
	if tray != nil {
		tray.Dispose()
		tray = nil
	}
}

// quitApp exits explicitly (tray/menu): saves config, releases resources, then
// terminates the process. It does not rely on walk close callbacks, avoiding the
// re-entrancy issues of closing the window from a tray-menu context.
func quitApp() {
	quitting = true
	applyUiConfig()
	shutdown()
	os.Exit(0)
}

// setLang switches the UI language (persisted to config.json; applied live).
func setLang(code string) {
	if i18n.GetLang() == code {
		return
	}
	i18n.SetLang(code)
	p := config.Load()
	p.Lang = code
	if err := config.Save(p); err != nil {
		log.Error(i18n.F("log.saveConfigFail", err))
		return
	}
	savedConfig.Lang = code
	retranslateUI()
	log.Print(i18n.F("log.langChanged", i18n.LangName(code)))
}

// retranslateUI re-applies the current-language strings to every text-bearing
// widget (menus, groups, labels, buttons, tooltips, tray), so switching the
// language takes effect immediately without a restart.
func retranslateUI() {
	if mainWin != nil {
		mainWin.SetTitle(i18n.T("ui.title"))
	}
	setText := func(a *walk.Action, key string) {
		if a != nil {
			a.SetText(i18n.T(key))
		}
	}
	setText(fileMenuAct, "menu.file")
	setText(exitAct, "menu.exit")
	setText(toolsMenuAct, "menu.tools")
	setText(testSleepAct, "menu.testSleep")
	setText(langMenuAct, "menu.language")
	setText(helpMenuAct, "menu.help")
	setText(checkUpdAct, "menu.checkUpdate")
	setText(githubAct, "menu.github")
	setText(aboutAct, "menu.about")
	setText(trayOpenAct, "menu.showWindow")
	setText(trayExitAct, "menu.exit")

	setGroup := func(g *walk.GroupBox, key string) {
		if g != nil {
			g.SetTitle(i18n.T(key))
		}
	}
	setGroup(statusGroup, "status.group")
	setGroup(configGroup, "config.group")
	setGroup(logGroup, "log.group")

	if delayLbl != nil {
		delayLbl.SetText(i18n.T("param.delay"))
	}
	if startupLbl != nil {
		startupLbl.SetText(i18n.T("section.startup"))
	}
	if windowLbl != nil {
		windowLbl.SetText(i18n.T("section.window"))
	}

	if delayEdit != nil {
		delayEdit.SetToolTipText(i18n.T("tt.delay"))
	}
	setCb := func(cb *walk.CheckBox, key, ttKey string) {
		if cb != nil {
			cb.SetText(i18n.T(key))
			cb.SetToolTipText(i18n.T(ttKey))
		}
	}
	setCb(autoBootCb, "cb.autoboot", "tt.autoboot")
	setCb(autoGuardCb, "cb.autoguard", "tt.autoguard")
	setCb(minTrayCb, "cb.mintotray", "tt.mintotray")
	setCb(closeTrayCb, "cb.closetotray", "tt.closetotray")

	if tray != nil {
		tray.SetToolTip(i18n.T("tray.tooltip"))
	}

	updateStatus("")
	updateStartStopBtn()
	updateLangMenu()
}

// langMenuItems builds one checkable menu action per supported language.
func langMenuItems() []MenuItem {
	codes := i18n.Codes()
	langActs = make([]*walk.Action, len(codes))
	items := make([]MenuItem, 0, len(codes))
	for i, code := range codes {
		i, code := i, code
		langActs[i] = &walk.Action{}
		items = append(items, Action{
			AssignTo:    &langActs[i],
			Text:        i18n.LangName(code),
			Checkable:   true,
			Checked:     i18n.GetLang() == code,
			OnTriggered: func() { setLang(code) },
		})
	}
	return items
}

// updateLangMenu checks the menu item matching the current language.
func updateLangMenu() {
	cur := i18n.GetLang()
	for _, act := range langActs {
		if act != nil {
			act.SetChecked(act.Text() == i18n.LangName(cur))
		}
	}
}

// toggleEngine is the combined start/stop button.
func toggleEngine() {
	if engine.Running() {
		stopEngine()
	} else {
		startEngine()
	}
}

var saveTimer *time.Timer

// scheduleConfigSave debounced save: persists after 250ms with no further
// change.
func scheduleConfigSave() {
	uiMu.Lock()
	w := mainWin
	uiMu.Unlock()
	if w == nil {
		return
	}
	if saveTimer != nil {
		saveTimer.Stop()
	}
	saveTimer = time.AfterFunc(250*time.Millisecond, func() {
		w.Synchronize(func() { applyUiConfig() })
	})
}

// applyUiConfig reads the controls into the engine params and persists them.
func applyUiConfig() {
	pc := config.Default()
	if delayEdit != nil {
		pc.DelayMS = int(delayEdit.Value())
	}
	if autoGuardCb != nil {
		pc.AutoStartGuard = autoGuardCb.Checked()
	}
	if minTrayCb != nil {
		pc.MinToTray = minTrayCb.Checked()
	}
	if closeTrayCb != nil {
		pc.CloseToTray = closeTrayCb.Checked()
	}
	pc.Lang = i18n.GetLang()

	// Enforce only a safety floor; no upper limit.
	if pc.DelayMS < 100 {
		pc.DelayMS = 100
	}
	engine.SetDelay(pc.DelayMS)

	// Persist only on actual change; applyUiConfig fires on every window
	// deactivation/close, and rewriting config.json each time would be wasted
	// disk I/O.
	if pc == savedConfig {
		return
	}
	if err := config.Save(pc); err != nil {
		log.Error(i18n.F("log.saveConfigFail", err))
		return
	}
	savedConfig = pc
}

func updateStatus(s string) {
	if statusLbl == nil {
		return
	}
	if s == "" {
		s = "stopped"
		if engine.Running() {
			s = "running"
		}
	}
	if s == "running" {
		statusLbl.SetText(i18n.T("status.guarding"))
		statusLbl.SetTextColor(statusRun)
	} else {
		statusLbl.SetText(i18n.T("status.stopped"))
		statusLbl.SetTextColor(statusStopRed)
	}
}

func updateStartStopBtn() {
	if startStopBtn == nil {
		return
	}
	if engine.Running() {
		startStopBtn.SetText(i18n.T("btn.stopGuard"))
	} else {
		startStopBtn.SetText(i18n.T("btn.startGuard"))
	}
}

// startEngine starts the engine; locks the params after starting.
func startEngine() {
	if engine.Running() {
		log.Print(i18n.T("log.engineAlready"))
		return
	}
	applyUiConfig()
	engine.Start()
	delayEdit.SetEnabled(false)
	log.Info(i18n.F("log.engineStarted", engine.Delay()))
	updateStatus("running")
	updateStartStopBtn()
}

// stopEngine stops the engine; returns whether it was actually running.
func stopEngine() bool {
	if !engine.Stop() {
		return false
	}
	delayEdit.SetEnabled(true)
	log.Info(i18n.T("log.engineStopped"))
	updateStatus("stopped")
	updateStartStopBtn()
	return true
}
