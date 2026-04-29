# Deck

A lightweight web-based media controller for your system. Control volume and
media playback from any device on your local network.

## Features

- Media playback controls (play/pause, next, previous)
- Volume control and mute toggle
- QR code for easy device connection
- Cross-platform support (macOS, Linux)

## System Requirements

- macOS
    - Built-in system audio support via CoreAudio
- Linux
    - Volume control, at least one of:
      - `wpctl` (WirePlumber/PipeWire)
      - `pactl` (PulseAudio)
      - `amixer` (ALSA)
    - Media playback control:
      - `playerctl` (MPRIS/D-Bus)

## Installation

```bash
go install github.com/nobe4/deck/cmd/deck@latest
```

Or build from source:

```bash
git clone https://github.com/nobe4/deck
cd deck
go build -o deck ./cmd/deck
```

## Usage

Basic usage:
```bash
deck
deck -h
```

## Custom Templates

You can provide a custom HTML template using the `-template` flag.

See [index.html](./internal/web/static/index.html) for the default template structure.
See [server.go](./internal/web/server.go) for the available config.

## License

MIT
