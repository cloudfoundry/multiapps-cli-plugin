# multiapps-cli-plugin — CLAUDE.md

## Project Role

This is the **multiapps-cli-plugin**: a Go-based Cloud Foundry CLI plugin that extends
the CF CLI (v8) with commands for managing MTA deployments (`deploy`, `undeploy`,
`bg-deploy`, `mtas`, `mta-ops`, `dmol`). It communicates with the
**multiapps-controller** backend via auto-generated OpenAPI clients (go-openapi/runtime).

## Security Boundary

**This is an OPEN SOURCE repository.**
Never introduce proprietary logic, credentials, or internal company context into this
codebase. No SAP-internal service integrations, no hardcoded tokens or URLs.

## Tech Stack

- **Language**: Go 1.25.4 (go modules — `go.mod` is the source of truth)
- **CF CLI integration**: `code.cloudfoundry.org/cli/v8`
- **REST clients**: auto-generated via `swagger generate client` (see `regen-client.sh`)
  — do NOT hand-edit files under `clients/`
- **Testing**: Ginkgo v1 + Gomega (`github.com/onsi/ginkgo`, `github.com/onsi/gomega`)

## Key Directory Layout

```
commands/        — CF CLI command implementations (deploy, undeploy, etc.)
clients/         — AUTO-GENERATED OpenAPI REST clients (v1 + v2); never edit manually
cli/             — CF CLI plugin registration and bootstrapping
util/            — shared helpers (HTTP, logging, formatting)
secure_parameters/ — secrets collection logic (--collect-secrets flag)
cfg/VERSION      — authoritative version string consumed by build.sh
regen-client.sh  — regenerate clients/ from multiapps-controller Swagger specs
```

## Build Commands

**Local dev binary (current platform):**
```bash
go build -ldflags "-X main.Version=$(cat cfg/VERSION)" -o multiapps-plugin
```

**All platforms (cross-compiled, output to `build/`):**
```bash
./build.sh
```
`build.sh` produces both static (`CGO_ENABLED=0`) and non-static variants for
linux/32, linux/64, linux/arm64, windows/32, windows/64, darwin/amd64, darwin/arm64.

## Running Tests

```bash
go test ./...
```
Or with Ginkgo's verbose output:
```bash
ginkgo -r ./...
```

## Regenerating the Go REST Client

Run `regen-client.sh` after the multiapps-controller Swagger specs change. The script
rebuilds the controller API module, then calls `swagger generate client` for both v1
and v2. Copy the generated output into `clients/`.

## Adding CLI Flags

New flags must be registered in the command's `GetPluginCommand()` `Options` map and
documented in the CF CLI help text. Update the parent `CLAUDE.md` documentation
checklist when adding user-visible flags.

## Mandatory Formatting Rule

**Before completing any task or committing code, you MUST run:**
```bash
go fmt ./...
```
All committed Go source must conform to standard `gofmt` formatting.