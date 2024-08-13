# VMStorageClassWarning

## Meaning

When running VMs using ODF storage with 'rbd' mounter or 'rbd.csi.ceph.com'
provisioner, Windows VMs may cause reports of bad crc/signature errors due to
certain I/O patterns. Cluster performance can be severely degraded if the number
of re-transmissions due to crc errors causes network saturation. 'krbd:rxbounce'
should be configured for the VM storage class to prevent these crc errors, the
"ocs-storagecluster-ceph-rbd-virtualization" storage class uses this option by
default, if available.

## Impact

Cluster may report a huge number of CRC errors and the cluster might experience
major service outages.

## Diagnosis

Obtain a list of VirtualMachines with an incorrectly configured storage class by
running the following PromQL query:

**Note:** You can use the Openshift metrics explorer available at
'https://{OPENSHIFT_BASE_URL}/monitoring/query-browser'.

```promql
$ kubevirt_ssp_vm_rbd_block_volume_without_rxbounce == 1
```

Example output:

```plaintext
kubevirt_ssp_vm_rbd_block_volume_without_rxbounce{name="testvmi-gwgdqp22k7", namespace="test_ns_1"} 1
kubevirt_ssp_vm_rbd_block_volume_without_rxbounce{name="testvmi-rjqbjg47a8", namespace="test_ns_2"} 1
```

The output displays a list of VirtualMachines that use a storage class without
`rxbounce_enabled`.

Obtain the VM volumes by running the following command:

```bash
$ oc get vm <vm-name> -o json | jq -r '.spec.template.spec.volumes[] | if .dataVolume then "DataVolume - " + .dataVolume.name elif .persistentVolumeClaim then "PersistentVolumeClaim - " + .persistentVolumeClaim.claimName else empty end'
```

## Mitigation

It is recommended to create a dedicated StorageClass with "krbd:rxbounce" map
option for the disks of virtual machines, to use a bounce buffer when receiving
data. The default behavior is to read directly into the destination buffer. A
bounce buffer is required if the stability of the destination buffer cannot be
guaranteed.

Note that changing the used storage class will not have any impact on existing
PVCs/VMs, meaning that new VMs will be provisioned with the optimized settings
but existing VMs need to be transitioned or the alert will continue to fire for
those.

```bash
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  name: vm-sc
parameters:
  # ...
  mounter: rbd
  mapOptions: "krbd:rxbounce"
provisioner: openshift-storage.rbd.csi.ceph.com
# ...
```

See [Optimizing ODF PersistentVolumes for Windows VMs](https://access.redhat.com/articles/6978371)
for details.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.