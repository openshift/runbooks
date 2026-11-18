# OCSCephClusterPhaseFailure

## Meaning

A `CephCluster` custom resource has been in the `Failure` phase for more
than 10 minutes. `CephCluster` is the top-level CR that represents the
entire Ceph cluster (mons, mgrs, OSDs) and is created/managed by the
`StorageCluster` via the Rook operator. Every other Ceph CR (pools,
filesystems, object stores, NFS, etc.) depends on this cluster being
healthy in order to reconcile successfully.

The Rook cluster controller sets this phase when it fails to validate or
reconcile the cluster, including:

- Mon quorum failing to start/establish.
- OSD daemons failing to start or prepare.
- The running/desired Ceph version check failing (unsupported version, or
  an attempted downgrade).
- An invalid cluster spec (incorrect mon count for the number of available
  nodes, invalid stretch-cluster configuration, invalid network spec, or
  invalid KMS/encryption configuration).
- Failure to connect to or import an external Ceph cluster (for
  external-mode deployments).

## Impact

Because `CephCluster` is the root CR, impact ranges from mild to severe
depending on which step failed:

- If mon quorum or OSD start fails, the cluster may be unable to serve I/O
  at all — existing PVs can become unavailable and all CSI provisioning
  will fail.
- If the failure is limited to version detection or spec validation,
  existing storage I/O is typically unaffected, but no further
  configuration changes or upgrades will be applied.
- Other CR controllers (CephBlockPool, CephFilesystem, CephObjectStore,
  CephNFS, etc.) will themselves start failing to reconcile while the
  cluster is unhealthy, since they depend on it.

## Diagnosis

### Check the CephCluster status and events

```bash
oc get cephcluster -n openshift-storage
oc describe cephcluster <cluster_name> -n openshift-storage
```

Check `status.phase`/`status.message` and the `Events` section for the
specific failure reason.

### Check cluster-reported Ceph health

```bash
oc get cephcluster <cluster_name> -n openshift-storage -o jsonpath='{.status.ceph}'
```

This is populated independently by a periodic health check and reflects
`ceph health detail`-equivalent output, even if the overall reconcile is
stuck.

### Check mon, mgr, and OSD pods

```bash
oc get pods -n openshift-storage -l app=rook-ceph-mon
oc get pods -n openshift-storage -l app=rook-ceph-mgr
oc get pods -n openshift-storage -l app=rook-ceph-osd
oc get pods -n openshift-storage -l app=rook-ceph-osd-prepare
```

Mon quorum and OSD start/prepare failures are the most common root
causes; the operator log only reports that these steps failed, not the
daemon-level root cause, so check these pods' logs directly for crashes
or scheduling issues.

### Check the Rook operator logs

```bash
oc logs -n openshift-storage deploy/rook-ceph-operator | grep -i <cluster_name>
```

Look for `failed to start ceph monitors`, `failed to start ceph osds`,
`failed the ceph version check`, or a cluster spec validation error —
these correspond directly to the distinct error paths that can set this
phase.

### Check overall Ceph cluster health

See [Access the Toolbox](helpers/accessToolbox.md) and
[Troubleshoot Ceph](helpers/troubleshootCeph.md) to check `ceph status`
and `ceph health detail` directly.

## Mitigation

1. If mon quorum or OSD daemons are failing to start, check node health,
   storage device availability, and scheduling constraints (node
   selectors, taints, PodDisruptionBudgets) before escalating.
2. If the Ceph version check failed, confirm the configured Ceph image is
   a supported version and not an unsupported downgrade.
3. If a spec validation error is reported (mon count, stretch cluster,
   network, or KMS/encryption configuration), correct the corresponding
   setting in the `StorageCluster` CR.
4. For external-mode clusters, verify connectivity and credentials to the
   external Ceph cluster.
5. Once the root cause is resolved, the cluster reconciles back to
   `Ready` automatically; no manual phase reset is needed.
6. If the operator appears stuck despite a healthy cluster, restart it and
   re-check:
   ```bash
   oc rollout restart deployment rook-ceph-operator -n openshift-storage
   ```

If the issue persists, [gather logs](helpers/gatherLogs.md) for support.
