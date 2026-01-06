# Baileys Hammer - task runner
#
# Usage:
#   just --list
#   just dev
#
# Notes:
# - This `justfile` intentionally does NOT auto-load `.env` (to avoid surprises from global/local env).
# - `DATABASE_URL` defaults to a local sqlite file in `./tmp/data/`.

set shell := ["bash", "-cu"]

APP_NAME := "baileys-hammer"
PORT := "8080"
DATABASE_URL := "./tmp/data/devTEST.db"
PASS := "pass"

default:
  @just --list

help:
  @just --list

check:
  @command -v go >/dev/null || (echo "missing tool: go" >&2; exit 1)

check-dev: check
  @command -v templ >/dev/null || (echo "missing tool: templ (see: just install-tools)" >&2; exit 1)
  @command -v air >/dev/null || (echo "missing tool: air (see: just install-tools)" >&2; exit 1)

install-tools:
  @echo "Installs go tools into your GOPATH/bin (requires network):"
  @echo "  go install github.com/a-h/templ/cmd/templ@v0.3.960"
  @echo "  go install github.com/air-verse/air@latest"

dirs:
  @mkdir -p ./tmp ./tmp/data

gen: dirs
  templ generate

watch-templ: dirs
  templ generate --watch

run: dirs
  @echo "DATABASE_URL={{DATABASE_URL}}"
  @echo "PASS={{PASS}}"
  DATABASE_URL="{{DATABASE_URL}}" PASS="{{PASS}}" PORT="{{PORT}}" go run .

dev: check-dev dirs
  @echo "Starting dev server on :{{PORT}} (ctrl-c to stop)"
  @echo "DATABASE_URL={{DATABASE_URL}}"
  @echo "PASS={{PASS}}"
  @trap 'jobs -pr | xargs -r kill 2>/dev/null || true' EXIT INT TERM; \
    export DATABASE_URL="{{DATABASE_URL}}" PASS="{{PASS}}" PORT="{{PORT}}"; \
    templ generate --watch & \
    air

build: gen dirs
  go build -o ./tmp/main .

test: check
  go test ./...

fmt: check
  go fmt ./...

fmt-check: check
  @test -z "$$(gofmt -l $$(find . -name '*.go' -not -path './tmp/*' -not -path './vendor/*'))" || (echo "gofmt needed; run: just fmt" >&2; exit 1)

tidy: check
  go mod tidy

docker-build:
  docker build -t {{APP_NAME}} .

docker-run: dirs
  docker run --rm -p 8080:8080 \
    -e DATABASE_URL="{{DATABASE_URL}}" \
    -e PASS="{{PASS}}" \
    -v "$$(pwd)/tmp/data:/tmp/data" \
    {{APP_NAME}}

fly-deploy:
  fly deploy

fly-deploy-dev:
  fly deploy -c fly-dev.toml

nix-shell:
  nix develop .#devShells.shell

nix-dev:
  nix develop .#devShells.dev

nix-build:
  nix develop .#devShells.build

nix-deploy:
  nix develop .#devShells.deploy

nix-docker:
  nix develop .#devShells.dockerBuild

