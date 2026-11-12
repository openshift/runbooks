# StorageClientCertNotConfigured

This alert flags a weakness in the management (control) path between the provider
and consumer clusters. The data path is unaffected and no I/O is stopped.

## Meaning

mTLS client-certificate material required for authentication
is missing from the spec of a `StorageClient` on the client cluster.

The `cert_type` label identifies the missing material:

- `server_ca` — `spec.serverCAConfigMap` is empty. No ConfigMap holding the
  provider's server CA certificate is referenced, so the client cannot verify the
  storage provider's server certificate.
- `client_cert` — `spec.clientCertSecret` is empty. No secret holding the client
  certificate is referenced, so the client cannot present a certificate to the
  storage provider.

## Impact

The connection from the storage client to the storage provider is not secured
with mutual TLS. The client cannot authenticate itself to the provider with a
certificate and cannot verify the server's authenticity.

## Diagnosis

Identify the affected client and missing material from the alert labels
`storage_client` and `cert_type`.

On the consumer cluster, inspect the StorageClient spec:

```bash
oc get storageclient <storage_client> -n openshift-storage -o yaml
```

Check the field matching `cert_type`:

- `server_ca` — `spec.serverCAConfigMap` must select a key in a ConfigMap holding
  the PEM-encoded provider's server CA certificate.
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
   provider's server CA certificate. This CA is usually provided by the storage
   admin of the provider/hub cluster.
2. Set `spec.serverCAConfigMap` to select that ConfigMap and key.

For a missing client certificate (`cert_type=client_cert`):

1. Create (or locate) a `kubernetes.io/tls` secret in `openshift-storage`
   (keys `tls.crt`, `tls.key`) containing the client certificate. It must be
   signed by the CA that the corresponding StorageConsumer on the provider cluster
   references via `spec.clientCAConfigMap`, and carry the SAN set in that
   StorageConsumer's `spec.clientSAN` (see the provider-side
   [StorageConsumerClientCertNotConfigured](StorageConsumerClientCertNotConfigured.md)).
   This client certificate is usually created by your IT for this consumer cluster.
2. Set `spec.clientCertSecret.name` to that secret.

```bash
oc edit storageclient <storage_client> -n openshift-storage
```

Once the flagged field is set, the alert clears after the next evaluation.
If mTLS is intentionally not used for this client, the alert can be silenced per
site policy. Note that the future version will mandate mTLS for security and will
block upgrades until it is configured.

See also: [StorageConsumerClientCertNotConfigured](StorageConsumerClientCertNotConfigured.md)—
the provider-cluster side of the same mTLS configuration.
