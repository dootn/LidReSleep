#!/usr/bin/env sh
# Cross-compiles LidReSleep.exe for Windows.
# By default builds amd64, arm64 and 386; set ARCH to build a single
# architecture, e.g. ARCH=arm64 docker/build.sh.
# Uses Docker for isolation; module/build caches are persisted under .cache/.
set -eu

ROOT="$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)"
cd "$ROOT"

mkdir -p bin .cache/gomod .cache/gocache

ARCHS="${ARCH:-amd64,arm64,386}"

docker run --rm \
  -u "$(id -u):$(id -g)" \
  -e HOME=/tmp/home \
  -e GOFLAGS=-buildvcs=false \
  -e GOMODCACHE=/tmp/home/gomod \
  -e GOCACHE=/tmp/home/gocache \
  -e CGO_ENABLED=0 \
  -e GOOS=windows \
  -e ARCHS="$ARCHS" \
  -v "$ROOT":/workspace -w /workspace \
  -v "$ROOT/.cache/gomod":/tmp/home/gomod \
  -v "$ROOT/.cache/gocache":/tmp/home/gocache \
  golang:1.23-alpine \
  sh -c '
    set -eu
    go mod download
    for arch in $(echo "$ARCHS" | tr "," " "); do
      # Embed the manifest (common controls v6 + DPI awareness) and the app icon
      # (multi-size .ico) for the target architecture. rsrc runs on the host
      # architecture; only its -arch matters for the COFF output.
      env -u GOOS -u GOARCH go run github.com/akavel/rsrc@v0.10.2 -arch "$arch" -manifest docker/manifest.xml -ico internal/gui/icon.ico -o cmd/lidsleep/rsrc.syso
      out="bin/LidReSleep-$arch.exe"
      GOOS=windows GOARCH=$arch go build -trimpath -ldflags "-s -w -H=windowsgui" -o "$out" ./cmd/lidsleep
      echo "OK -> $out"
    done
  '

ls -lh "$ROOT"/bin/LidReSleep*.exe
