set windows-shell := ["powershell.exe", "-NoProfile", "-Command"]

default: build

build:
    New-Item -ItemType Directory -Force -Path dist | Out-Null
    go build -ldflags "-H=windowsgui" -o dist/voicemeeter-companion-background.exe .
    go build -o dist/voicemeeter-companion.exe .

build-release version:
    New-Item -ItemType Directory -Force -Path dist | Out-Null
    go build -ldflags "-X main.Version={{version}}" -o dist/voicemeeter-companion-{{version}}.exe .

run:
    go run -ldflags "-X main.PROTOCOL_VER=dev" . 

clean:
    Remove-Item -Recurse -Force -ErrorAction SilentlyContinue dist

tidy:
    go mod tidy