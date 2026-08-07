# LidReSleep

[English](./README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Français](./README.fr.md) | [Deutsch](./README.de.md) | **Español** | [Русский](./README.ru.md) | [Português](./README.pt.md) | [Italiano](./README.it.md)

Una pequeña herramienta de Windows en segundo plano que evita que tu portátil se caliente durante la noche después de cerrar la tapa.

Muchos portátiles modernos (Modern Standby) no duermen de verdad al cerrar la tapa: la pantalla se apaga pero el sistema sigue conectado en bajo consumo, y puede ser fácilmente **despertado y mantenido despierto** por solicitudes de red, tareas en segundo plano, etc., causando calor y consumo de batería toda la noche.

Lo que hace LidReSleep es simple: **duerme al cerrar la tapa; si se despierta de forma inesperada con la tapa aún cerrada, vuelve a dormirse automáticamente hasta que abras la tapa.**

- Aplicación portátil de un solo archivo, sin instalación
- 10 idiomas de interfaz (简体中文 / English / 日本語 / 한국어 / Français / Deutsch / Español / Русский / Português / Italiano), detectados automáticamente del idioma del sistema
- Windows 10/11 (x64 / ARM64 / x86), elige la versión que coincida con tu CPU

![Captura de pantalla de LidReSleep](screenshot.jpg)

## Inicio rápido

1. Descarga la versión para tu CPU y haz doble clic en ella (p. ej. `LidReSleep-amd64.exe`) para abrir el panel.
2. Haz clic en **「▶ Iniciar protección」**; el estado cambia a `● Protegiendo`.
3. Cierra la tapa y listo.

Después: cerrar la tapa → dormir de inmediato; despertar con la tapa aún cerrada → volver a dormir tras ~3 segundos (predeterminado); abrir la tapa → cancelar y continuar con normalidad.

## Guía de la interfaz

### Estado
- `● Detenido / ● Protegiendo`, indica si la protección está activa.
- Botón `▶ Iniciar protección / ■ Detener protección`, el texto cambia según el estado.

### Configuración

| Ajuste | Descripción | Predeterminado |
|---|---|---|
| Retardo de suspensión (ms) | esperar este tiempo antes de volver a dormir tras despertar | `3000` |
| ☑ Ejecutar al inicio | ejecutarse automáticamente al iniciar Windows (nivel de sistema, registro) | No |
| ☑ Protección automática al iniciar sesión | iniciar la protección y minimizar a la bandeja al abrir | No |
| ☑ Minimizar a la bandeja | ocultar en la bandeja al minimizar | Sí |
| ☑ Cerrar en la bandeja | ocultar en la bandeja en lugar de salir al cerrar | Sí |

- Los cambios se **guardan automáticamente**, no hace falta guardar manualmente.
- Ejecutar al inicio se aplica de inmediato (registro); el resto se guarda en `config.json` junto al exe.

### Menús
- **Archivo**: Salir
- **Herramientas**: Probar suspensión (dormir una vez para verificar)
- **Idioma**: 10 idiomas (la marca muestra el actual; se aplica de inmediato)
- **Ayuda**: Buscar actualizaciones (comprueba GitHub), Página del proyecto, Acerca de (funciones y explicación de Modern Standby)

### Registro
Cada línea tiene una marca de tiempo y un nivel, mostrando en tiempo real los eventos de tapa/despertar/re-suspensión, con desplazamiento automático y límite de 200 KB.

| Nivel | Significado |
|---|---|
| `INFO` | información general |
| `EVENT` | eventos del sistema (tapa/suspensión/despertar) |
| `ACTION` | acciones del programa (programar/cancelar/dormir) |
| `ERROR` | error (con motivo) |

### Bandeja del sistema
- Se oculta en la bandeja al minimizar/cerrar (si está activado).
- Clic izquierdo en el icono: restaurar la ventana principal; clic derecho: Mostrar ventana principal / Salir.

## Preguntas frecuentes

**¿Por qué mi portátil sigue calentándose después de cerrar la tapa?**
Probablemente Modern Standby lo despierta. Usa Herramientas → Probar suspensión para verificar; tras cerrar la tapa con la protección activa, busca `ACTION Ejecutando suspensión` en el registro.

**¿Qué es Modern Standby?**
Consulta la explicación sencilla en Ayuda → Acerca de.

**¿Cómo salgo completamente?**
Clic derecho en el icono de la bandeja → Salir (si "Cerrar en la bandeja" está activado, ✕ solo minimiza).

**¿El inicio automático no funciona?**
Depende de la cuenta de usuario actual; los fallos (p. ej. al ejecutarse desde una ubicación restringida) se registran con el motivo.

## Usuarios avanzados: reducir los despertamientos (opcional)

La herramienta se encarga de "volver a dormir tras despertar". Para reducir los despertamientos, ejecuta en un PowerShell elevado:

```powershell
# deshabilitar temporizadores de activación
powercfg /setacvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
powercfg /setdcvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
# inspeccionar fuentes de activación
powercfg /waketimers
powercfg /devicequery wake_armed
# restaurar valores predeterminados
powercfg /restoredefaultschemes
```

---

## Descarga

Obtén la última versión desde GitHub Releases: [LidReSleep Releases](https://github.com/dootn/LidReSleep/releases/latest)

| Archivo | CPU | Descarga |
|---|---|---|
| `LidReSleep-amd64.exe` | Intel/AMD 64 bits (la mayoría de PC) | [amd64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-amd64.exe) |
| `LidReSleep-arm64.exe` | ARM64 (p. ej. Surface Pro X, PC Snapdragon) | [arm64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-arm64.exe) |
| `LidReSleep-386.exe` | x86 de 32 bits | [386](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-386.exe) |

> Página del proyecto: https://github.com/dootn/LidReSleep
