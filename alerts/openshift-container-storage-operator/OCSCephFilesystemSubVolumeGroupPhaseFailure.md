# OCSCephFilesystemSubVolumeGroupPhaseFailure

## Meaning

A `CephFilesystemSubVolumeGroup` custom resource has been in the `Failure`
phase for more than 10 minutes. Each `StorageConsumer` gets a dedicated
subvolume group inside the shared `CephFilesystem`, created automatically by
the consumer's onboarding workflow (not created directly by a user). The
Rook controller sets this phase when it fails to create or update the
subvolume group in Ceph (for example `ceph fs subvolumegroup create`
failures).

## Impact

The affected storage consumer/client cannot get a working, isolated CephFS
subvolume group. New CephFS volume provisioning for that specific consumer
can fail or stall, even though the parent `CephFilesystem` itself may be
healthy.

## Diagnosis

### Check the CephFilesystemSubVolumeGroup status and events

```bash
oc get cephfilesystemsubvolumegroup -n openshift-storage
oc describe cephfilesystemsubvolumegroup <subvolumegroup_name> -n openshift-storage
```

Note the `spec.filesystemName` (the parent filesystem) and the owning
`StorageConsumer`, and check `status.phase` plus the `Events` section for
the failure reason.

### Check the parent CephFilesystem is Ready

The subvolume group cannot be created until its parent filesystem is
`Ready`:

```bash
oc get cephfilesystem <filesystem_name> -n openshift-storage -o jsonpath='{.status.phase}'
```

If the parent filesystem itself is failing, troubleshoot
[OCSCephFilesystemPhaseFailure](OCSCephFilesystemPhaseFailure.md) first.

### Check the Rook operator logs

```bash
oc logs -n openshift-storage deploy/rook-ceph-operator | grep -i <subvolumegroup_name>
```

This contains the underlying Ceph error from the `ceph fs subvolumegroup
create`/update call (for example, an invalid quota, pinning, or data pool
layout setting).

### Check overall Ceph cluster health

See [Access the Toolbox](helpers/accessToolbox.md) and
[Troubleshoot Ceph](helpers/troubleshootCeph.md) to check `ceph status` and
`ceph health detail`.

## Mitigation

1. Resolve the parent `CephFilesystem`'s health first if it isn't `Ready`.
2. Resolve any underlying Ceph cluster health issues identified above (OSDs
   down, cluster full, unresponsive mons).
3. Once the root cause is resolved, the subvolume group reconciles back to
   `Ready` automatically; no manual phase reset is needed.
4. If the operator appears stuck despite a healthy cluster, restart it and
   re-check:
   ```bash
   oc rollout restart deployment rook-ceph-operator -n openshift-storage
   ```
5. If the resource remains stuck in `Failure` with no clear Ceph-side error,
   escalate with the owning `StorageConsumer` name and the operator logs
   rather than deleting the `CephFilesystemSubVolumeGroup` directly, since it
   is managed automatically by the consumer's onboarding workflow.

If the issue persists, [gather logs](helpers/gatherLogs.md) for support.
