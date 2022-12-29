# LowKVMNodesCount

## Meaning

This alert fires when fewer than two nodes in the cluster have KVM resources.

## Impact

The cluster must have at least two nodes with KVM resources for live migration.

Virtual machines cannot be scheduled or run if no nodes have KVM resources.

## Diagnosis

- Identify the nodes with KVM resources:

  ```bash
  $ oc get nodes -o jsonpath='{.items[*].status.allocatable}' | \
    grep devices.kubevirt.io/kvm
  ```

## Mitigation

Install KVM on the nodes without KVM resources.