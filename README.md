# LidReSleep

**English** | [简体中文](./README.zh-CN.md)

A tiny Windows background tool that stops your laptop from heating up overnight after you close the lid.

Many modern laptops (Modern Standby) do not truly sleep when the lid is closed — the screen just goes off while the system stays connected at low power, and can easily be **woken and kept awake** by network requests, background tasks, etc., causing all-night heat and battery drain.

What LidReSleep does is simple: **sleep when the lid closes; if it is woken unexpectedly while the lid is still closed, sleep again automatically until you open the lid.**

- Portable single-file app, no install needed
- Chinese & English UI, auto-detected from system language
- Windows 10/11 (x64 / ARM64 / x86), pick the build matching your CPU

![LidReSleep screenshot](screenshot.jpg)

## Quick Start

1. Download the build for your CPU, then double-click it (e.g. `LidReSleep-amd64.exe`) to open the panel.
2. Click **「▶ Start Guard」**; the status turns `● Guarding`.
3. Close the lid and go.

After that: close the lid → sleep immediately; woken with the lid still closed → sleep again after ~3 seconds (default); open the lid → cancel and resume normally.

## UI Guide

### Status
- `● Stopped / ● Guarding`, shows whether the guard is running.
- `▶ Start Guard / ■ Stop Guard` button, text toggles with state.

### Settings

| Setting | Description | Default |
|---|---|---|
| Sleep delay (ms) | wait this long before sleeping after being woken | `3000` |
| ☑ Run at startup | run automatically at Windows login (system-level, registry) | No |
| ☑ Auto-guard after login | start guarding and minimize to tray on launch | No |
| ☑ Minimize to tray | hide to tray when minimizing | Yes |
| ☑ Close window to tray | hide to tray instead of exiting when closing | Yes |

- Changes are **saved automatically**, no manual save needed.
- Run-at-startup applies immediately (registry); other settings are stored in `config.json` next to the exe.

### Menus
- **File**: Exit
- **Tools**: Test Sleep (sleep once to verify)
- **Language**: 中文 / English (check shows current; takes effect after restart)
- **Help**: Check for Updates (checks GitHub for a new version), Project Page, About (features & Modern Standby explainer)

### Log
Each line has a timestamp and level tag, showing lid/wake/re-sleep events in real time, auto-scrolling, capped at 200KB.

| Level | Meaning |
|---|---|
| `INFO` | general info |
| `EVENT` | system events (lid/sleep/wake) |
| `ACTION` | program actions (schedule/cancel/sleep) |
| `ERROR` | error (with reason) |

### System Tray
- Hiding to tray on minimize/close (when enabled).
- Icon left-click: restore the main window; right-click: Show Main Window / Exit.

## FAQ

**Why does my laptop still get hot after closing the lid?**
Likely Modern Standby waking it. Use Tools → Test Sleep to verify; after closing the lid while guarding, check the log for `ACTION Executing sleep`.

**What is Modern Standby?**
See the plain-language explanation in Help → About.

**How do I exit completely?**
Right-click the tray icon → Exit (if "Close to tray" is enabled, ✕ only minimizes).

**Run at startup not working?**
It depends on the current user account; failures (e.g. running from a restricted location) are logged with a reason.

## Power users: reduce wake-ups (optional)

The tool handles "sleep again after waking". To reduce wake-ups themselves, run in an elevated PowerShell:

```powershell
# disable wake timers
powercfg /setacvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
powercfg /setdcvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
# inspect wake sources
powercfg /waketimers
powercfg /devicequery wake_armed
# restore defaults
powercfg /restoredefaultschemes
```

---

## Download

Get the latest version from GitHub Releases: [LidReSleep Releases](https://github.com/dootn/LidReSleep/releases/latest)

| File | CPU | Download |
|---|---|---|
| `LidReSleep-amd64.exe` | Intel/AMD 64-bit (most PCs) | [amd64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-amd64.exe) |
| `LidReSleep-arm64.exe` | ARM64 (e.g. Surface Pro X, Snapdragon PCs) | [arm64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-arm64.exe) |
| `LidReSleep-386.exe` | 32-bit x86 | [386](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-386.exe) |

> Project page: https://github.com/dootn/LidReSleep
