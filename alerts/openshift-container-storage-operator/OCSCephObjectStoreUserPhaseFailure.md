# OCSCephObjectStoreUserPhaseFailure

## Meaning

A `CephObjectStoreUser` custom resource has been in a `Failed` or
`ReconcileFailed` phase for more than 10 minutes. `CephObjectStoreUser`
creates a Ceph RGW (S3) user and its access/secret key `Secret`, created/
managed by the `StorageCluster`/Rook operator or by a consuming
application.

The Rook object-user controller sets `ReconcileFailed` as a single
catch-all for any error while reconciling, including: the parent
`CephObjectStore` not yet being ready, a cluster-info/Ceph version lookup
failure, a referenced key-pair `Secret` being malformed, `spec` validation
failure (missing `spec.store`, or an account-linked user's `displayName`
not matching the IAM-compatible pattern `[\w+=,.@-]+`), an account
reference (`spec.accountRef`) that cannot be resolved, or the underlying
`radosgw-admin`/Ceph RGW user create/update call failing. The phase does
not distinguish which of these occurred; check the operator logs for the
specific cause.

## Impact

The S3 access/secret key `Secret` for this user
(`rook-ceph-object-user-<store>-<name>` in the `openshift-storage`
namespace) is not created or not kept up to date. Applications relying on
this user's credentials to access the S3-compatible endpoint can fail to
authenticate, or may be using stale credentials if the user previously
existed and is now failing to reconcile an update.

## Diagnosis

### Check the CephObjectStoreUser status and events

```bash
oc get cephobjectstoreuser -n openshift-storage
oc describe cephobjectstoreuser <user_name> -n openshift-storage
```

Check `status.phase`/`status.message` and the `Events` section for the
specific failure reason.

### Check the parent CephObjectStore is Ready

```bash
oc get cephobjectstore <store_name> -n openshift-storage -o jsonpath='{.status.phase}'
```

If the parent object store itself is unhealthy, the RGW endpoint it needs
to call to create/update the user will be unreachable.

### Check the generated Secret

```bash
oc get secret rook-ceph-object-user-<store_name>-<user_name> -n openshift-storage
```

If this `Secret` is missing or stale, the user was never successfully
created/updated in Ceph.

### Check the Rook operator logs

```bash
oc logs -n openshift-storage deploy/rook-ceph-operator | grep -i <user_name>
```

Look for `invalid CephObjectStoreUser ... spec` (a spec validation
failure), `failed to generate UserKeySpec` (a referenced key `Secret`
problem), or an error from the RGW admin API call itself.

### Check overall Ceph cluster health

See [Access the Toolbox](helpers/accessToolbox.md) and
[Troubleshoot Ceph](helpers/troubleshootCeph.md) to check `ceph status`
and `ceph health detail`.

## Mitigation

1. If the parent `CephObjectStore` isn't `Ready`, resolve that first.
2. If the operator log shows a spec validation error, correct the
   `CephObjectStoreUser` spec (for example, set `spec.store`, or fix
   `spec.displayName` to be IAM-compatible when `spec.accountRef` is set).
3. If a referenced key-pair `Secret` is malformed, correct or remove the
   `spec.keys` reference.
4. Once the root cause is resolved, the resource reconciles back to
   `Ready` automatically; no manual phase reset is needed.
5. If the operator appears stuck despite a healthy cluster, restart it and
   re-check:
   ```bash
   oc rollout restart deployment rook-ceph-operator -n openshift-storage
   ```

If the issue persists, [gather logs](helpers/gatherLogs.md) for support.
