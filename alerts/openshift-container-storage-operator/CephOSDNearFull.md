# CephOSDNearFull

## Meaning

Utilization of the back-end storage device OSD has crossed 75% on host
`<hostname>`. Free up some space or expand the storage cluster or contact
support.

## Impact

- **Severity:** Warning
- **Potential Customer Impact:** High

The OSD storage devices nearing full capacity can impact the overall performance
and availability of the Ceph storage system.

## Diagnosis

The alert is triggered when the utilization of the back-end storage device OSD
exceeds 75%. Detailed diagnosis involves checking whether all OSDs are up and
running.

### Prerequisites

1. Verify cluster access:
   - Check the output to ensure you are in the correct context for the cluster
     mentioned in the alert.
   - List clusters you have permission to access:

     ```bash
     ocm list clusters
     ```

   - From the list, find the cluster ID of the mentioned cluster.

2. Check Alerts:
   - Get the route to this cluster’s alert manager:

     ```bash
     MYALERTMANAGER=$(oc -n openshift-monitoring get routes/alertmanager-main
     --no-headers | awk '{print $2}')
     ```

   - Check all alerts:

     ```bash
     curl -k -H "Authorization: Bearer $(oc -n openshift-monitoring sa get-token
     prometheus-k8s)" https://${MYALERTMANAGER}/api/v1/alerts | jq '.data[] |
     select( .labels.alertname) | { ALERT: .labels.alertname, STATE:
     .status.state}'
     ```

3. (Optional) Document OCS Ceph Cluster Health:
   - You may check OCS Ceph Cluster health using the rook-ceph toolbox:
     - Check and document ceph cluster health:

       ```bash
       TOOLS_POD=$(oc get pods -n openshift-storage -l app=rook-ceph-tools -o name)
       oc rsh -n openshift-storage $TOOLS_POD
       ceph status
       ceph osd status
       exit
       ```

## Mitigation

### 1. Delete Data

The following instructions only apply
to OCS clusters that are near or full but NOT in readonly mode. Readonly mode
would prevent any changes including deleting data (i.e. PVC/PV deletions).

Delete some data, and the cluster will resolve the alert through
self-healing processes.

### 2. Current size < 1 TB, Expand to 4 TB

The user may increase capacity via the addon, and the cluster will resolve
the alert through self-healing processes.

### 3. Current size = 4TB

Please contact Dedicated Support.

Document Ceph Cluster health check:
[gather_logs](helpers/gatherLogs.md)

```bash
oc adm must-gather --image=<must-gather-image-name>
```
