# StorageClientCertNotConfigured

## Meaning

mTLS client-certificate material required for authentication
is missing from the spec of a `StorageClient` on the client cluster.

The `cert_type` label identifies the missing material:

- `server_ca` — `spec.serverCAConfigMap` is empty. No ConfigMap holding the peer
  server CA certificate is referenced, so the client cannot verify the storage
  provider's server certificate.
- `client_cert` — `spec.clientCertSecret` is empty. No secret holding the client
  certificate is referenced, so the client cannot present a certificate to the
  storage provider.

## Impact

The connection from the storage client to the storage provider is not secured
with mutual TLS. The client cannot authenticate itself to the provider with a
certificate, so the connection runs without client-certificate authentication
and without per-client cryptographic isolation.

## Diagnosis

Identify the affected client and missing material from the alert labels
`storage_client` and `cert_type`.

On the consumer cluster, inspect the StorageClient spec:

```bash
oc get storageclient <storage_client> -n openshift-storage -o yaml
```

Check the field matching `cert_type`:

- `server_ca` — `spec.serverCAConfigMap` must select a key in a ConfigMap holding
  the PEM-encoded peer server CA certificate.
- `client_cert` — `spec.clientCertSecret.name` must reference a
  `kubernetes.io/tls` secret (keys `tls.crt`, `tls.key`) with the client
  certificate presented to the provider.

Confirm the metric and which material is flagged:

```bash
# 1 = not configured, 0 = configured
ocs_client_operator_storage_client_cert_status{storage_client="<name>"}
```

If a field is set, verify the referenced object exists:

```bash
oc get configmap <serverCAConfigMap-name> -n openshift-storage   # cert_type=server_ca
oc get secret <clientCertSecret-name> -n openshift-storage        # cert_type=client_cert
```

## Mitigation

Populate the missing material in the StorageClient spec.

For a missing server CA (`cert_type=server_ca`):

1. Create (or locate) a ConfigMap in `openshift-storage` holding the PEM-encoded
   peer server CA certificate.
2. Set `spec.serverCAConfigMap` to select that ConfigMap and key.

For a missing client certificate (`cert_type=client_cert`):

1. Create (or locate) a `kubernetes.io/tls` secret in `openshift-storage`
   (keys `tls.crt`, `tls.key`) containing the client certificate. It must be
   signed by a CA trusted by the provider and carry the SAN the provider expects
   (see the provider-side
   [StorageConsumerClientCertNotConfigured](StorageConsumerClientCertNotConfigured.md)).
2. Set `spec.clientCertSecret.name` to that secret.

```bash
oc edit storageclient <storage_client> -n openshift-storage
```

Once the flagged field is set,the alert clears after the next evaluation.
If mTLS is intentionally not used for this client,the alert can be silenced per
site policy.

See also: [StorageConsumerClientCertNotConfigured](StorageConsumerClientCertNotConfigured.md)—
the provider-cluster side of the same mTLS configuration.
