# NodeFilesystemAlmostOutOfSpace

## Meaning

The `NodeFilesystemAlmostOutOfSpace` alert is similar to the
[NodeFilesystemSpaceFillingUp][1] alert, but rather
than being based on a prediction that a file system will become full in a
certain amount of time, it uses simple static thresholds. This alert triggers
at a `warning` level when 5% of space remains in the file system, and at a
`critical` level when 3% of space remains.

## Impact

A node's file system becoming full can have a widespread negative impact. This
issue can cause any or all of the applications scheduled to that node to
experience anything from degraded performance to becoming fully inoperable.
Depending on the node and file system involved, this issue can pose a critical
threat to the stability of the cluster.

## Diagnosis

Refer to the [NodeFilesystemSpaceFillingUp][1] runbook.

## Mitigation

Refer to the [NodeFilesystemSpaceFillingUp][1] runbook.

[1]: https://github.com/openshift/runbooks/blob/master/alerts/cluster-monitoring-operator/NodeFilesystemSpaceFillingUp.md
