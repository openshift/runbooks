# KubeVirtNoAvailableNodesToRunVMs

## Meaning

This alert fires when the node CPUs in the cluster do not support virtualization
or the virtualization extensions are not enabled.

## Impact

The nodes must support virtualization and the virtualization features must be
enabled in the BIOS to run virtual machines (VMs).

## Diagnosis

- Check the nodes for hardware virtualization support:

  ```bash
  $ oc get nodes -o json|jq '.items[]|{"name": .metadata.name, "kvm": .status.allocatable["devices.kubevirt.io/kvm"]}'
  ```

  Example output:

  ```text
  {
    "name": "shift-vwpsz-master-0",
    "kvm": null
  }
  {
    "name": "shift-vwpsz-master-1",
    "kvm": null
  }
  {
    "name": "shift-vwpsz-master-2",
    "kvm": null
  }
  {
    "name": "shift-vwpsz-worker-8bxkp",
    "kvm": "1k"
  }
  {
    "name": "shift-vwpsz-worker-ctgmc",
    "kvm": "1k"
  }
  {
    "name": "shift-vwpsz-worker-gl5zl",
    "kvm": "1k"
  }
  ```

  Nodes with `"kvm": null` or `"kvm": 0` do not support virtualization extensions.

  Nodes with `"kvm": "1k"` do support virtualization extensions

## Mitigation

Ensure that hardware and CPU virtualization extensions are enabled on all nodes
and that the nodes are correctly labeled.

See [OpenShift Virtualization reports no nodes are available, cannot start VMs](https://access.redhat.com/solutions/5106121)
for details.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case.

