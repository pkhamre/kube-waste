## What's new in v0.2.0

kube-waste scans a Kubernetes cluster for wasted resources and reports estimated monthly savings — read-only, no agents, no sidecars.

### Highlights

- **More resilient scans.** If a resource kind can't be listed (RBAC, transient errors), the scan continues and skips only the affected detectors — partial results are still priced and reported.
- **A cleaner, testable pipeline.** Fetch a snapshot → run detectors → price results. Pricing is one module behind a single interface, cluster access goes through one narrow seam, and the "is this referenced?" checks now share one references module.
- **Fresh README.** New hero, terminal, and workflow diagrams (pure SVG).
- **New docs.** A domain glossary (`CONTEXT.md`) and agent documentation.

### Breaking changes

None — the CLI interface and output format are unchanged.

### Assets

- `kube-waste-linux-amd64` · `kube-waste-windows-amd64.exe`
- `kube-waste-darwin-amd64` · `kube-waste-darwin-arm64`
- `SHA256SUMS`

kube-waste is experimental and under active development. Treat savings estimates as rough guidance — verify against your own cloud bill.
