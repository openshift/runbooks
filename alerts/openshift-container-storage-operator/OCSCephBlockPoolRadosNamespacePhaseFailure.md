# OCSCephBlockPoolRadosNamespacePhaseFailure

## Meaning

A `CephBlockPoolRadosNamespace` custom resource has been in the `Failure`
phase for more than 10 minutes. Each `StorageConsumer` gets a dedicated RADOS
namespace inside the shared `CephBlockPool`, created automatically by the
primary consumer's workflow (not created directly by a user). The Rook
controller sets this phase when it fails to create or update the RADOS
namespace in Ceph (for example `rbd namespace create` failures), or when
mirroring configuration for the namespace fails to apply.

## Impact

The affected storage consumer/client cannot get a working, isolated RBD
namespace in the shared pool. New RBD volume provisioning for that specific
consumer can fail or stall, even though the parent `CephBlockPool` itself may
be healthy.

## Diagnosis

### Check the CephBlockPoolRadosNamespace status and events

```bash
oc get cephblockpoolradosnamespace -n openshift-storage
oc describe cephblockpoolradosnamespace <radosnamespace_name> -n openshift-storage
```

Note the `spec.blockPoolName` (the parent pool) and the owning
`StorageConsumer`, and check `status.phase` plus the `Events` section for the
failure reason.

### Check the parent CephBlockPool is Ready

The RADOS namespace cannot be created until its parent pool is `Ready`:

```bash
oc get cephblockpool <pool_name> -n openshift-storage -o jsonpath='{.status.phase}'
```

If the parent pool itself is failing, troubleshoot
[OCSCephBlockPoolPhaseFailure](OCSCephBlockPoolPhaseFailure.md) first.

### Check the Rook operator logs

```bash
oc logs -n openshift-storage deploy/rook-ceph-operator | grep -i <radosnamespace_name>
```

This contains the underlying Ceph error from the `rbd namespace create`/
update call, or a mirroring configuration failure if
`spec.mirroring` is set.

### Check overall Ceph cluster health

See [Access the Toolbox](helpers/accessToolbox.md) and
[Troubleshoot Ceph](helpers/troubleshootCeph.md) to check `ceph status` and
`ceph health detail`.

## Mitigation

1. Resolve the parent `CephBlockPool`'s health first if it isn't `Ready`.
2. Resolve any underlying Ceph cluster health issues identified above (OSDs
   down, cluster full, unresponsive mons).
3. Once the root cause is resolved, the RADOS namespace reconciles back to
   `Ready` automatically; no manual phase reset is needed.
4. If the operator appears stuck despite a healthy cluster, restart it and
   re-check:
   ```bash
   oc rollout restart deployment rook-ceph-operator -n openshift-storage
   ```
5. If the resource remains stuck in `Failure` with no clear Ceph-side error,
   escalate with the owning `StorageConsumer` name and the operator logs
   rather than deleting the `CephBlockPoolRadosNamespace` directly, since it
   is managed automatically by the consumer's onboarding workflow.

If the issue persists, [gather logs](helpers/gatherLogs.md) for support.
