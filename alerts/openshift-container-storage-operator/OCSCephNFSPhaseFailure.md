# OCSCephNFSPhaseFailure

## Meaning

A `CephNFS` custom resource has been in a `Failed` or `ReconcileFailed`
phase for more than 10 minutes. `CephNFS` resources represent an NFS
Ganesha server exporting CephFS data, created/managed by the
`StorageCluster` via the Rook operator.

Note: the Rook NFS controller currently only sets `status.phase` to empty
(on first creation) or `Ready` (on reconcile success); reconcile errors are
reported only as Kubernetes `Events`/operator log lines, not as a `Failed`/
`ReconcileFailed` phase. This alert is therefore defense-in-depth for a
future/alternate Rook behavior and will rarely fire today. If the NFS
server is unhealthy but this alert is not firing, rely on the pod- and
export-level checks below instead of `status.phase`.

## Impact

Clients connected to the affected NFS export can lose access, and new NFS
exports/shares for this `CephNFS` may fail to provision. Existing CephFS
data is not affected, only the NFS Ganesha gateway in front of it.

## Diagnosis

### Check the CephNFS status and events

```bash
oc get cephnfs -n openshift-storage
oc describe cephnfs <nfs_name> -n openshift-storage
```

### Check the NFS Ganesha pods

```bash
oc get pods -n openshift-storage -l app=rook-ceph-nfs,ceph_nfs=<nfs_name>
oc logs -n openshift-storage -l app=rook-ceph-nfs,ceph_nfs=<nfs_name>
```

Look for `CrashLoopBackOff`/restarts, or Ganesha export errors in the pod
logs.

### Check the Rook operator logs

```bash
oc logs -n openshift-storage deploy/rook-ceph-operator | grep -i <nfs_name>
```

### Check overall Ceph cluster health

The CephFS backing pools and MDS daemons must be healthy for NFS exports
to work. See [Access the Toolbox](helpers/accessToolbox.md) and
[Troubleshoot Ceph](helpers/troubleshootCeph.md) to check `ceph status`,
`ceph health detail`, and `ceph nfs cluster info <nfs_name>`.

## Mitigation

1. Resolve any underlying Ceph cluster or MDS health issues identified
   above.
2. If the Ganesha pod(s) are crash-looping, check the pod logs for a
   configuration error (invalid export path, unsupported squash/access
   mode) and correct the `CephNFS` spec via the `StorageCluster` CR.
3. Once the root cause is resolved, the resource reconciles back to
   `Ready` automatically; no manual phase reset is needed.
4. If the operator appears stuck despite a healthy cluster, restart it and
   re-check:
   ```bash
   oc rollout restart deployment rook-ceph-operator -n openshift-storage
   ```

If the issue persists, [gather logs](helpers/gatherLogs.md) for support.
