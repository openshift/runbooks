# OCSCephBlockPoolPhaseFailure

## Meaning

A `CephBlockPool` custom resource has been in the `Failure` phase for more
than 10 minutes. `CephBlockPool` resources back the RBD storage classes (for
example `ocs-storagecluster-ceph-rbd`), the Ceph metadata pool, and NVMe-oF
pools, and are created/managed by the `StorageCluster` via the Rook operator.
The Rook pool controller sets this phase when it fails to create or
reconcile the pool's configuration in the Ceph cluster (for example `ceph
osd pool create`/`ceph osd pool set` failures).

## Impact

New volumes cannot be provisioned from the storage class(es) backed by this
pool, and changes to the pool's configuration (replica size, compression,
mirroring, etc.) are not being applied. Existing PVs on the pool are
typically unaffected unless the underlying Ceph cluster itself is unhealthy.

## Diagnosis

### Check the CephBlockPool status and events

```bash
oc get cephblockpool -n openshift-storage
oc describe cephblockpool <pool_name> -n openshift-storage
```

Look at `status.phase`/`status.message` and the `Events` section at the
bottom of the `describe` output for the specific failure reason.

### Check the Rook operator logs

```bash
oc logs -n openshift-storage deploy/rook-ceph-operator | grep -i <pool_name>
```

This usually contains the underlying Ceph error (for example, a pool
replica size that cannot be satisfied, an invalid erasure-coding profile,
or the cluster rejecting the request because it is unhealthy).

### Check overall Ceph cluster health

Pool creation/reconciliation requires a responsive, healthy cluster. See
[Access the Toolbox](helpers/accessToolbox.md) and
[Troubleshoot Ceph](helpers/troubleshootCeph.md) to check `ceph status` and
`ceph health detail`.

Common root causes:

- The Ceph cluster is not yet healthy/quorate (check `CephClusterErrorState`/
  `CephClusterWarningState` alerts first if firing).
- The pool's requested replica size or failure domain cannot be satisfied by
  the current number of OSDs/failure domains (for example, requesting a
  3-way replica across zones with fewer than 3 zones available).
- OSDs are down, full, or there isn't enough raw capacity to create the pool.

## Mitigation

1. Resolve any underlying Ceph cluster health issues identified above (add
   OSD nodes/zones, bring OSDs back up, free up capacity).
2. If the pool configuration itself is incorrect (for example, an
   unsupported replica count or device class for the environment), update
   the corresponding settings in the `StorageCluster` CR.
3. Once the root cause is resolved, the pool reconciles back to `Ready`
   automatically within a few minutes; no manual phase reset is needed.
4. If the operator appears stuck despite a healthy cluster, restart it and
   re-check:
   ```bash
   oc rollout restart deployment rook-ceph-operator -n openshift-storage
   ```

If the issue persists, [gather logs](helpers/gatherLogs.md) for support.
