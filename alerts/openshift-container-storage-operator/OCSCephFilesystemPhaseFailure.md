# OCSCephFilesystemPhaseFailure

## Meaning

A `CephFilesystem` custom resource has been in the `Failure` phase for more
than 10 minutes. `CephFilesystem` resources back the CephFS storage classes
and are created/managed by the `StorageCluster` via the Rook operator. The
Rook controller sets this phase when it fails to create the filesystem's
metadata/data pools, run `ceph fs new`, and deploy the MDS daemons, or when
cephfs-mirror bootstrap peer configuration fails (if mirroring is enabled).

## Impact

New volumes cannot be provisioned from the CephFS storage class(es) backed
by this filesystem, and MDS daemons for it may not be running. Existing PVs
on the filesystem are typically unaffected unless the underlying Ceph
cluster itself is unhealthy.

## Diagnosis

### Check the CephFilesystem status and events

```bash
oc get cephfilesystem -n openshift-storage
oc describe cephfilesystem <filesystem_name> -n openshift-storage
```

Look at `status.phase`/`status.message` and the `Events` section at the
bottom of the `describe` output for the specific failure reason.

### Check the MDS deployments and pods

```bash
oc get deployment -n openshift-storage -l rook_file_system=<filesystem_name>
oc get pods -n openshift-storage -l app=rook-ceph-mds,rook_file_system=<filesystem_name>
```

If the MDS pods are not `Running`, describe them to check for scheduling or
resource issues.

### Check the Rook operator logs

```bash
oc logs -n openshift-storage deploy/rook-ceph-operator | grep -i <filesystem_name>
```

This usually contains the underlying Ceph error (for example, a failure to
create the metadata/data pools, an unsupported pool/replica configuration,
or a cephfs-mirror bootstrap peer secret failure).

### Check overall Ceph cluster health

See [Access the Toolbox](helpers/accessToolbox.md) and
[Troubleshoot Ceph](helpers/troubleshootCeph.md) to check `ceph status` and
`ceph health detail`.

Common root causes:

- The Ceph cluster is not yet healthy/quorate (check `CephClusterErrorState`/
  `CephClusterWarningState` alerts first if firing).
- The filesystem's metadata/data pool replica size or failure domain cannot
  be satisfied by the current number of OSDs/failure domains.
- OSDs are down, full, or there isn't enough raw capacity to create the
  pools.
- Mirroring is enabled but the cephfs-mirror bootstrap peer secret could not
  be created.

## Mitigation

1. Resolve any underlying Ceph cluster health issues identified above (add
   OSD nodes/zones, bring OSDs back up, free up capacity).
2. If the filesystem configuration itself is incorrect (for example, an
   unsupported replica count or device class for the environment), update
   the corresponding settings in the `StorageCluster` CR.
3. Once the root cause is resolved, the filesystem reconciles back to
   `Ready` automatically within a few minutes; no manual phase reset is
   needed.
4. If the operator appears stuck despite a healthy cluster, restart it and
   re-check:
   ```bash
   oc rollout restart deployment rook-ceph-operator -n openshift-storage
   ```

If the issue persists, [gather logs](helpers/gatherLogs.md) for support.
