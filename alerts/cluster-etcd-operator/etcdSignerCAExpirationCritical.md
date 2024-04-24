# etcdSignerCAExpirationCritical

## Meaning

This alert fires when the signer CA certificate of etcd (metrics or server)
is having one year left before expiration.

## Impact

When the etcd server signer certificate expires, your cluster will become
unavailable and very hard to recover.
If the metrics signer certificate expires, you will lose metrics and alerts
for etcd.

## Mitigation

Please follow the [OpenShift documentation here][manualRota] on how to manually
 rotate the signer certificates.

[manualRota]: https://docs.openshift.com/container-platform/4.16/security/certificate_types_descriptions/etcd-certificates.html

