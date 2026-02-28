set windows-shell := ["powershell.exe", "-NoProfile", "-Command"]

default: build

# Development build (console window visible for debugging)
build:
    New-Item -ItemType Directory -Force -Path dist | Out-Null
    go build -o dist/voicemeeter-companion.exe .

# Release build (no console window, for GitHub releases)
build-release:
    New-Item -ItemType Directory -Force -Path dist | Out-Null
    go build -ldflags "-H=windowsgui" -o dist/voicemeeter-companion.exe .

run:
    go run -ldflags "-X main.PROTOCOL_VER=dev" .

clean:
    Remove-Item -Recurse -Force -ErrorAction SilentlyContinue dist

tidy:
    go mod tidy