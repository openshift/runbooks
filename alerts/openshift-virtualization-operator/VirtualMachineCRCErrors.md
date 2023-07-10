# VirtualMachineCRCErrors

## Meaning

This alert fires when the storage class is incorrectly configured.
A system-wide, shared dummy page causes CRC errors when data is
written and read across different processes or threads.

## Impact

A large number of CRC errors might cause the cluster to display
severe performance degradation.

## Diagnosis

1. Get the volume name from a virtual machine:

   ```bash
   $ oc get vm <vm_name> -o jsonpath='{.spec.template.spec.volumes}'
   ```

2. Get the storage class name from the volume:

   ```bash
   $ oc get pvc <volume> -o jsonpath='{.spec.storageClassName}'
   ```

3. Get the storage class configuration:

   ```bash
   $ oc get sc <storage_class> -o yaml
   ```

4. Check the storage class configuration for the `krbd:rxbounce` map option.

## Mitigation

Add the `krbd:rxbounce` map option to the storage class configuration:

```bash
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: example-storage-class
parameters:
  # ...
  mounter: rbd
  mapOptions: "krbd:rxbounce"
provisioner: openshift-storage.rbd.csi.ceph.com
# ...
```

The `krbd:rxbounce` option creates a bounce buffer to receive data. The default
behavior is for the destination buffer to receive data directly. A bounce buffer
is required if the destination buffer is unstable.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.

