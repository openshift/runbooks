# MultipathPathCritical

## Meaning

For 2 minutes, at least one path and at least half of the paths for a
DM-multipath device have been outside the active state at each evaluation.
The failed paths can change between evaluations. This includes total path
loss and failure of a device with only one path. The metrics come from the
node-exporter `dmmultipath` collector.

## Impact

The device has severely reduced path redundancy or has already lost all
paths. Applications and VMs using it are at risk of I/O disruption or blocked
I/O. `MultipathPathDegraded` can fire alongside this alert.

## Diagnosis

Use the alert's `instance`, `device`, and `sysfs_name` labels to identify the
node and DM-multipath device. The instance label identifies the node-exporter
scrape target; confirm the corresponding node before inspecting it.

On the affected node, inspect the device and its underlying paths:

```bash
multipath -ll <device>
dmsetup info <device>
dmsetup status <device>
journalctl -u multipathd
```

Check storage fabric connectivity, HBA and switch ports, array controllers,
and iSCSI sessions where applicable. Correlate with storage maintenance and
recent configuration changes.

In Prometheus, compare the failed and total path metrics for the same instance,
device, and sysfs name. Run each query separately:

```promql
node_dmmultipath_device_paths_failed{job="node-exporter",instance="<instance>",device="<device>",sysfs_name="<sysfs_name>"}
```

```promql
node_dmmultipath_device_paths{job="node-exporter",instance="<instance>",device="<device>",sysfs_name="<sysfs_name>"}
```

If the paths recover before you run these queries, the instant results show
zero failed paths and do not establish the earlier condition. To inspect the
incident window, graph both selectors in the Prometheus console over the
alert's firing period and compare failed and total paths over time.

The path metrics classify SCSI `running` and NVMe `live` as active. They report
underlying device state, not multipathd path-checker or ALUA state. Missing
metrics do not establish that storage is healthy.

## Mitigation

Immediately work with the storage administrator to restore connectivity.
Prioritize devices with no usable paths. Investigate fabric, HBA, switch,
storage session, and array controller failures. Verify path recovery and
affected workload I/O before closing the incident. Do not remove a mapped
device while workloads are using it.

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.