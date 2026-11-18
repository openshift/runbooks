# OCSCephRBDMirrorPhaseFailure

## Meaning

A `CephRBDMirror` custom resource has been in a `Failed` phase for more
than 10 minutes. `CephRBDMirror` deploys the `rbd-mirror` daemon(s) that
replicate RBD images to a peer Ceph cluster for disaster recovery, managed
by the `StorageCluster`/Rook operator.

The Rook rbd-mirror controller sets `Failed` as a single catch-all for any
error encountered while reconciling, including: an invalid `CephRBDMirror`
spec, failure to load cluster info, failure to detect the running/desired
Ceph version, failure to add a bootstrap peer (`spec.peers`), a cephx
daemon-key rotation failure, or failure to create/update the rbd-mirror
daemon deployment(s). The phase does not distinguish which of these
occurred; check the operator logs for the specific cause.

## Impact

RBD image mirroring to the peer cluster stops making progress or may not
start at all. This does not affect local I/O to RBD volumes, but it puts
disaster-recovery replication at risk — mirrored images can fall behind or
stop syncing entirely, increasing potential data loss (RPO) if a failover
to the peer cluster is needed while this is unresolved.

## Diagnosis

### Check the CephRBDMirror status and events

```bash
oc get cephrbdmirror -n openshift-storage
oc describe cephrbdmirror <mirror_name> -n openshift-storage
```

### Check the rbd-mirror daemon pods

```bash
oc get pods -n openshift-storage -l app=rook-ceph-rbd-mirror
oc logs -n openshift-storage -l app=rook-ceph-rbd-mirror
```

### Check bootstrap peer and mirroring status

```bash
oc get cephrbdmirror <mirror_name> -n openshift-storage -o jsonpath='{.spec.peers}'
```

From the [Toolbox](helpers/accessToolbox.md), check the actual mirroring
state in Ceph, independent of the CR's reconcile status:

```bash
ceph rbd mirror pool status <pool_name> --verbose
```

### Check the Rook operator logs

```bash
oc logs -n openshift-storage deploy/rook-ceph-operator | grep -i <mirror_name>
```

Look for `failed to add ceph rbd mirror peer`, `failed to populate cluster
info`, `failed to detect running and desired ceph version`, or `failed to
create ceph rbd mirror deployments` — these correspond directly to the
distinct error paths that set `Failed`.

### Check overall Ceph cluster health

See [Access the Toolbox](helpers/accessToolbox.md) and
[Troubleshoot Ceph](helpers/troubleshootCeph.md) to check `ceph status`
and `ceph health detail`.

## Mitigation

1. If the operator log shows a bootstrap peer failure, verify the peer
   secret/token referenced in `spec.peers` is still valid and the peer
   cluster is reachable.
2. If it shows a cluster-info or Ceph version detection failure, resolve
   the underlying Ceph cluster health issue first.
3. If it shows a deployment-creation failure, check for resource
   constraints (node selectors, taints, insufficient CPU/memory) preventing
   the rbd-mirror pod from scheduling.
4. Once the root cause is resolved, the resource reconciles back to
   `Ready` automatically; no manual phase reset is needed.
5. If the operator appears stuck despite a healthy cluster, restart it and
   re-check:
   ```bash
   oc rollout restart deployment rook-ceph-operator -n openshift-storage
   ```

If the issue persists, [gather logs](helpers/gatherLogs.md) for support.
