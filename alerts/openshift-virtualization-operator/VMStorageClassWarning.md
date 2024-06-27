# VMStorageClassWarning

## Meaning

This alert fires when the storage class is incorrectly configured.
A system-wide, shared dummy page causes CRC errors when data is
written and read across different processes or threads.

## Impact

A large number of CRC errors might cause the cluster to display
severe performance degradation.

## Diagnosis

1. Navigate to **Observe** -> **Metrics** in the web console.

2. Obtain a list of virtual machines with incorrectly configured storage classes
   by running the following PromQL query:
   ```text
   kubevirt_ssp_vm_rbd_volume{rxbounce_enabled="false", volume_mode="Block"} == 1
   ```

   The output displays a list of virtual machines that use a storage
   class without `rxbounce_enabled`.

   Example output:
   ```text
   kubevirt_ssp_vm_rbd_volume{name="testvmi-gwgdqp22k7", namespace="test_ns", pv_name="testvmi-gwgdqp22k7", rxbounce_enabled="false", volume_mode="Block"} 1
   ```

3. Obtain the storage class name by running the following command:

   ```bash
   $ oc get pv <pv_name> -o=jsonpath='{.spec.storageClassName}'
   ```

## Mitigation

Create a default OpenShift Virtualization storage class with the `krbd:rxbounce` map option. See [Optimizing ODF PersistentVolumes for Windows VMs](https://access.redhat.com/articles/6978371) for details.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
