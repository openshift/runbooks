# CephOSDSlowOps

## Meaning

OSD (Object Storage Daemon) requests are taking an extended amount of time to
process, indicating potential performance issues.

## Impact

**Severity:** Warning
**Potential Customer Impact:** Medium

This alert suggests that OSDs are experiencing delays in processing requests,
potentially affecting the overall performance of the Ceph storage system.

## Diagnosis

The alert is triggered when OSD requests take longer to process than the time
defined by the `osd_op_complaint_time` parameter, which is set to 30 seconds by
default. To gather more information about the slow requests, access the OSD pod
terminal and issue the following commands:

```bash
ceph daemon osd.<id> ops
ceph daemon osd.<id> dump_historic_ops
```

Note: Replace `<id>` with the OSD number, which can be found in the pod name
(e.g., `rook-ceph-osd-0-5d86d4d8d4-zlqkx`, where `<0>` is the OSD ID).*

## Mitigation

### Recommended Actions

1. **Check Hardware/Infrastructure:** Investigate problems with the underlying
   hardware/infrastructure, such as disk drives, hosts, racks, or network
   switches. Use the Openshift monitoring console to find alerts/errors about
   cluster resources.

2. **Check Network Issues:** Verify if the slow OSD operations are related to
   network problems. Follow the steps in the provided Standard Operating
   Procedure (SOP) -
   [Check Ceph Network Connectivity SOP](helpers/networkConnectivity.md).
   Escalate to the ODF team if it is a network issue.

3. **Review System Load:** Use the Openshift console to review metrics of the
   OSD pod and the Node running the OSD. If needed, add/assign more resources to
   address system load issues.
