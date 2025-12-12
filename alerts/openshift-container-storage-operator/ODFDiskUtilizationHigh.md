# ODFDiskUtilizationHigh

## Meaning

A Ceph OSD disk is >90% busy (as measured by %util from iostat
semantics via node_disk_io_time_seconds_total), indicating heavy I/O load.

## Impact

* Increased I/O latency for block/object/file clients.
* Reduced cluster throughput during peak workloads.
* Potential for “slow request” warnings in Ceph logs.

## Diagnosis

1. Identify node and device from alert labels.
2. Check disk model and type:
```bash
oc debug node/<node>
lsblk -d -o NAME,ROTA,MODEL
# Confirm it’s an expected OSD device (HDD/SSD/NVMe)
```
3. Monitor real-time I/O:
```bash
iostat -x 2 5
```
4. Correlate with Ceph:
```bash
ceph osd df tree  # check weight and reweight
ceph osd perf     # check commit/apply latency
```

## Mitigation

* Add more disks to the cluster enhance the performance.
* Move the workloads to another storage system.

