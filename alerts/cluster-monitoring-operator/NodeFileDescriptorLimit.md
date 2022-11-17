# NodeFileDescriptorLimit

## Meaning

The `NodeFileDescriptorLimit` alert is triggered when a node's kernel is
running out of available file descriptors. A `warning` level alert triggers at
greater than 70% usage, and a `critical` level alert triggers at greater than
90% usage.

## Impact

Applications on the node might no longer be able to open and operate on
files, which is likely to have severe negative consequences for anything
scheduled on this node.

## Diagnosis

Open a shell on the node and use standard Linux utilities to diagnose the issue:

```console
$ NODE_NAME='<value of instance label from alert>'

$ oc debug "node/$NODE_NAME"
# sysctl -a | grep 'fs.file-'
fs.file-max = 1597016
fs.file-nr = 7104       0       1597016
# lsof -n
```

## Mitigation

Reduce the number of files opened simultaneously either by adjusting application
configuration or by moving some applications to other nodes.
