# KubeMacPoolDuplicateMacsFound

## Meaning

This alert fires when `KubeMacPool` detects duplicate MAC addresses.

`KubeMacPool` is responsible for allocating MAC addresses and preventing MAC
address conflicts. When `KubeMacPool` starts, it scans the cluster for the MAC
addresses of virtual machines (VMs) in managed namespaces.

## Impact

Duplicate MAC addresses on the same LAN might cause network issues.

## Diagnosis

1. Obtain the namespace and the name of the `kubemacpool-mac-controller` pod:

   ```bash
   $ oc get pod -A -l control-plane=mac-controller-manager --no-headers \
     -o custom-columns=":metadata.namespace,:metadata.name"
   ```
2. Obtain the duplicate MAC addresses from the `kubemacpool-mac-controller`
logs:

   ```bash
   $ oc logs -n <namespace> <kubemacpool_mac_controller> | grep "already allocated"
   ```

   Example output:

   ```text
   mac address 02:00:ff:ff:ff:ff already allocated to vm/kubemacpool-test/testvm, br1,
   conflict with: vm/kubemacpool-test/testvm2, br1
   ```

## Mitigation

1. Update the VMs to remove the duplicate MAC addresses.
2. Restart the `kubemacpool-mac-controller` pod:

   ```bash
   $ oc delete pod -n <namespace> <kubemacpool_mac_controller>
   ```