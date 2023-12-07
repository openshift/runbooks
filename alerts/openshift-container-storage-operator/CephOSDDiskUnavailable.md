# CephOSDDiskUnavailable

## Meaning

A disk device on one of the hosts is inaccessible, leading to the corresponding
OSD being marked out by the Ceph cluster.

## Impact

**Severity:** Error
**Potential Customer Impact:** High

This alert indicates that a disk device is not accessible on one of the hosts,
resulting in the Ceph cluster marking the corresponding OSD as out. The alert is
triggered when a Ceph node fails to recover within 10 minutes.

## Diagnosis

The alert is raised when a disk device becomes inaccessible on a host, causing
the associated OSD to be marked out by the Ceph cluster. To determine which node
has failures, follow the procedure outlined in
[Determine Failed OCS Node](helpers/determineFailedOcsNode.md).

**Prerequisites:** [Prerequisites](helpers/diagnosis.md)

## Mitigation

### Recommended Actions

1. **Determine Failed OCS Node:** Follow the procedure in
   [Determine Failed OCS Node](helpers/determineFailedOcsNode.md) to identify
   the node with failures.

2. **Gather Logs:** Collect logs using the [Gather Logs](helpers/gatherLogs.md)
   procedure for further analysis.

3. **Check for Network Issues:** Verify if the issue is related to a network
   problem by following the steps in the provided Standard Operating Procedure
   (SOP) -
   [Check Ceph Network Connectivity SOP](helpers/networkConnectivity.md).
   Escalate to the ODF team if it is a network issue.
