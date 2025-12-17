# VirtualMachineInstanceHasEphemeralHotplugVolume

## Meaning

The `VirtualMachineInstanceHasEphemeralHotplugVolume` alert is triggered when a
Virtual Machine Instance (VMI) contains an Ephemeral Hotplug Volume, which is
defined
as a hotplug volume that only exists in the VMI and will not persist during VM
restart

## Impact

The `HotplugVolumes` Feature Gate will be deprecated in a future release and
will
be replaced by the `DeclarativeHotplugVolumes` Feature Gate. The two are
mutually
exclusive, and when `DeclarativeHotplugVolumes` is enabled, any remaining
ephemeral hotplug volumes will automatically be unplugged from any VMIs.

If this alert is triggered, it is to inform of this future deprecation and to
suggest steps to convert ephermeral volumes to persist ones.

## Diagnosis

To diagnose the cause of this alert, the following steps can be taken:

1. Find each VM that contains an ephemeral hotplug volume.
   This command will print out list entries in format [vm-name, namespace]
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

To mitigate the impact of this alert, consider converting these ephemeral
hotplug volumes to instead persist within the VM.

To do so:
``` bash
$ virtctl addvolume <vm-name> --volume-name=<volume-name> --persist
```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.