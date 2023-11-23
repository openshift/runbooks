# CephOSDCriticallyFull

## Meaning

Utilization of back-end storage device (OSD) has crossed 80%.
Immediately free up some space or expand the storage cluster or contact support.

## Impact

One of the OSD Storage size has crossed 80% of the total capacity. Expand the
cluster immediately.

## Diagnosis

Alert message will have enough information about the underlying failure.
It should show the name of 'ceph-daemon', 'device-class' and the 'host-name'.

A sample alert message is provided below,

    Utilization of storage device <ceph-daemon-name> of device_class type
<device-class-name> has crossed 80% on host <host-name>. Immediately free up
some space or add capacity of type <device-class>.

## Mitigation

### Delete data

The customer may delete data and the cluster will resolve the alert through self
healing processes.

### Expand the storage capacity

Customer may assess their ability to expand. Here are some points,

**Current Storage Size < 1TB**:  
The customer may increase capacity via the addon and the cluster will resolve
the alert through self healing processes.

**Current Size itself is 1TB**:
Please contact your dedicated customer care support.

[gather_logs](helpers/gatherLogs.md)

