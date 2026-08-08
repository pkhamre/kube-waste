# Domain Context

## Glossary

- **kube-waste** — a read-only CLI that scans a Kubernetes cluster for wasted resources and reports estimated monthly savings.
- **Scan** — one pass over the cluster: fetch a snapshot, run every detector, aggregate results. The module `Scan(ctx, clientset) ScanResult`.
- **Snapshot** — the immutable cluster state a scan runs against: one fetch of each resource kind (deployments, pods, PVCs, services, endpoint slices), plus which kinds failed to list. Detectors are pure functions over it.
- **Kind** — a cluster resource kind the snapshot can hold (deployments, pods, persistentvolumeclaims, services, endpointslices). A kind that failed to list is *unavailable*; detectors that need it are skipped, never run against empty data.
- **Detector** — a pure function that finds one type of waste inside a snapshot. Registered in a table with the kinds it needs.
- **Waste item** — one piece of detected waste: type, namespace, name, details, estimated monthly cost.
- **Waste type** — the four things this tool detects: **Unused PVC**, **Unused Service**, **Orphaned Pod**, **Zombie Deployment**.
  - *Unused PVC* — a Bound or Pending PersistentVolumeClaim no running pod mounts.
  - *Unused Service* — a Service with no backing pods (no endpoint slice endpoints), skipping the `kubernetes` service.
  - *Orphaned Pod* — a Running pod with no owner reference.
  - *Zombie Deployment* — a Deployment scaled to zero replicas.
- **Potential savings** — the sum of estimated monthly costs across waste items; the headline output.
- **Estimated** — a cost computed from default resource assumptions (no requests set) rather than observed values.

## Decisions

- See `docs/adr/` for architecture decision records (none yet).
