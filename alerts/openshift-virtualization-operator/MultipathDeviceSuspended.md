# MultipathDeviceSuspended

## Meaning

A DM-multipath device has remained suspended for 1 minute. The node-exporter
`dmmultipath` collector reports `node_dmmultipath_device_active` as `0` for a
suspended device and `1` for an active device.

## Impact

I/O through the device is blocked. Applications and VMs using it may hang or
experience I/O timeouts. This alert does not establish why the device was
suspended.

## Diagnosis

Use the alert's `instance`, `device`, and `sysfs_name` labels to identify the
node and device. Confirm the node corresponding to the node-exporter scrape
target. On that node, inspect the device and multipath status:

```bash
dmsetup info <device>
dmsetup status <device>
multipath -ll <device>
journalctl -u multipathd
```

Check whether a storage management operation intentionally suspended the
device, and inspect underlying paths, storage connectivity, and system logs.
Coordinate with the administrator or component that manages the device.

Check the device state in Prometheus:

```promql
node_dmmultipath_device_active{job="node-exporter",instance="<instance>",device="<device>",sysfs_name="<sysfs_name>"}
```

## Mitigation

Resolve the underlying storage or management issue before resuming the device.
Coordinate restoration with the device's managing component or administrator;
a suspended device can be part of an in-progress storage operation.

Verify that the device is active, workload I/O resumes, and
`node_dmmultipath_device_active` returns to `1`. Do not remove or recreate the
device while workloads are using it.