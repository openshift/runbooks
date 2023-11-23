# CephMgrIsMissingReplicas

## Meaning

Ceph Manager is missing replicas. That means Storage metrics collector
service doesn't have required no of replicas.

## Impact

This impacts the health status reporting and will cause some of the information
reported by `ceph status` to be missing or stale. In addition, the ceph manager
is responsible for a Manager framework aimed at expanding the Ceph existing capabilities.

## Diagnosis

To resolve this alert, you will need to determine the cause of the disappearance
of the Ceph Manager through the logs and restart if necessary.

## Mitigation

Check the manager pod's logs. Verify the rook-ceph-mgr pod is failing and restart
if necessary. If the ceph mgr pod restart fails, use general basic pod troubleshooting
to resolve.

[pod debug](helpers/podDebug.md) [gather_logs](helpers/gatherLogs.md)

