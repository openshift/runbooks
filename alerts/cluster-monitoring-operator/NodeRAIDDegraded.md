# NodeRAIDDegraded

## Meaning

The `NodeRAIDDegraded` alert triggers when a node has a storage configuration
with a RAID array and that array is reporting a degraded state because of one or
more disk failures.

## Impact

The affected node might go offline at any moment if the RAID array fully fails
because of issues with disks.

## Diagnosis

Open a shell on the node and use standard Linux utilities to diagnose the
issue. Note that you might also need to install additional software in the
debug container:

```console
$ NODE_NAME='<value of instance label from alert>'

$ oc debug "node/$NODE_NAME"
$ cat /proc/mdstat
```

## Mitigation

See the Red Hat Enterprise Linux [documentation][1] to see mitigation steps for
failing RAID arrays.

[1]: https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/8/html/managing_storage_devices/managing-raid_managing-storage-devices
