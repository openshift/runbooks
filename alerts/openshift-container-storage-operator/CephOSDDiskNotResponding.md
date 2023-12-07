# CephOSDDiskNotResponding

## Meaning

A disk device on one of the hosts is not responding, potentially impacting the
performance and availability of the OSD.

## Impact

**Severity:** Error
**Potential Customer Impact:** Medium

This alert signals that a disk device is not responding on a host, and it may
affect the proper functioning of the associated OSD (Object Storage Daemon).

## Diagnosis

The alert is raised when a disk device is identified as not responding. To
diagnose the issue, check whether all OSDs are up and running.

**Prerequisites:** [Prerequisites](helpers/diagnosis.md)

## Mitigation

### Recommended Actions

1. **Check for Network Issues:** Verify if the unresponsive disk issue is
   related to a network problem by following the steps in the provided Standard
   Operating Procedure (SOP) -
   [Check Ceph Network Connectivity SOP](helpers/networkConnectivity.md).
   Escalate to the ODF team if it is a network issue.

2. **Generic Debugging:** Follow the general pod debug workflow outlined below
   to identify and address potential issues.
   - [Pod Debug Workflow](helpers/podDebug.md)
   - [Gather Logs](helpers/gatherLogs.md)

**Additional Resources:**

- [Troubleshooting](helpers/troubleshootCeph.md)
