# LidReSleep

[English](./README.md) | [简体中文](./README.zh-CN.md) | **日本語** | [한국어](./README.ko.md) | [Français](./README.fr.md) | [Deutsch](./README.de.md) | [Español](./README.es.md) | [Русский](./README.ru.md) | [Português](./README.pt.md) | [Italiano](./README.it.md)

ノートパソコンの蓋を閉じた後の夜間の発熱を防ぐ、小さな Windows バックグラウンドツールです。

多くの現代のノートパソコン（Modern Standby）は、蓋を閉じても本当にはスリープせず、画面が消えるだけでシステムは低消費電力で接続されたままになり、ネットワークリクエストやバックグラウンドタスクなどによって**簡単に起動され、起動されたまま維持される**ことがあります。その結果、一晩中発熱やバッテリーの消耗が続きます。

LidReSleep が行うことはシンプルです。**蓋を閉じたらスリープします。万一蓋が閉じたまま予期せず起動された場合は、蓋を開けるまで自動的に再びスリープします。**

- ポータブルな単一ファイルのアプリ、インストール不要
- 10 言語の UI（简体中文 / English / 日本語 / 한국어 / Français / Deutsch / Español / Русский / Português / Italiano）、システム言語から自動検出
- Windows 10/11（x64 / ARM64 / x86）、CPU に合ったビルドを選択

![LidReSleep スクリーンショット](screenshot.jpg)

## クイックスタート

1. 自分の CPU に合ったビルドをダウンロードし、ダブルクリックします（例：`LidReSleep-amd64.exe`）。
2. **「▶ 監視を開始」** をクリックすると、状態が `● 監視中` になります。
3. 蓋を閉じてそのままにします。

その後は：蓋を閉じる → すぐにスリープ。蓋が閉じたまま起動される → 既定で約 3 秒後に再びスリープ。蓋を開ける → キャンセルして通常に戻ります。

## UI ガイド

### 状態
- `● 停止中 / ● 監視中` で、監視が動作しているかを表示します。
- `▶ 監視を開始 / ■ 監視を停止` ボタン、状態に応じて文字が切り替わります。

### 設定

| 設定 | 説明 | 既定 |
|---|---|---|
| スリープ遅延(ms) | 起動された後にこの時間待ってからスリープします | `3000` |
| ☑ 起動時に実行 | Windows ログイン時に自動実行（システムレベル、レジストリ） | いいえ |
| ☑ ログイン後に自動で監視 | 起動時に監視を開始し、トレイに最小化 | いいえ |
| ☑ 最小化でトレイに格納 | 最小化時にトレイに隠す | はい |
| ☑ 閉じるボタンでトレイに格納 | 閉じる時に終了せずトレイに隠す | はい |

- 変更は**自動的に保存**されます。手動保存は不要です。
- 起動時に実行は即時に反映されます（レジストリ）。それ以外の設定は exe と同じ場所の `config.json` に保存されます。

### メニュー
- **ファイル**: 終了
- **ツール**: スリープをテスト（検証のため一度だけスリープ）
- **言語**: 10 言語（現在の言語にチェックが付きます。即時に適用）
- **ヘルプ**: 更新を確認（GitHub で新バージョンを確認）、プロジェクトページ、について（機能と Modern Standby の解説）

### ログ
各行にタイムスタンプとレベルタグがあり、蓋の開閉・起動・再スリープのイベントをリアルタイムに表示します。自動スクロールし、最大 200KB で切り詰められます。

| レベル | 意味 |
|---|---|
| `INFO` | 一般的な情報 |
| `EVENT` | システムイベント（蓋の開閉・スリープ・起動） |
| `ACTION` | プログラムの動作（スケジュール・キャンセル・スリープ） |
| `ERROR` | エラー（理由付き） |

### システムトレイ
- 最小化または閉じる時に（有効な場合）トレイに隠れます。
- アイコン左クリック：メインウィンドウを復元。右クリック：メインウィンドウを表示 / 終了。

## よくある質問

**蓋を閉じたのにノートパソコンがまだ熱いのはなぜ？**
Modern Standby が起動している可能性が高いです。「ツール → スリープをテスト」で確認できます。監視中に蓋を閉じて、ログに `ACTION スリープを実行` が表示されるか確認してください。

**Modern Standby とは？**
「ヘルプ → について」にある分かりやすい解説をご覧ください。

**完全に終了するには？**
トレイアイコンを右クリック → 終了（「閉じるボタンでトレイに格納」が有効な場合、✕ は最小化のみ）。

**起動時に実行が機能しない？**
現在のユーザーアカウントに依存します。失敗（制限された場所から実行している場合など）は理由とともにログに記録されます。

## 上級者向け：起動回数を減らす（任意）

このツールは「起動された後に再びスリープする」を処理します。起動自体を減らすには、管理者権限の PowerShell で実行してください。

```powershell
# 起動タイマーを無効化
powercfg /setacvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
powercfg /setdcvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
# 起動元を調査
powercfg /waketimers
powercfg /devicequery wake_armed
# 既定に戻す
powercfg /restoredefaultschemes
```

---

## ダウンロード

最新版は GitHub Releases から：[LidReSleep Releases](https://github.com/dootn/LidReSleep/releases/latest)

| ファイル | CPU | ダウンロード |
|---|---|---|
| `LidReSleep-amd64.exe` | Intel/AMD 64 ビット（ほとんどの PC） | [amd64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-amd64.exe) |
| `LidReSleep-arm64.exe` | ARM64（Surface Pro X、Snapdragon PC など） | [arm64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-arm64.exe) |
| `LidReSleep-386.exe` | 32 ビット x86 | [386](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-386.exe) |

> プロジェクトページ: https://github.com/dootn/LidReSleep
