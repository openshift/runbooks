# CephClusterWarningState

## Meaning

Storage cluster is in warning state for more than 10m.

The Storage cluster has been in a warning state for an unacceptable amount of
time. While the storage operations will continue to function in this state, it
is recommended to fix errors so that the cluster does not get into error state
impacting operations. Check for other alerts that would have triggered prior
to this one and troubleshoot those alerts first.

## Impact

Cluster services not available. Errors in operations retrieving/putting
information in the cluster

## Diagnosis

See [general diagnosis document](helpers/diagnosis.md)

## Mitigation

### Check if it is a Network Issue

Check if it is a [network issue](helpers/networkConnectivity.md)

### Make changes to solve alert

General troubleshooting will be required in order to determine the cause of this
 alert. This alert will trigger along with other (usually many other) alerts.
Please view and troubleshoot the other alerts first.

### Review pods

[pod debug](helpers/podDebug.md)

If the basic health of the running pods, node affinity and resource availability
on the nodes have been verified, run Ceph tools for status of the storage
components.

#### Troubleshoot Ceph

[Troubleshoot_ceph_err](helpers/troubleshootCeph.md) and
[gather_logs](helpers/gatherLogs.md) to provide more information to support
teams.
