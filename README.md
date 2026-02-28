# voicemeeter-ha-windows-companion

A small Windows background app that bridges [Voicemeeter](https://vb-audio.com/Voicemeeter/) to your Home Assistant instance over WebSocket. It's the required companion to the [ha-voicemeeter](https://github.com/Skrubbadubba/ha-voicemeeter) integration.

---

## Requirements

- Windows 10 or later
- [Voicemeeter](https://vb-audio.com/Voicemeeter/) (any edition: Basic, Banana, or Potato) installed and running

---

## Installation

1. Download `voicemeeter-companion.exe` from the [latest release](https://github.com/Skrubbadubba/voicemeeter-ha-windows-companion/releases/latest).
2. Place it somewhere permanent — for example `C:\Program Files\VoicemeeterCompanion\`.
3. Make sure Voicemeeter is running, then launch `voicemeeter-companion.exe`.

The app runs silently in the background and appears in the system tray. It starts a WebSocket server on port **27001** that the Home Assistant integration connects to.

---

## Auto-start on Windows

To have the companion start automatically when you log in:

1. Right-click `voicemeeter-companion.exe` and select **Create shortcut**.
2. Press <kbd>Win</kbd> + <kbd>R</kbd>, type `shell:startup`, and press <kbd>Enter</kbd>.
3. Move the shortcut into that folder.

> **Note:** Voicemeeter itself should also be set to start on login, and it must be running before the companion app starts.

---

## Building from source

You'll need [Go](https://go.dev/dl/) installed.

[Just](https://github.com/casey/just) is optional, the scripts are easy enough. You can install it on windows with winget:

```ps
winget install --id Casey.Just --exact
```

```ps
git clone https://github.com/Skrubbadubba/voicemeeter-ha-windows-companion
cd voicemeeter-ha-windows-companion
just build
```

This produces `dist/voicemeeter-companion.exe` — a console build with a visible terminal window with.

---

## How it works

The companion app reads Voicemeeter's state by calling the official `VoicemeeterRemote64.dll` that ships with Voicemeeter. It exposes that state over a local WebSocket server. When Home Assistant connects, it receives a full state dump immediately, followed by push updates whenever something changes. Commands from HA (muting a strip, adjusting gain, etc.) are sent back the other way.

See [PROTOCOL.md](PROTOCOL.md) for a more in-depth description.