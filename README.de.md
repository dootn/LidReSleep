# LidReSleep

[English](./README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Français](./README.fr.md) | **Deutsch** | [Español](./README.es.md) | [Русский](./README.ru.md) | [Português](./README.pt.md) | [Italiano](./README.it.md)

Ein kleines Windows-Hintergrundtool, das verhindert, dass Ihr Laptop nach dem Schließen des Deckels über Nacht heiß wird.

Viele moderne Laptops (Modern Standby) schlafen nicht wirklich, wenn der Deckel geschlossen wird – der Bildschirm geht nur aus, während das System mit geringem Stromverbrauch verbunden bleibt und leicht durch Netzwerkanfragen, Hintergrundaufgaben usw. **aufgeweckt und wachgehalten werden** kann, was zu Hitze und Akkuverbrauch über Nacht führt.

Was LidReSleep tut, ist einfach: **Schlafen, wenn der Deckel geschlossen wird; falls er unerwartet geweckt wird, während der Deckel noch geschlossen ist, automatisch wieder schlafen, bis Sie den Deckel öffnen.**

- Portables Einzeldatei-Programm, keine Installation nötig
- 10 UI-Sprachen (简体中文 / English / 日本語 / 한국어 / Français / Deutsch / Español / Русский / Português / Italiano), automatisch aus der Systemsprache erkannt
- Windows 10/11 (x64 / ARM64 / x86), wählen Sie den Build passend zu Ihrer CPU

![LidReSleep-Screenshot](screenshot.jpg)

## Schnellstart

1. Laden Sie den Build für Ihre CPU herunter und doppelklicken Sie darauf (z. B. `LidReSleep-amd64.exe`), um das Panel zu öffnen.
2. Klicken Sie auf **「▶ Wache starten」**; der Status wechselt zu `● Wache aktiv`.
3. Deckel schließen und fertig.

Danach: Deckel schließen → sofort schlafen; mit geschlossenem Deckel geweckt → nach ~3 Sekunden wieder schlafen (Standard); Deckel öffnen → abbrechen und normal weiterarbeiten.

## Bedienungsanleitung

### Status
- `● Gestoppt / ● Wache aktiv`, zeigt, ob die Wache läuft.
- Button `▶ Wache starten / ■ Wache stoppen`, der Text wechselt je nach Zustand.

### Einstellungen

| Einstellung | Beschreibung | Standard |
|---|---|---|
| Ruheverzögerung (ms) | nach dem Wecken diese Zeit warten, bevor erneut geschlafen wird | `3000` |
| ☑ Beim Start ausführen | automatisch bei Windows-Anmeldung starten (Systemebene, Registrierung) | Nein |
| ☑ Nach Anmeldung automatisch bewachen | beim Start bewachen und in das Infobereich minimieren | Nein |
| ☑ Beim Minimieren in das Infobereich | beim Minimieren im Infobereich ausblenden | Ja |
| ☑ Beim Schließen in das Infobereich | beim Schließen statt Beenden im Infobereich ausblenden | Ja |

- Änderungen werden **automatisch gespeichert**, keine manuelle Speicherung nötig.
- Autostart gilt sofort (Registrierung); andere Einstellungen werden in `config.json` neben der exe gespeichert.

### Menüs
- **Datei**: Beenden
- **Extras**: Ruhezustand testen (einmal schlafen, um zu prüfen)
- **Sprache**: 10 Sprachen (Häkchen zeigt die aktuelle; gilt sofort)
- **Hilfe**: Nach Updates suchen (prüft GitHub auf eine neue Version), Projektseite, Über (Funktionen und Modern-Standby-Erklärung)

### Protokoll
Jede Zeile hat einen Zeitstempel und eine Stufe, zeigt Deckel-/Aufweck-/Wiederschlaf-Ereignisse in Echtzeit, mit automatischem Scrollen, begrenzt auf 200 KB.

| Stufe | Bedeutung |
|---|---|
| `INFO` | allgemeine Informationen |
| `EVENT` | Systemereignisse (Deckel/Schlaf/Aufwecken) |
| `ACTION` | Programmaktionen (Planung/Abbrechen/Schlaf) |
| `ERROR` | Fehler (mit Grund) |

### Infobereich
- Beim Minimieren/Schließen (falls aktiviert) im Infobereich verbergen.
- Symbol linksklicken: Hauptfenster wiederherstellen; rechtsklicken: Hauptfenster anzeigen / Beenden.

## FAQ

**Warum wird mein Laptop nach dem Schließen des Deckels trotzdem heiß?**
Wahrscheinlich weckt Modern Standby ihn auf. Verwenden Sie Extras → Ruhezustand testen zur Überprüfung; nach dem Schließen des Deckels während der Wache im Protokoll nach `ACTION Schlaf ausführen` suchen.

**Was ist Modern Standby?**
Siehe die verständliche Erklärung unter Hilfe → Über.

**Wie beende ich das Programm vollständig?**
Symbol im Infobereich rechtsklicken → Beenden (wenn „Beim Schließen in das Infobereich" aktiviert ist, minimiert ✕ nur).

**Autostart funktioniert nicht?**
Das hängt vom aktuellen Benutzerkonto ab; Fehler (z. B. beim Start von einem eingeschränkten Ort) werden mit einem Grund protokolliert.

## Power-User: Aufwecken reduzieren (optional)

Das Tool kümmert sich um „nach dem Aufwecken wieder schlafen". Um die Aufweckvorgänge selbst zu reduzieren, führen Sie in einer erweiterten PowerShell aus:

```powershell
# Wecktimer deaktivieren
powercfg /setacvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
powercfg /setdcvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
# Aufweckquellen prüfen
powercfg /waketimers
powercfg /devicequery wake_armed
# Standardwerte wiederherstellen
powercfg /restoredefaultschemes
```

---

## Download

Die neueste Version von GitHub Releases: [LidReSleep Releases](https://github.com/dootn/LidReSleep/releases/latest)

| Datei | CPU | Download |
|---|---|---|
| `LidReSleep-amd64.exe` | Intel/AMD 64-Bit (die meisten PCs) | [amd64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-amd64.exe) |
| `LidReSleep-arm64.exe` | ARM64 (z. B. Surface Pro X, Snapdragon-PCs) | [arm64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-arm64.exe) |
| `LidReSleep-386.exe` | 32-Bit-x86 | [386](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-386.exe) |

> Projektseite: https://github.com/dootn/LidReSleep
