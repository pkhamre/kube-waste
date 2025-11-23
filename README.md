# kube-waste

> Stop burning money on Kubernetes. Find unused PVCs, abandoned Load Balancers, and zombie workloads.

A lightweight CLI tool written in Go that scans your Kubernetes cluster for wasteful resources and estimates potential cost savings. **No agents. No sidecars. Read-only.**

## Features

- 🔍 **Detects 4 types of waste:**
  - Unused PersistentVolumeClaims (PVCs)
  - Services with no backing pods
  - Zombie deployments (scaled to 0 replicas)
  - Orphaned pods (running without owner references)

- 💵 **Real dollar estimates** - Calculates monthly savings based on standard cloud pricing
- ⚡ **Fast** - Scans thousands of resources in milliseconds
- 🔒 **100% Safe** - Read-only access using your local kubeconfig
- 🚀 **Zero installation** - Single binary, no cluster-side components

## Installation

### Download Binary

Download the latest release for your platform from the [releases page](https://github.com/pkhamre/kube-waste/releases).

```bash
# Linux/macOS
curl -LO https://github.com/pkhamre/kube-waste/releases/latest/download/kube-waste
chmod +x kube-waste
sudo mv kube-waste /usr/local/bin/
```

### Build from Source

Requires Go 1.24 or later:

```bash
git clone https://github.com/pkhamre/kube-waste.git
cd kube-waste
go build -o kube-waste .
```

## Usage

Simply run the binary - it will automatically detect your kubeconfig:

```bash
./kube-waste
```

### Example Output

```
ℹ️  Using config from current directory: kubeconfig
🔍 Scanning cluster for waste (Go Version)...

TYPE               NAMESPACE    NAME                 DETAILS                  EST. SAVINGS
----               ---------    ----                 -------                  ------------
Unused PVC         default      old-mongo-backup     50Gi (Bound)             $5.00/mo
Unused Service     payments     legacy-lb-api        LoadBalancer (Unused)    $15.00/mo
Zombie Deployment  backend      worker-v1            Replicas: 0              $0.00/mo
Orphaned Pod       default      debug-pod            0.1 vCPU / 0.1 GB        $2.30/mo
-------------------------------------------------------------------------
💰 TOTAL POTENTIAL SAVINGS: $22.30 / month
```

## Kubeconfig Resolution

The tool looks for kubeconfig in this order:

1. `./kubeconfig` (current directory)
2. `./.kubeconfig` (current directory, hidden)
3. `~/.kube/config` (home directory)
4. `$KUBECONFIG` environment variable

## Detection Logic

### Unused PVCs
- Lists all PersistentVolumeClaims
- Checks if any pod mounts them
- Reports unattached PVCs with storage size
- **Cost**: $0.10/GB/month

### Unused Services
- Lists all Services (ClusterIP, NodePort, LoadBalancer)
- Queries EndpointSlices for backing pods
- Reports services with zero endpoints
- **Cost**: $15.00/month for LoadBalancers

### Zombie Deployments
- Lists all Deployments
- Identifies those with `replicas: 0`
- Suggests cleanup of scaled-down workloads
- **Cost**: $0.00 (operational waste)

### Orphaned Pods
- Lists all running Pods
- Identifies those without OwnerReferences
- Calculates resource usage (CPU/Memory)
- **Cost**: $20.00/vCPU + $3.00/GB RAM per month

## Development

### Project Structure

```
.
├── main.go                  # CLI entry point & kubeconfig handling
├── pkg/
│   └── analyzer/
│       ├── types.go         # WasteItem struct & constants
│       ├── deployment.go    # Zombie deployment detection
│       ├── pod.go          # Orphaned pod detection
│       ├── pvc.go          # Unused PVC detection
│       └── service.go      # Unused service detection
└── AGENTS.md               # Development guidelines
```

### Build & Test

```bash
# Build
go build -o kube-waste .

# Run
./kube-waste

# Test
go test ./...

# Format
go fmt ./...

# Lint
go vet ./...
```

See [AGENTS.md](AGENTS.md) for detailed coding guidelines.

## Cost Calculation Assumptions

Default pricing estimates (configurable in `pkg/analyzer/`):

- **Storage (PVC)**: $0.10/GB/month
- **LoadBalancer**: $15.00/month
- **Compute (vCPU)**: $20.00/month
- **Memory (RAM)**: $3.00/GB/month

These are conservative estimates based on major cloud providers (AWS/GCP/Azure). Actual costs may vary.

## Limitations

- **Metrics-free**: Does not analyze actual usage patterns, only configuration
- **Stateless**: No historical tracking between runs
- **Pricing estimates**: Based on generic cloud pricing, not your specific rates
- **Permissions**: Requires read access to pods, services, deployments, PVCs, and endpoint slices

## Why Go?

- **Single binary distribution** - No runtime dependencies
- **Fast startup** - Sub-millisecond initialization
- **Native K8s support** - Official `client-go` library
- **Cross-platform builds** - Easy to compile for Linux/macOS/Windows

## Contributing

Contributions welcome! Please:

1. Fork the repository
2. Create a feature branch
3. Follow the guidelines in [AGENTS.md](AGENTS.md)
4. Submit a pull request

## License

MIT License - see [LICENSE](LICENSE) for details

## Roadmap

- [ ] JSON/CSV export formats
- [ ] Custom cost configuration via flags
- [ ] Historical tracking & trend analysis
- [ ] Automated Slack/email alerts
- [ ] Support for StatefulSets and DaemonSets
- [ ] Integration with cloud provider APIs for exact pricing

---

**Disclaimer**: This tool provides estimates only. Always verify actual costs with your cloud provider before making infrastructure changes.
