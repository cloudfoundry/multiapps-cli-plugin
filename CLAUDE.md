# multiapps-cli-plugin — CF CLI Plugin

## Project Role

Go CF CLI plugin that exposes MTA deployment commands (`cf deploy`, `cf undeploy`, etc.) to
the user. It has no business logic of its own — it translates CLI arguments into REST API calls
against the `multiapps-controller` backend and streams the results back to the terminal.

## Security Boundary

This is an **OPEN SOURCE** repository. Never introduce proprietary logic, credentials,
or internal company context into this codebase.

## Contributing Rules

- **Every commit and PR must reference a JIRA item.** No commit or PR should be merged without
  a backlog reference in its description.
- **At least one commit in the change must contain a real JIRA key in the form
  `LMCROSSITXSADEPLOY-1234`** (the `LMCROSSITXSADEPLOY-<number>` format), pointing to an actual
  Jira item — not a placeholder.

## What Claude Code MUST NOT Do Without Explicit Confirmation

Claude Code must **NOT** perform any of the following actions unless the user has explicitly
confirmed that specific action in the current conversation. Do not infer permission from an
earlier, unrelated approval — ask, then wait.

- **Do NOT push** to any remote (no `git push`, no branch/tag pushes, no opening PRs) without
  explicit confirmation.
- **Do NOT cut, tag, or trigger a release** — releases are handled by the maintainers' dedicated
  process.
- **Do NOT update version files** — `cfg/VERSION` is injected at build time via `-ldflags`; leave
  it (and any other version markers) untouched unless explicitly told to change it.
- **Do NOT bump, publish, or otherwise mutate the project version** or any similar
  project-state-changing action (release automation, changelog/version tagging, dependency
  version bumps) without explicit confirmation.
- **Do NOT hand-edit generated clients** under `clients/mtaclient/`, `clients/mtaclient_v2/`, or
  the specs in `clients/swagger/` — regenerate them instead.

When in doubt about whether an action falls into the above, treat it as requiring confirmation
and ask first.

## Tech Stack

- **Go** / go modules — repo lives at `github.com/cloudfoundry/multiapps-cli-plugin` (the `go.mod` module path is still the legacy `github.com/cloudfoundry-incubator/multiapps-cli-plugin`)
- **CF CLI plugin SDK** (`code.cloudfoundry.org/cli/v8/plugin`)
- **go-openapi** generated REST clients — do not hand-edit files under `clients/mtaclient/` or `clients/mtaclient_v2/`
- **Ginkgo v1 + Gomega** for tests — suites bootstrap via `go test` (`RunSpecs` in `*_suite_test.go`)

## Package Layout

| Package | Purpose |
|-----------------------------------|--------------------------------------------------------------|
| `multiapps_plugin.go` | Plugin entry point — `Commands` slice, `Run()`, `GetMetadata()` |
| `commands/` | One file per CF command; all embed `*BaseCommand` |
| `commands/fakes/` | Fakes for command-level interfaces |
| `clients/baseclient/` | Shared client plumbing — token factory, user-agent transport, client errors |
| `clients/mtaclient/` | go-openapi generated MTA REST client (API v1) |
| `clients/mtaclient_v2/` | go-openapi generated MTA REST client (API v2) |
| `clients/mtaclient/fakes/` | Fake v1 client builder for tests |
| `clients/mtaclient_v2/fakes/` | Fake v2 client builder for tests |
| `clients/restclient/` | Lower-level REST client (file upload, operations polling) |
| `clients/cfrestclient/` | CF-specific REST client (space/org resolution) |
| `clients/csrf/` | CSRF token handling |
| `clients/models/` | Shared model types |
| `clients/swagger/` | OpenAPI/Swagger specs (`mta_rest.yaml`, `rest.yaml`) the clients are generated from |
| `configuration/properties/` | One file per env-var-backed config property |
| `secure_parameters/` | Handling of sensitive deploy parameters |
| `util/` | CF target resolution, URL calculator, file splitter, user-agent |
| `ui/` | Terminal output helpers |
| `testutil/` | Shared test helpers (output capturer, table formatter, fake transport) |
| `cfg/VERSION` | Version string — injected via `-ldflags` at build time, do not edit manually |

## Developer Workflow

### Adding a New Command

1. **Create** `commands/<name>_command.go`:
   - Define a struct embedding `*BaseCommand`
   - Implement `GetPluginCommand()` — sets `Name`, `HelpText`, `UsageDetails.Usage`, and `UsageDetails.Options`
   - Implement `defineCommandOptions(flags *flag.FlagSet)` — declare flags
   - Implement `executeInternal(positionalArgs []string, dsHost string, flags *flag.FlagSet, cfTarget util.CloudFoundryTarget) ExecutionStatus`
   - Use `c.NewMtaClient()` or `c.NewMtaV2Client()` to get the REST client
   - Use `ui.Say()`, `ui.Ok()`, `ui.Failed()` for terminal output
   - See `commands/mtas_command.go` as the canonical simple example

2. **Register** the command in `multiapps_plugin.go` `Commands` slice — without this the command is invisible to CF CLI:
   ```go
   var Commands = []commands.Command{
       ...
       commands.NewYourCommand(),
   }
   ```

3. **Write tests** in `commands/<name>_command_test.go`:
   - Package `commands_test`, Ginkgo `Describe/Context/It` structure
   - Use `cli_fakes.NewFakeCliConnectionBuilder()` for the CF connection
   - Use fake client builders from `clients/mtaclient/fakes/` or `clients/mtaclient_v2/fakes/`
   - Use `testutil.NewUIOutputCapturer()` to capture and assert terminal output
   - Use `util_fakes.NewDeployServiceURLFakeCalculator()` for the URL calculator
   - Initialize with `command.InitializeAll(...)` — not `command.Initialize(...)`
   - See `commands/mtas_command_test.go` as the canonical simple example

Then build, test, and verify as described below.

### Build, Test & Verify (Any Change)

These steps apply to any change, not just new commands:

1. **Run tests:**
   ```bash
   go test ./commands/...
   # or all packages:
   go test ./...
   ```

2. **Format:**
   ```bash
   gofmt -w cli clients commands testutil ui util
   ```

3. **Build:**
   ```bash
   go build -ldflags "-X main.Version=$(cat cfg/VERSION)" -o multiapps-plugin .
   ```

4. **Install into CF CLI:**
   ```bash
   cf install-plugin ./multiapps-plugin -f
   ```

5. **Verify manually** against a CF environment with a running multiapps-controller (e.g. run the affected `cf <command>`).

### Full Cross-Platform Release Build

```bash
./build.sh
```

Produces static and non-static binaries for all platforms + `checksums.txt` in `build/`.

## Key Implementation Details

A few behaviors are non-obvious from the package layout — check the referenced files before touching them:

- **File upload** (`util/file_splitter.go`) — MTAR archives are split into at most **50 chunks**
  (`MaxFileChunkCount`); default chunk size **45 MB** (`MULTIAPPS_UPLOAD_CHUNK_SIZE`), uploaded in
  parallel unless `MULTIAPPS_UPLOAD_CHUNKS_SEQUENTIALLY=true`.
- **Controller URL resolution** (`util/deploy_service_url_calculator.go`) — resolved in order:
  `-u` flag → `MULTIAPPS_CONTROLLER_URL` → auto-derived from the CF API host
  (`https://api.cf.example.com` → `deploy-service.cf.example.com`).
- **REST API contract** — client targets `/api/v1/spaces/{spaceGuid}/` (full operation set) and
  `/api/v2/...` (MTA listing with namespace filtering). REST model or endpoint changes in
  `multiapps-controller` require regenerating the Go clients; breaking changes need a new version
  path (`/api/v3/`) coordinated across both repos.

## Environment Variables

| Variable | Default | Purpose |
|----------------------------------------|------------------|--------------------------------------|
| `MULTIAPPS_CONTROLLER_URL` | *(auto-derived)* | Override backend URL |
| `MULTIAPPS_UPLOAD_CHUNK_SIZE` | `45` (MB) | Size of each upload chunk |
| `MULTIAPPS_UPLOAD_CHUNKS_SEQUENTIALLY` | `false` | Upload chunks one at a time |
| `MULTIAPPS_DISABLE_UPLOAD_PROGRESS_BAR` | `false` | Suppress progress bar |
| `MULTIAPPS_USER_AGENT_SUFFIX` | *(empty)* | Append string to HTTP User-Agent |
| `DEBUG` | `false` | Set to `1` to enable HTTP request logging |

## Commands

| CF command | Go file | Alias |
|------------------------------|--------------------------------------|-------|
| `cf deploy` | `deploy_command.go` | |
| `cf bg-deploy` | `blue_green_deploy_command.go` | |
| `cf undeploy` | `undeploy_command.go` | |
| `cf mtas` | `mtas_command.go` | |
| `cf mta` | `mta_command.go` | |
| `cf mta-ops` | `mta_operations_command.go` | |
| `cf download-mta-op-logs` | `download_mta_op_logs_command.go` | `dmol` |
| `cf purge-mta-config` | `purge_config_command.go` | |
| `cf rollback-mta` | `rollback_mta_command.go` | |
