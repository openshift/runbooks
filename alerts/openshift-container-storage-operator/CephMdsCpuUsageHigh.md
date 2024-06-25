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
increase the CPU allocation or run multiple active metadata servers.

## Diagnosis

To diagnose the alert, click on the workloads->pods and select the
corresponding MDS pod and click on the metrics tab.
You should be able to see the allocated and used CPU. By default,
the alert is fired if the used CPU is 67% of allocated CPU for 6 hours.
If this is the case take the steps mentioned in mitigation.

## Mitigation

We need to either increase the allocated CPU or run multiple active MDS. The
below command describes how to set the number of allocated CPU for MDS server.

```bash
oc patch -n openshift-storage storagecluster ocs-storagecluster \
    --type merge \
    --patch '{"spec": {"resources": {"mds": {"limits": {"cpu": "8"},
    "requests": {"cpu": "8"}}}}}'
```

In order to run multiple active MDS servers, use below command

```bash
oc patch -n openshift-storage cephfilesystem ocs-storagecluster-cephfilesystem\
    --type merge \
    --patch '{"spec": {"metadataServer": {"activeCount": 2}}}'

Make sure we have enough CPU provisioned for MDS depending on the load.
```

Always increase the `activeCount` by 1. The scaling of activeCount works only
if you have more than one PV. If there is only one PV that is causing CPU load,
look at increasing the cpu resource as described above.
