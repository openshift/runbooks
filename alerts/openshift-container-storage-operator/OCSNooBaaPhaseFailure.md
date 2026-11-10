# OCSNooBaaPhaseFailure

## Meaning

A `NooBaa` custom resource has been in the `Rejected` phase for more than
10 minutes. NooBaa is the Multi-Cloud Gateway (MCG) backing ODF's S3
object storage interface. `Rejected` means the NooBaa operator's reconcile
never got past its first "Verifying" phase, because the CR itself is
invalid — this is a CR-authoring/day-0 configuration error, not a
transient runtime failure, and is **not retried on a timer**.

The operator only rejects the CR for one of these reasons:

- The `NooBaa` CR's name does not match the single, fixed system name the
  operator expects (normally `noobaa`, one per namespace).
- `spec.image` does not parse as a valid container image reference.
- `spec.image` uses the official NooBaa core image name, but its tag's
  version does not satisfy the operator's supported version range.

The exact reason and a human-readable message are recorded as a `Warning`
Kubernetes `Event` on the `NooBaa` object and in its `status.conditions`.

## Impact

Because `Rejected` happens before the Creating/Connecting/Configuring
phases, the reconciler is not managing the NooBaa core `StatefulSet`,
mgmt/S3 services, default `BackingStore`/`BucketClass`, or the OBC storage
class. For a brand-new/initial `NooBaa` CR, the S3/MCG object storage
service never comes up. For a previously-`Ready` system whose CR was
later edited into an invalid state, the existing core pod and its running
S3 endpoint are not torn down, so already-provisioned buckets/objects
typically keep serving, but no further reconciliation (scaling, config
drift correction, upgrades) happens until the spec is fixed.

## Diagnosis

### Check the NooBaa status, conditions, and events

```bash
oc get noobaa -n openshift-storage
oc get noobaa <noobaa_name> -n openshift-storage -o yaml
oc describe noobaa <noobaa_name> -n openshift-storage
oc get events -n openshift-storage --field-selector involvedObject.name=<noobaa_name>
```

The `status.conditions` message and the `Warning` event will state exactly
which check failed: `InvalidSystemName`, `InvalidImage`, or
`InvalidImageVersion`.

### Confirm the CR name

The `NooBaa` CR must use the operator's expected system name (typically
`noobaa`, one per namespace). A renamed or duplicate CR will be rejected
with `InvalidSystemName`.

### Confirm the image reference

```bash
oc get noobaa <noobaa_name> -n openshift-storage -o jsonpath='{.spec.image}'
```

Verify it is a well-formed image reference, and if it is the official
NooBaa core image, that its tag falls within the version range this
operator release supports.

## Mitigation

1. Edit the `NooBaa` CR (directly, or via the `StorageCluster` CR that
   manages it) to fix whichever condition was reported:
   - Correct the CR name to the expected singleton name, or
   - Fix `spec.image` to a valid, parseable image reference, or
   - Use a supported image version/tag.
2. No restart or manual phase reset is needed — the operator watches the
   `NooBaa` CR and re-reconciles automatically as soon as the spec is
   corrected, progressing through Verifying/Creating/Connecting/
   Configuring back to `Ready`.
3. If `oc describe`/`oc get events` do not show one of the three reasons
   above, the operator version may differ from what this runbook assumes;
   check the NooBaa operator logs directly:
   ```bash
   oc logs -n openshift-storage deploy/noobaa-operator
   ```

If the issue persists, [gather logs](helpers/gatherLogs.md) for support.
