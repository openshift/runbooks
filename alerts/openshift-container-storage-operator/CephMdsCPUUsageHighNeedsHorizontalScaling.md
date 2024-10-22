# MDSCpuUsageHighNeedsHorizontalScaling

## Meaning

Ceph MDS CPU usage for the MDS daemon has exceeded the threshold for adequate
performance.

## Impact

MDS serves filesystem metadata. The MDS is crucial for any file creation,
rename, deletion and update operations.
MDS is by default is allocated 2 or 3 CPUS.
It is okay as long as metadata operations are not too many.
As metadata server get loaded due to high number of requests to the server,
the metadata operation load increases and eventually triggers this alert.
It means we need to scale horizontally by adding more metadata servers.
This will help parallely serve the client requests much efficiently without
overwhelming any single MDS pod.

## Diagnosis

To diagnose the alert, click on the workloads->pods and select the
corresponding MDS pod and click on the metrics tab.
You should be able to see the allocated and used CPU. By default,
the alert is fired if the used CPU is 67% of allocated CPU and there
is high rate of mds requests for past 6 hours.
If this is the case take the steps mentioned in mitigation.

## Mitigation

In this case, we need to add more active metadata servers. The below
command describes how to add multiple active MDS servers,

```bash
oc patch -n openshift-storage storagecluster ocs-storagecluster\
    --type merge \
    --patch '{"spec": {"managedResources": {"cephFilesystems":{"activeMetadataServers": 2}}}}'
```
PS: Make sure we have enough CPU provisioned for MDS depending on the load.

Always increase the `activeMetadataServers` by 1 and analyze the load.
The scaling of activeMetadataServers work only if you have more than one PV.
If there is only one PV that is causing CPU load, look at increasing the
CPU resource as described in
[VerticalScaling](CephMdsCPUUsageHighNeedsVerticalScaling.md) file
