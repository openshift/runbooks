# NodeFilesystemFilesFillingUp

## Meaning

The `NodeFilesystemFilesFillingUp` alert is similar to the
[NodeFilesystemSpaceFillingUp][1] alert, but predicts that the file system will
run out of inodes rather than bytes of storage space. The alert triggers at a
`critical` level when the file system is predicted to run out of available
inodes within four hours.

## Impact

A node's file system becoming full can have a widespread negative impact. The
issue might cause any or all of the applications scheduled to that node to
experience anything from degraded performance to becoming fully inoperable,
Depending on the node and file system involved, this issue can pose a critical
threat to the stability of a cluster.

## Diagnosis

Note the `instance` and `mountpoint` labels from the alert. You can graph the
usage history of this file system by using the following query in the OpenShift
web console:

```text
node_filesystem_files_free{
  instance="<value of instance label from alert>",
  mountpoint="<value of mountpoint label from alert>"
}
```

You can also open a debug session on the node and use standard Linux utilities
to locate the source of the usage:

```console
$ MOUNT_POINT='<value of mountpoint label from alert>'
$ NODE_NAME='<value of instance label from alert>'

$ oc debug "node/$NODE_NAME"
$ df -hi "/host/$MOUNT_POINT"
```

Note that in many cases a file system that is running out of inodes will still
have available storage. Running out of inodes is often caused when an
application creates many small files.

## Mitigation

The number of inodes allocated to a file system is usually based on the storage
size. You might be able to solve the problem, or at least delay the negative
impact of the problem, by increasing the size of the storage volume. Otherwise,
determine the application that is creating large numbers of small files and
then either adjust its configuration or provide it with dedicated storage.

[1]: https://github.com/openshift/runbooks/blob/master/alerts/cluster-monitoring-operator/NodeFilesystemSpaceFillingUp.md
