# LidReSleep

[English](./README.md) | **简体中文** | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Français](./README.fr.md) | [Deutsch](./README.de.md) | [Español](./README.es.md) | [Русский](./README.ru.md) | [Português](./README.pt.md) | [Italiano](./README.it.md)

一款小巧的 Windows 后台工具，解决笔记本"合上盖子还在发热"的问题。

很多现代笔记本（Modern Standby）合盖后并没有真正睡眠，只是屏幕熄了，系统仍在低功耗联网，很容易被网络请求、后台任务等**意外唤醒并保持唤醒**，导致整夜发热、耗电。

LidReSleep 做的事很简单：**合盖就睡眠；万一被意外唤醒而盖子仍关闭，就自动再次入睡，直到你开盖。**

- 单文件绿色软件，双击即用，无需安装
- 界面支持 10 种语言（简体中文 / English / 日本語 / 한국어 / Français / Deutsch / Español / Русский / Português / Italiano），自动按系统语言显示
- Windows 10/11（x64 / ARM64 / x86），按 CPU 选择对应版本

![LidReSleep 界面截图](screenshot.jpg)

## 快速开始

1. 按 CPU 选对应版本下载，双击 `LidReSleep-amd64.exe`（Intel/AMD 64 位），打开控制面板。
2. 点 **「▶ 启动守护」**，状态变为 `● 正在守护`。
3. 合上盖子走人即可。

之后：合盖 → 立即睡眠；被意外唤醒且盖子仍闭合 → 默认约 3 秒后再次入睡；开盖 → 取消，恢复正常。

## 界面说明

### 状态区
- `● 已停止 / ● 正在守护`，直观显示守护状态。
- `▶ 启动守护 / ■ 停止守护` 按钮，文字随状态切换。

### 配置区

| 配置项 | 说明 | 默认 |
|---|---|---|
| 入睡延迟(ms) | 被唤醒后等待这么久再入睡 | `3000` |
| ☑ 开机启动 | 随 Windows 登录自动运行（系统级，写注册表） | 否 |
| ☑ 登录后自动守护 | 启动后直接进入守护并最小化到托盘 | 否 |
| ☑ 最小化到托盘 | 点最小化按钮时隐藏到系统托盘 | 是 |
| ☑ 关闭窗口隐藏到托盘 | 点窗口 ✕ 时隐藏到托盘而非退出 | 是 |

- 参数改动**自动保存**（按回车或离开输入框后写入），无需手动保存。
- 开机启动为系统级设置，勾选后立即生效；其余配置记在 exe 同目录 `config.json`。

### 菜单栏
- **文件**：退出
- **工具**：测试睡眠（立即试睡一次验证）
- **语言**：10 种语言可选（当前语言带勾选；切换后立即生效）
- **帮助**：检查更新（从 GitHub 检查是否有新版本）、项目主页、关于（含功能与 Modern Standby 科普）

### 运行日志
每行含时间戳与等级标签，实时显示合盖/唤醒/再入睡全过程，自动滚动，最多保留 200KB。

| 等级 | 含义 |
|---|---|
| `INFO` | 一般信息 |
| `EVENT` | 系统事件（盖子开关/睡眠/唤醒） |
| `ACTION` | 程序动作（排程/取消/执行入睡） |
| `ERROR` | 出错（含原因） |

### 系统托盘
- 最小化或关窗（勾选后）即隐藏到托盘。
- 图标左键：恢复主窗口；右键：显示主窗口 / 退出。

## 常见问题

**为什么我合盖后电脑还在发热？**
多半是 Modern Standby 设备合盖后仍被唤醒。点「工具 → 测试睡眠」立即验证；守护中合盖后看日志是否出现 `ACTION 执行睡眠操作`。

**什么是 Modern Standby？**
见「帮助 → 关于」里的通俗说明。

**怎么彻底退出？**
托盘图标右键 → 退出（若勾选"关闭窗口隐藏到托盘"，点 ✕ 只是最小化）。

**勾选了开机启动但没生效？**
开机启动依赖当前用户账户；从 U 盘/受限目录运行时写入失败，会在日志提示原因。

## 进阶：从源头减少唤醒（可选）

工具负责"唤醒后再入睡"。想减少被唤醒次数，可在管理员 PowerShell 中执行：

```powershell
# 关闭"允许唤醒定时器"
powercfg /setacvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
powercfg /setdcvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
# 查看唤醒源
powercfg /waketimers
powercfg /devicequery wake_armed
# 恢复默认
powercfg /restoredefaultschemes
```

---

## 下载

前往 GitHub Releases 获取最新版本：[LidReSleep Releases](https://github.com/dootn/LidReSleep/releases/latest)

| 文件 | CPU | 下载 |
|---|---|---|
| `LidReSleep-amd64.exe` | Intel/AMD 64 位（主流） | [amd64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-amd64.exe) |
| `LidReSleep-arm64.exe` | ARM64（Surface Pro X、骁龙本等） | [arm64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-arm64.exe) |
| `LidReSleep-386.exe` | 32 位 x86 | [386](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-386.exe) |

> 项目主页：https://github.com/dootn/LidReSleep
