# VirtualMachineInstanceHasEphemeralHotplugVolume

## Meaning

The `VirtualMachineInstanceHasEphemeralHotplugVolume` alert is triggered when a
virtual machine instance (VMI) contains an ephemeral hotplug volume. An
ephemeral hotplug volume only exists in the VMI and does not persist when
restarting a virtual machine (VM).

## Impact

The `HotplugVolumes` feature gate will be deprecated in a future release and
will
be replaced by the `DeclarativeHotplugVolumes` feature gate. The two are
mutually
exclusive, and when `DeclarativeHotplugVolumes` becomes enabled, any remaining
ephemeral hotplug volumes will be automatically unplugged from all VMIs.

This alert is triggered to inform users about the future deprecation and to
suggest steps to convert ephermeral volumes to persistent ones.

## Diagnosis

1. Find each VM that contains an ephemeral hotplug volume.
   This command returns a list with entries in the [vm-name, namespace] format.
    ``` bash
    $ oc get vmis -A -o json | jq -r '.items[].metadata | select(.annotations |
has("kubevirt.io/ephemeral-hotplug-volumes")) | [.name , .namespace] | @tsv'
    ```
2. For each VM listed, find the volumes that need to be patched.
   ``` bash
   $ oc get vmis <vm-name> -n <namespace> -o json | jq -r
'.metadata.annotations."kubevirt.io/ephemeral-hotplug-volumes"'
   ```

## Mitigation

To mitigate the impact of this alert, consider converting the ephemeral
hotplug volumes in the VM to persistent volumes instead.

To convert ephemeral volumes to persistent volumes, run the following command:
``` bash
$ virtctl addvolume <vm-name> --volume-name=<volume-name> --persist
```

If you cannot resolve the issue, log in to the [Red Hat Customer
Portal](link:https://access.redhat.com)
and open a support case, attaching the artifacts gathered during the diagnosis
procedure.