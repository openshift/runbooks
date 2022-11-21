# NodeFilesystemSpaceFillingUp

## Meaning

The `NodeFilesystemSpaceFillingUp` alert triggers when two conditions are met:  

* The current file system usage exceeds a certain threshold.
* An extrapolation algorithm predicts that the file system will run out of
  space within a certain amount of time. If the time period is less than 24
  hours, this is a `Warning` alert. If the time is less than 4
  hours, this is a `Critical` alert.

## Impact

As a file system starts to get low on space, system performance usually
degrades gradually.

If a file system fills up and runs out of space, processes that need to write
to the file system can no longer do so, which can result in lost data and
system instability.

## Diagnosis

* Study recent trends of file system usage on a dashboard. Sometimes, a periodic
pattern of writing and cleaning up in the file system can cause the linear
prediction algorithm to trigger a false alert.

* Use the Linux operating system tools and utilities to investigate what
directories are using the most space in the file system. Is the issue an
irregular condition, such as a process failing to clean up behind itself and
using a large amount of space? Or does the issue seem to be related to
organic growth?

To assist in your diagnosis, watch the following metric in PromQL:

```console
node_filesystem_free_bytes
```

Then, check the `mountpoint` label for the alert.

## Mitigation

If the `mountpoint` label is `/`, `/sysroot` or `/var`, remove unused images to
resolve the issue:

1. Debug the node by accessing the node file system:

    ```console
    $ NODE_NAME=<instance label from alert>
    $ oc -n default debug node/$NODE_NAME
    $ chroot /host
    ```

1. Remove dangling images:

    ```console
    $ podman images -q -f dangling=true | xargs --no-run-if-empty podman rmi
    ```

1. Remove unused images:

    ```console
    $ podman images | grep -v -e registry.redhat.io -e "quay.io/openshift" -e registry.access.redhat.com -e docker-registry.usersys.redhat.com -e docker-registry.ops.rhcloud.com -e rhmap | xargs --no-run-if-empty podman rmi 2>/dev/null
    ```

1. Exit debug:

    ```console
    $ exit
    $ exit
    ```
