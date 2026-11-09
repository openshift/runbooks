# MultipathPathDegraded

## Meaning

For 5 minutes, at least one path for a DM-multipath device has been outside
the active state at each evaluation. The failed paths can change between
evaluations. The metrics come from the node-exporter `dmmultipath` collector.
When that collector is unavailable or disabled, or no multipath devices
exist, this alert does not fire.

## Impact

Path redundancy is degraded. I/O can continue through usable paths, but
further failures may disrupt access to the device. If all paths are already
failed, applications and VMs using the device may have blocked I/O.

`MultipathPathCritical` can fire alongside this alert.

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

Restore the failed storage paths with the storage administrator. Check cable
and transceiver health, zoning, HBA ports, storage sessions, and array ports.
Before closing the incident, recheck `multipath -ll` on the affected node and
confirm that affected workload I/O is healthy: the failed-path count and the
alert can clear while multipath still marks a path unusable. Do not remove
the mapped device while workloads are using it.