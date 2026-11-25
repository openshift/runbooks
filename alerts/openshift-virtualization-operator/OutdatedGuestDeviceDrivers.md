# OutdatedGuestDeviceDrivers

## Meaning

This alert fires when outdated VirtIO drivers are detected in a Windows virtual
machine (VM) for more than 1 hour.

The QEMU guest agent reports the name and version of the driver of each VirtIO
device in the guest. A driver is outdated when its build number, the last
component of its version, is lower than the minimum build configured for the
cluster for that device and Windows version. Equal and newer builds do not
alert.

The alert fires once per VM and device.

## Impact

The VM continues to run, but its drivers might lack performance, stability and
security fixes that are included in the minimum driver version configured for
the cluster.

The cluster and other VMs are not affected.

## Diagnosis

Run the following PromQL queries, replacing the placeholders with the values of
the alert labels.

**Note:** You can run the queries in the web console by navigating to
**Observe** -> **Metrics**.

1. List all outdated drivers of the VM:

   ```promql
   vmi:kubevirt_vmi_guest_device_driver_build_number:max{namespace="<namespace>", name="<name>"}
   ```

   Each result is one outdated driver. The `driver_name` and `driver_version`
   labels show the installed driver, and the value is its build number.

2. Obtain the minimum driver version configured for the cluster for a device
   and Windows version:

   ```promql
   kubevirt_guest_device_driver_minimum_build_number{device_id="<device_id>", guest_os_version_id="<guest_os_version_id>"}
   ```

   The `driver_version` label shows the minimum version, and the value is its
   build number.

If the alert fires for many VMs at the same time, the minimum driver versions
configured for the cluster have most likely been raised, for example after a
cluster upgrade. In this case, all affected VMs need a driver update.

## Mitigation

The VM owner must update all VirtIO drivers in the VM to the latest drivers
provided by the cluster. The cluster provides the drivers as a container disk.

Update the drivers by following
[Update VirtIO drivers and the guest agent on a Windows
VM](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/virtualization/managing-vms#virt-updating-virtio-drivers-guest-agent-windows_virt-update-virtio-drivers).

The alert resolves after the guest agent reports the updated driver versions.

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.