# Agent Guidelines for kube-waste-go

## Build & Test Commands
- **Build**: `go build -o kube-waste .`
- **Run**: `./kube-waste`
- **Test all**: `go test ./...`
- **Test single package**: `go test ./pkg/analyzer`
- **Test with verbose**: `go test -v ./...`
- **Format**: `go fmt ./...`
- **Lint**: `go vet ./...`

## Code Style & Conventions
- **Go version**: 1.24+ with standard Go formatting (gofmt)
- **Imports**: Group stdlib, then third-party (k8s.io), then internal (github.com/pkhamre/kube-waste/pkg)
- **Error handling**: Return errors, don't panic (except fatal startup errors). Check errors immediately after function calls
- **Naming**: Use camelCase for functions/vars, PascalCase for exported types. Descriptive names (e.g., `DetectZombieDeployments`, `WasteItem`)
- **Pointers**: K8s structs use pointers for optional fields (e.g., `*dep.Spec.Replicas`). Always nil-check before dereferencing
- **Context**: Use `context.TODO()` for K8s API calls
- **Types**: Use typed constants (e.g., `WasteType string`). Prefer explicit types over `interface{}`
- **Comments**: Add doc comments to exported functions explaining what they detect/calculate
- **Cost calculation**: Use const values defined in analyzer package (CostPerGB, CostPerVCPU, etc.)

## Agent skills

### Issue tracker

Issues live as markdown files under `.scratch/<feature>/`. See `docs/agents/issue-tracker.md`.

### Triage labels

Default vocabulary: needs-triage, needs-info, ready-for-agent, ready-for-human, wontfix. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context: one `CONTEXT.md` + `docs/adr/` at the repo root. See `docs/agents/domain.md`.
