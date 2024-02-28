# MDSCpuUsageHigh

## Meaning

Ceph MDS cpu usage for the MDS daemon has exceeded the threshold for adequate
performance.

## Impact

MDS serves filesystem metadata. The MDS is crucial for any file creation,
rename, deletion and update operations.
MDS is by default is allocated 2 or 3 CPUS.
It is okay as long as metadata operations are not too many.
When the metadata operation load increases enough to trigger this alert,
it means the default CPU allocation is unable to cope with load and we need to
increase the CPU allocation.

## Diagnosis

To diagnose the alert, click on the workloads->pods and select the
corresponding MDS pod and click on the metrics tab.
You should be able to see the allocated and used CPU. By default,
the alert is fired if the used CPU is 67% of allocated CPU for 6 hours.
If this is the case take the steps mentioned in mitigation.

## Mitigation

We need to increase the allocated CPU. The below command describes
how to set the number of allocated CPU for MDS server.

```bash
oc patch -n openshift-storage storagecluster ocs-storagecluster \
    --type merge \
    --patch '{"spec": {"resources": {"mds": {"limits": {"cpu": "8"},
    "requests": {"cpu": "8"}}}}}'
```
