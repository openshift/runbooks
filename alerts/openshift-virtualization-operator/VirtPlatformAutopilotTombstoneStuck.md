# VirtPlatformAutopilotTombstoneStuck

## Meaning

This alert fires when virt-platform-autopilot cannot remove a resource it is
trying to delete (a "tombstone") for more than 30 minutes, either because
deletion failed (`-1`) or because the resource lacks the required management
label (`-2`).

The alert is driven by the `kubevirt_autopilot_tombstone_status` gauge,
labelled by `kind`, `name`, and `namespace`:

- `1` = the resource still exists (deletion pending)
- `0` = the resource is deleted or the CRD is absent
- `-1` = deletion error (an API error, finalizer, webhook, or RBAC blocks it)
- `-2` = skipped due to a label mismatch (the resource lacks the
  `platform.kubevirt.io/managed-by=virt-platform-autopilot` label)

The alert triggers on `kubevirt_autopilot_tombstone_status < 0` held for `30m`.

## Impact

This is a `warning` alert. The impact depends on the firing condition:

- **Deletion error (`-1`):** an obsolete resource lingers in the cluster but is
  not actively used. Over time, accumulated tombstones cause configuration drift
  and can conflict with a resource later recreated under the same name.
- **Label mismatch (`-2`):** the resource may be an active user-owned resource
  that the operator refuses to delete as a safety mechanism. This is not stale
  cleanup: verify ownership before treating it as a tombstone.

## Diagnosis

1. Export the `NAMESPACE` environment variable:

   ```bash
   $ NAMESPACE="$(oc get deployment -A -l app=virt-platform-autopilot \
       -o jsonpath='{.items[*].metadata.namespace}')"
   $ [ "$(echo "$NAMESPACE" | wc -w)" -eq 1 ] || \
       { echo "Error: expected 1 deployment, found in: ${NAMESPACE:-none}"; exit 1; }
   $ export NAMESPACE
   ```

2. Read the `kind`, `name`, and `namespace` from the alert labels, and note
   the metric value: `-1` (deletion error) or `-2` (label mismatch).

3. For a **deletion error (`-1`)**, check what is blocking removal as stuck
   finalizers, an in-progress deletion, RBAC, and the operator logs:

   ```bash
   $ oc get <kind> <name> -n <namespace> \
       -o jsonpath='{.metadata.finalizers} {.metadata.deletionTimestamp}'
   $ oc auth can-i get <resource-type> -n <namespace> \
       --as system:serviceaccount:$NAMESPACE:virt-platform-autopilot
   $ oc auth can-i delete <resource-type> -n <namespace> \
       --as system:serviceaccount:$NAMESPACE:virt-platform-autopilot
   $ oc -n $NAMESPACE logs --tail=200 \
       -l app=virt-platform-autopilot | grep -i tombstone
   ```

   Replace `<resource-type>` with the lowercased, pluralized form of `<kind>`
   (e.g., `services`, `configmaps`). For CRDs, use the format shown by
   `oc api-resources` (e.g., `kubedeschedulers`, `machineconfigs`). Omit
   `-n <namespace>` for cluster-scoped resources like `KubeDescheduler` or
   `MachineConfig`.

4. For a **label mismatch (`-2`)**, confirm the resource lacks the management
   label. This is the safety check working as intended: the operator refuses to
   delete a resource it did not create. Omit `-n <namespace>` for cluster-scoped
   resources like `KubeDescheduler` or `MachineConfig`.

   ```bash
   $ oc get <kind> <name> -n <namespace> \
       -o jsonpath='{.metadata.labels.platform\.kubevirt\.io/managed-by}'
   ```

## Mitigation

- **Deletion error (`-1`):** clear the blocker. Remove a finalizer left behind
  by a dead controller, fix or bypass the blocking webhook, or restore the
  operator's delete permissions. The operator retries deletion on every
  reconcile, so no manual delete is needed once the blocker is gone:

  Inspect the finalizers first (see the Diagnosis step) and confirm the one
  you intend to clear is stale, left by a controller that no longer exists and
  whose cleanup is already done. Then remove only that finalizer by its index,
  leaving any others in place. Use a test operation to verify the expected
  finalizer is still at that index:

  ```bash
  $ oc patch <kind> <name> -n <namespace> --type=json \
      -p='[{"op":"test","path":"/metadata/finalizers/<index>","value":"<expected-finalizer>"},
          {"op":"remove","path":"/metadata/finalizers/<index>"}]'
  ```

  The patch will fail if the finalizer array changed since inspection — if
that
  happens, re-inspect and adjust the index. Removing a live finalizer skips the
  cleanup it guards, so only do this when you understand the consequences.

- **Label mismatch (`-2`):** if the resource should be removed by the operator,
  add the management label so the next reconcile deletes it. Otherwise, the
  resource is owned by someone else and must be handled manually. Omit
  `-n <namespace>` for cluster-scoped resources like `KubeDescheduler` or
  `MachineConfig`.

  ```bash
  $ oc label <kind> <name> -n <namespace> \
      platform.kubevirt.io/managed-by=virt-platform-autopilot
  ```

The `for: 30m` clause only delays firing; it does not delay resolution. The
alert clears on the next evaluation after the expression becomes false, that
is, once `kubevirt_autopilot_tombstone_status` is no longer negative (the
resource is deleted or the CRD is absent, `0`, or exists again, `1`).

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.