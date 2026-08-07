# LidReSleep

[English](./README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | **한국어** | [Français](./README.fr.md) | [Deutsch](./README.de.md) | [Español](./README.es.md) | [Русский](./README.ru.md) | [Português](./README.pt.md) | [Italiano](./README.it.md)

노트북을 닫은 후 밤새 발열되는 문제를 해결하는 작은 Windows 백그라운드 도구입니다.

많은 최신 노트북(Modern Standby)은 덮개를 닫아도 실제로 절전되지 않고, 화면만 꺼진 채 시스템이 저전력으로 연결 상태를 유지합니다. 네트워크 요청, 백그라운드 작업 등으로 **쉽게 깨어나 깨어난 상태가 유지**되어 밤새 발열과 배터리 소모가 발생할 수 있습니다.

LidReSleep이 하는 일은 간단합니다. **덮개를 닫으면 절전합니다. 덮개가 닫힌 상태에서 예기치 않게 깨어나면, 덮개를 열 때까지 자동으로 다시 절전합니다.**

- 설치가 필요 없는 휴대용 단일 파일 앱
- 10개 언어 UI(简体中文 / English / 日本語 / 한국어 / Français / Deutsch / Español / Русский / Português / Italiano), 시스템 언어에서 자동 감지
- Windows 10/11(x64 / ARM64 / x86), CPU에 맞는 빌드 선택

![LidReSleep 스크린샷](screenshot.jpg)

## 빠른 시작

1. CPU에 맞는 빌드를 다운로드한 후 더블 클릭합니다(예: `LidReSleep-amd64.exe`).
2. **「▶ 보호 시작」** 을 클릭하면 상태가 `● 보호 중`이 됩니다.
3. 덮개를 닫고 가면 됩니다.

이후: 덮개 닫기 → 즉시 절전. 덮개가 닫힌 상태로 깨어남 → 기본 약 3초 후 다시 절전. 덮개 열기 → 취소하고 정상 동작.

## UI 가이드

### 상태
- `● 중지됨 / ● 보호 중`, 보호가 실행 중인지 표시합니다.
- `▶ 보호 시작 / ■ 보호 중지` 버튼, 상태에 따라 텍스트가 바뀝니다.

### 설정

| 설정 | 설명 | 기본값 |
|---|---|---|
| 절전 지연(ms) | 깨어난 후 이 시간을 기다렸다가 절전합니다 | `3000` |
| ☑ 시작 시 실행 | Windows 로그인 시 자동 실행(시스템 수준, 레지스트리) | 아니요 |
| ☑ 로그인 후 자동 보호 | 시작 시 보호를 시작하고 트레이로 최소화 | 아니요 |
| ☑ 최소화 시 트레이로 | 최소화할 때 트레이로 숨김 | 예 |
| ☑ 닫을 때 트레이로 | 닫을 때 종료하지 않고 트레이로 숨김 | 예 |

- 변경 사항은 **자동 저장**됩니다. 수동 저장이 필요 없습니다.
- 시작 시 실행은 즉시 적용됩니다(레지스트리). 나머지 설정은 exe 옆의 `config.json`에 저장됩니다.

### 메뉴
- **파일**: 종료
- **도구**: 절전 테스트(확인을 위해 한 번 절전)
- **언어**: 10개 언어(현재 언어에 체크 표시, 즉시 적용)
- **도움말**: 업데이트 확인(GitHub에서 새 버전 확인), 프로젝트 페이지, 정보(기능 및 Modern Standby 설명)

### 로그
각 줄에 타임스탬프와 레벨 태그가 있으며, 덮개/깨어남/재절전 이벤트를 실시간으로 표시합니다. 자동 스크롤, 최대 200KB로 제한됩니다.

| 레벨 | 의미 |
|---|---|
| `INFO` | 일반 정보 |
| `EVENT` | 시스템 이벤트(덮개/절전/깨어남) |
| `ACTION` | 프로그램 동작(예약/취소/절전) |
| `ERROR` | 오류(이유 포함) |

### 시스템 트레이
- 최소화 또는 닫을 때(사용 시) 트레이로 숨겨집니다.
- 아이콘 왼쪽 클릭: 기본 창 복원. 오른쪽 클릭: 기본 창 표시 / 종료.

## FAQ

**덮개를 닫았는데도 노트북이 여전히 뜨거운 이유는?**
Modern Standby가 깨어나게 하는 경우가 많습니다. 도구 → 절전 테스트로 확인하세요. 보호 중 덮개를 닫은 후 로그에 `ACTION 절전 실행`이 표시되는지 확인합니다.

**Modern Standby란 무엇인가요?**
도움말 → 정보의 쉬운 설명을 참고하세요.

**완전히 종료하려면?**
트레이 아이콘 오른쪽 클릭 → 종료("닫을 때 트레이로"가 켜져 있으면 ✕는 최소화만 합니다).

**시작 시 실행이 안 되나요?**
현재 사용자 계정에 따라 다릅니다. 실패(예: 제한된 위치에서 실행)하면 이유와 함께 로그에 기록됩니다.

## 고급 사용자: 깨어남 줄이기(선택 사항)

이 도구는 "깨어난 후 다시 절전"을 처리합니다. 깨어남 자체를 줄이려면 관리자 권한 PowerShell에서 실행하세요.

```powershell
# 깨우기 타이머 비활성화
powercfg /setacvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
powercfg /setdcvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
# 깨우기 소스 확인
powercfg /waketimers
powercfg /devicequery wake_armed
# 기본값 복원
powercfg /restoredefaultschemes
```

---

## 다운로드

최신 버전은 GitHub Releases에서: [LidReSleep Releases](https://github.com/dootn/LidReSleep/releases/latest)

| 파일 | CPU | 다운로드 |
|---|---|---|
| `LidReSleep-amd64.exe` | Intel/AMD 64비트(대부분의 PC) | [amd64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-amd64.exe) |
| `LidReSleep-arm64.exe` | ARM64(예: Surface Pro X, Snapdragon PC) | [arm64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-arm64.exe) |
| `LidReSleep-386.exe` | 32비트 x86 | [386](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-386.exe) |

> 프로젝트 페이지: https://github.com/dootn/LidReSleep
