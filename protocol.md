# Voicemeeter Companion App — WebSocket Protocol

Note that this document is simply for me to note down and help me not forget the protocol. I may or may not forget to update. The source of truth is the code! Integration - companion compatibility is based on the communicated version.

## Transport

- Protocol: WebSocket (`ws://`)
- Default port: `27001`
- Endpoint: `/ws`
- Message encoding: JSON, UTF-8
- Direction labels used below: **C→H** (companion → HA), **H→C** (HA → companion)

---

## Connection Lifecycle

1. Companion app listens on `/ws`
2. HA connects on startup and reconnects automatically after any disconnect
3. On each new connection, companion app **must** immediately send a full `state` message
4. After that, companion app sends `update` messages as parameters change
5. HA sends `set` messages when the user changes something in the UI
6. Either side may close the connection at any time; HA will reconnect

---

## Message Types

### `state` — Full state dump (C→H)

Sent immediately after a client connects. Describes the complete current state

of Voicemeeter.

```json
{
  "type": "state",
  "kind": "banana",
  "version": "0.1",
  "strips": [
    {
      "index": 0,
      "label": "Stereo Input 1",
      "mute": true,
      "gain": -6,
      "virtual": false,
      "a1": false,
      "a2": false,
      "a3": true,
      "a4": false,
      "a5": false,
      "b1": false,
      "b2": false,
      "b3": false
    },
    {
      "index": 1,
      "label": "",
      "mute": false,
      "gain": 0,
      "virtual": false,
      "a1": false,
      "a2": false,
      "a3": false,
      "a4": false,
      "a5": false,
      "b1": false,
      "b2": false,
      "b3": false
    },
    {
      "index": 2,
      "label": "",
      "mute": false,
      "gain": 0,
      "virtual": false,
      "a1": false,
      "a2": false,
      "a3": false,
      "a4": false,
      "a5": false,
      "b1": false,
      "b2": false,
      "b3": false
    },
    {
      "index": 3,
      "label": "",
      "mute": false,
      "gain": 0,
      "virtual": true,
      "a1": true,
      "a2": true,
      "a3": true,
      "a4": false,
      "a5": false,
      "b1": false,
      "b2": false,
      "b3": false
    },
    {
      "index": 4,
      "label": "",
      "mute": false,
      "gain": 0,
      "virtual": true,
      "a1": true,
      "a2": true,
      "a3": true,
      "a4": false,
      "a5": false,
      "b1": true,
      "b2": true,
      "b3": false
    }
  ],
  "buses": [
    {
      "index": 0,
      "label": "",
      "mute": false,
      "gain": 0
    },
    {
      "index": 1,
      "label": "",
      "mute": true,
      "gain": 0
    },
    {
      "index": 2,
      "label": "",
      "mute": true,
      "gain": 0
    },
    {
      "index": 3,
      "label": "",
      "mute": true,
      "gain": 0
    },
    {
      "index": 4,
      "label": "",
      "mute": true,
      "gain": 0
    }
  ]
}
```

| Field    | Type   | Description                                          |
| -------- | ------ | ---------------------------------------------------- |
| `type`   | string | Always `"state"`                                     |
| `kind`   | string | Voicemeeter variant: `"basic"` `"banana"` `"potato"` |
| `strips` | array  | All input strips, ordered by index                   |
| `buses`  | array  | All output buses, ordered by index                   |

**Strip object:**

| Field       | Type   | Description                                            |
| ----------- | ------ | ------------------------------------------------------ |
| `index`     | int    | Zero-based position                                    |
| `label`     | string | User-defined name in Voicemeeter UI                    |
| `mute`      | bool   | Whether the strip is muted                             |
| `gain`      | float  | Fader level in dB. Range: `-60.0` to `+12.0`           |
| `virtual`   | bool   | `true` for virtual inputs, `false` for hardware inputs |
| `a1-5,b1-3` | bool   | Whether the strip - bus route is toggled               |

**Bus object:**

| Field   | Type   | Description                                  |
| ------- | ------ | -------------------------------------------- |
| `index` | int    | Zero-based position                          |
| `label` | string | User-defined name in Voicemeeter UI          |
| `mute`  | bool   | Whether the bus is muted                     |
| `gain`  | float  | Fader level in dB. Range: `-60.0` to `+12.0` |

---

### `update` — Single parameter changed (C→H)

Sent whenever a parameter changes in Voicemeeter, whether triggered by the

user in the Voicemeeter UI or by a `set` message from HA.

```json
{"type": "update", "target": "strip", "index": 0, "param": "mute", "value": true}
{"type": "update", "target": "bus",   "index": 2, "param": "gain", "value": -12.0}
{"type": "set", "target": "strip", "index": 0, "param": "a1", "value": true}
```

| Field    | Type          | Description                          |
| -------- | ------------- | ------------------------------------ |
| `type`   | string        | Always `"update"`                    |
| `target` | string        | `"strip"` or `"bus"` or bus labels   |
| `index`  | int           | Zero-based index of the strip or bus |
| `param`  | string        | Parameter name — see table below     |
| `value`  | bool or float | New value                            |

---

### `set` — Change a parameter (H→C)

Sent by HA when the user changes something in the UI.

```json
{"type": "set", "target": "strip", "index": 0, "param": "mute", "value": true}
{"type": "set", "target": "bus",   "index": 2, "param": "gain", "value": -6.0}
{"type": "set", "target": "strip", "index": 0, "param": "a1", "value": true}
```

| Field    | Type          | Description                          |
| -------- | ------------- | ------------------------------------ |
| `type`   | string        | Always `"set"`                       |
| `target` | string        | `"strip"` or `"bus"`                 |
| `index`  | int           | Zero-based index of the strip or bus |
| `param`  | string        | Parameter name — see table below     |
| `value`  | bool or float | Value to set                         |

After applying a `set`, the companion app should emit a corresponding `update`

message confirming the new value. This keeps any other connected clients in

sync and gives HA a confirmation that the change was applied.

---

## Valid Parameters

| target | param         | value type | range / notes         |
| ------ | ------------- | ---------- | --------------------- |
| strip  | `mute`        | bool       |                       |
| strip  | `gain`        | float      | `-60.0` to `+12.0` dB |
| strip  | `<bus label>` | bool       | `                     |
| bus    | `mute`        | bool       |                       |
| bus    | `gain`        | float      | `-60.0` to `+12.0` dB |

---

## Strip and Bus Counts by Kind

| Kind   | Hardware strips | Virtual strips | Total strips | Buses |
| ------ | --------------- | -------------- | ------------ | ----- |
| basic  | 2               | 1              | 3            | 2     |
| banana | 3               | 2              | 5            | 5     |
| potato | 5               | 3              | 8            | 8     |

---

## Error Handling

- HA does not send any acknowledgement for `update` messages
- The companion app does not send any acknowledgement for `set` messages —

  confirmation comes in the form of the subsequent `update`
- If the companion app cannot apply a `set` (e.g. Voicemeeter is not running),

  it should not send an `update` — HA will retain the last known value
- Unknown message types on either side should be silently ignored

---

## Versioning

This is version 0.1. The version will be sent together with the state message. Its a simplified semantic versioning. Any x.a version will be compatible with any x.b version. When breaking changes occur, the digit is bumped. 