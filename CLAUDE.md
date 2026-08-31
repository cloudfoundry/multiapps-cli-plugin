# multiapps-cli-plugin — CF CLI Plugin

## Project Role

Go CF CLI plugin that exposes MTA deployment commands (`cf deploy`, `cf undeploy`, etc.) to
the user. It has no business logic of its own — it translates CLI arguments into REST API calls
against the `multiapps-controller` backend and streams the results back to the terminal.

## Security Boundary

This is an **OPEN SOURCE** repository. Never introduce proprietary logic, credentials,
or internal company context into this codebase.

## Tech Stack

- **Go** / go modules (`github.com/cloudfoundry-incubator/multiapps-cli-plugin`)
- **CF CLI plugin SDK** (`code.cloudfoundry.org/cli/v8/plugin`)
- **go-openapi** generated REST clients — do not hand-edit files under `clients/mtaclient/` or `clients/mtaclient_v2/`
- **Ginkgo v1 + Gomega** for tests — suites bootstrap via `go test` (`RunSpecs` in `*_suite_test.go`)

## Package Layout

| Package | Purpose |
|-----------------------------------|--------------------------------------------------------------|
| `multiapps_plugin.go` | Plugin entry point — `Commands` slice, `Run()`, `GetMetadata()` |
| `commands/` | One file per CF command; all embed `*BaseCommand` |
| `commands/fakes/` | Fakes for command-level interfaces |
| `clients/mtaclient/` | go-openapi generated MTA REST client (API v1) |
| `clients/mtaclient_v2/` | go-openapi generated MTA REST client (API v2) |
| `clients/mtaclient/fakes/` | Fake v1 client builder for tests |
| `clients/mtaclient_v2/fakes/` | Fake v2 client builder for tests |
| `clients/restclient/` | Lower-level REST client (file upload, operations polling) |
| `clients/cfrestclient/` | CF-specific REST client (space/org resolution) |
| `clients/csrf/` | CSRF token handling |
| `clients/models/` | Shared model types |
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

4. **Run tests:**
   ```bash
   go test ./commands/...
   # or all packages:
   go test ./...
   ```

5. **Format:**
   ```bash
   gofmt -w cli clients commands testutil ui util
   ```

6. **Build:**
   ```bash
   go build -ldflags "-X main.Version=$(cat cfg/VERSION)" -o multiapps-plugin .
   ```

7. **Install into CF CLI:**
   ```bash
   cf install-plugin ./multiapps-plugin -f
   ```

8. **Verify manually** by running `cf <your-command>` against a CF environment with a running multiapps-controller.

### Full Cross-Platform Release Build

```bash
./build.sh
```

Produces static and non-static binaries for all platforms + `checksums.txt` in `build/`.

## File Upload

MTAR archives are split into chunks before upload:

- Default chunk size: **45 MB** — controlled by `MULTIAPPS_UPLOAD_CHUNK_SIZE`
- Maximum **50 chunks** (`MaxFileChunkCount` in `util/file_splitter.go`) — minimum chunk size is enforced as `ceil(mtar_size_mb / 50)`
- Chunks upload **in parallel** by default — set `MULTIAPPS_UPLOAD_CHUNKS_SEQUENTIALLY=true` to serialize

## Controller URL Resolution

Resolved in this priority order (`util/deploy_service_url_calculator.go`):

1. `-u` flag on the command line
2. `MULTIAPPS_CONTROLLER_URL` environment variable
3. Auto-derived: strips the CF API host up to the first `.` and prepends `deploy-service.`
   (e.g. `https://api.cf.example.com` → `deploy-service.cf.example.com`)

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

## REST API Contract

- **v1** `/api/v1/spaces/{spaceGuid}/` — full operation set
- **v2** `/api/v2/spaces/{spaceGuid}/` — MTA listing with namespace filtering

REST model or endpoint changes in `multiapps-controller` require regenerating the Go clients.
Breaking changes require a new API version path (`/api/v3/`) with coordinated updates in both repos.