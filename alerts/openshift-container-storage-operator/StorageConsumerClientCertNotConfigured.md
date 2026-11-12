# StorageConsumerClientCertNotConfigured

## Meaning

Fields requied for mTLS authentication of client are missing from the spec of a
`StorageConsumer` on the storage provider cluster.

The `cert_type` label identifies which piece of material is missing:

- `client_ca` — `spec.clientCAConfigMap` is empty. No CA ConfigMap is referenced
  to validate the consumer's client certificate.
- `client_san` — `spec.clientSAN` is empty. No expected Subject Alternative Name
  is set to pin the consumer's identity.

## Impact

Client-certificate (mTLS) authentication is not in place for the affected
consumer. The provider cannot cryptographically validate and isolate that
consumer's connection, so the connection falls back to a weaker trust posture.
Per-consumer isolation and SAN pinning are not enforced until both the CA and
the SAN are configured.

## Diagnosis

Identify the affected consumer and which material is missing from the alert
labels `storage_consumer_name` and `cert_type`.

On the provider cluster, inspect the StorageConsumer spec:

```bash
oc get storageconsumer <storage_consumer_name> -n openshift-storage -o yaml
```

Check these fields:

- `spec.clientCAConfigMap.name` — must reference a ConfigMap holding the CA that
  signs the consumer's client certificate.
- `spec.clientSAN` — must hold the SAN expected in the consumer's client
  certificate.

Confirm the metric and which material is flagged:

```bash
# 1 = not configured, 0 = configured
ocs_storage_consumer_cert_status{storage_consumer_name="<name>"}
```

If `cert_type="client_ca"`, verify the referenced ConfigMap exists:

```bash
oc get configmap <name> -n openshift-storage
```

## Mitigation

Populate the missing material in the StorageConsumer spec.

For a missing CA (`cert_type=client_ca`):

1. Create (or locate) a ConfigMap in `openshift-storage` containing the CA that
   signs this consumer's client certificate.
2. Set `spec.clientCAConfigMap.name` to that ConfigMap.

For a missing SAN (`cert_type=client_san`):

1. Set `spec.clientSAN` to the SAN embedded in the consumer's client
   certificate.

```bash
oc edit storageconsumer <storage_consumer_name> -n openshift-storage
```

Once both `clientCAConfigMap` and `clientSAN` are set,the alert clears
after the next evaluation. If mTLS is intentionally not used for this
consumer, the alert can be silenced per site policy.

See also: [StorageClientCertNotConfigured](StorageClientCertNotConfigured.md) —
the consumer-cluster side of the same mTLS configuration.
