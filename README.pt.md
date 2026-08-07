# LidReSleep

[English](./README.md) | [简体中文](./README.zh-CN.md) | [日本語](./README.ja.md) | [한국어](./README.ko.md) | [Français](./README.fr.md) | [Deutsch](./README.de.md) | [Español](./README.es.md) | [Русский](./README.ru.md) | **Português** | [Italiano](./README.it.md)

Uma pequena ferramenta de Windows em segundo plano que impede que seu laptop aqueça durante a noite depois de fechar a tampa.

Muitos laptops modernos (Modern Standby) não entram de fato em suspensão quando a tampa é fechada — a tela apenas desliga enquanto o sistema continua conectado em baixo consumo, podendo ser facilmente **ativado e mantido ativo** por solicitações de rede, tarefas em segundo plano etc., causando aquecimento e consumo de bateria a noite toda.

O que o LidReSleep faz é simples: **suspender ao fechar a tampa; se for ativado inesperadamente com a tampa ainda fechada, voltar a suspender automaticamente até você abrir a tampa.**

- Aplicativo portátil de arquivo único, sem instalação
- 10 idiomas de interface (简体中文 / English / 日本語 / 한국어 / Français / Deutsch / Español / Русский / Português / Italiano), detectados automaticamente do idioma do sistema
- Windows 10/11 (x64 / ARM64 / x86), escolha a versão correspondente à sua CPU

![Captura de tela do LidReSleep](screenshot.jpg)

## Início rápido

1. Baixe a versão para sua CPU e dê um clique duplo (ex.: `LidReSleep-amd64.exe`) para abrir o painel.
2. Clique em **「▶ Iniciar proteção」**; o status muda para `● Protegendo`.
3. Feche a tampa e pronto.

Depois: fechar a tampa → suspender imediatamente; ser ativado com a tampa ainda fechada → suspender novamente após ~3 segundos (padrão); abrir a tampa → cancelar e retomar normalmente.

## Guia da interface

### Status
- `● Parado / ● Protegendo`, mostra se a proteção está em execução.
- Botão `▶ Iniciar proteção / ■ Parar proteção`, o texto alterna conforme o estado.

### Configurações

| Configuração | Descrição | Padrão |
|---|---|---|
| Atraso de suspensão (ms) | aguardar este tempo antes de suspender novamente após ser ativado | `3000` |
| ☑ Executar na inicialização | executar automaticamente no login do Windows (nível de sistema, registro) | Não |
| ☑ Proteção automática após login | iniciar a proteção e minimizar para a bandeja na abertura | Não |
| ☑ Minimizar para a bandeja | ocultar na bandeja ao minimizar | Sim |
| ☑ Fechar para a bandeja | ocultar na bandeja em vez de sair ao fechar | Sim |

- As alterações são **salvas automaticamente**, sem necessidade de salvar manualmente.
- Executar na inicialização é aplicado imediatamente (registro); os demais ajustes ficam em `config.json` ao lado do exe.

### Menus
- **Arquivo**: Sair
- **Ferramentas**: Testar suspensão (suspender uma vez para verificar)
- **Idioma**: 10 idiomas (a marca mostra o atual; aplica imediatamente)
- **Ajuda**: Verificar atualizações (verifica o GitHub), Página do projeto, Sobre (recursos e explicação do Modern Standby)

### Registro
Cada linha tem um carimbo de data/hora e um nível, mostrando em tempo real os eventos de tampa/ativação/re-suspensão, com rolagem automática e limite de 200 KB.

| Nível | Significado |
|---|---|
| `INFO` | informações gerais |
| `EVENT` | eventos do sistema (tampa/suspensão/ativação) |
| `ACTION` | ações do programa (agendar/cancelar/suspender) |
| `ERROR` | erro (com motivo) |

### Bandeja do sistema
- Ocultar na bandeja ao minimizar/fechar (se ativado).
- Clique esquerdo no ícone: restaurar a janela principal; clique direito: Mostrar janela principal / Sair.

## Perguntas frequentes

**Por que meu laptop ainda esquenta depois de fechar a tampa?**
Provavelmente o Modern Standby está ativando-o. Use Ferramentas → Testar suspensão para verificar; depois de fechar a tampa com a proteção ativa, procure por `ACTION Executando suspensão` no registro.

**O que é Modern Standby?**
Veja a explicação em linguagem simples em Ajuda → Sobre.

**Como sair completamente?**
Clique direito no ícone da bandeja → Sair (se "Fechar para a bandeja" estiver ativado, ✕ apenas minimiza).

**Executar na inicialização não funciona?**
Depende da conta de usuário atual; falhas (ex.: execução de um local restrito) são registradas com o motivo.

## Usuários avançados: reduzir ativações (opcional)

A ferramenta cuida de "suspender novamente após ser ativado". Para reduzir as ativações em si, execute em um PowerShell elevado:

```powershell
# desabilitar temporizadores de ativação
powercfg /setacvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
powercfg /setdcvalueindex SCHEME_CURRENT SUB_SLEEP 0 0
# inspecionar fontes de ativação
powercfg /waketimers
powercfg /devicequery wake_armed
# restaurar padrões
powercfg /restoredefaultschemes
```

---

## Download

Obtenha a versão mais recente em GitHub Releases: [LidReSleep Releases](https://github.com/dootn/LidReSleep/releases/latest)

| Arquivo | CPU | Download |
|---|---|---|
| `LidReSleep-amd64.exe` | Intel/AMD 64 bits (maioria dos PCs) | [amd64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-amd64.exe) |
| `LidReSleep-arm64.exe` | ARM64 (ex.: Surface Pro X, PCs Snapdragon) | [arm64](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-arm64.exe) |
| `LidReSleep-386.exe` | x86 de 32 bits | [386](https://github.com/dootn/LidReSleep/releases/latest/download/LidReSleep-386.exe) |

> Página do projeto: https://github.com/dootn/LidReSleep
