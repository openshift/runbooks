# CephDataRecoveryTakingTooLong

## Meaning

Data recovery processes in Ceph are taking an extended amount of time,
indicating potential issues with the recovery speed.

## Impact

**Severity:** Warning
**Potential Customer Impact:** High

This alert indicates that the ongoing data recovery in Ceph is progressing
slower than expected, which may have a significant impact on system performance.

## Diagnosis

The alert is triggered when the data recovery process in Ceph is identified as
taking too long. It is recommended to check the status of all OSDs to ensure
they are up and running, as slow recovery may be related to OSD issues.

**Prerequisites:** [Prerequisites](helpers/diagnosis.md)

## Mitigation

### Recommended Actions

1. **Check for Network Issues:**
   Verify if the extended data recovery time is due to network issues by
   following the steps in the provided Standard Operating Procedure (SOP) -
   [Check Ceph Network Connectivity SOP](helpers/networkConnectivity.md).

2. **Generic Debugging:**
   Follow the general pod debug workflow outlined below to identify and
   address potential issues.

   - [Pod Debug Workflow](helpers/podDebug.md)
   - [Gather Logs](helpers/gatherLogs.md)

3. **Investigate Disk Utilization:** Check for high disk utilization (near 100%)
   during rebalancing or client I/O using tools like dstat. Assess if IOPS
   limitations on disks contribute to slow recovery.

4. **Object Count Impact:** Due to a higher object count (310 million),
   replication efforts are extensive, impacting recovery time. Consider
   adjusting parameters like "osd max push objects" to potentially accelerate
   the process.

**Additional Resources:**

- [Troubleshooting](helpers/troubleshootCeph.md)
