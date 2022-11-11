# NodeFilesystemAlmostOutOfFiles

## Meaning

The `NodeFilesystemAlmostOutOfFiles` alert is similar to the
[NodeFilesystemSpaceFillingUp][1] alert, but rather
than being based on a prediction that a filesystem will run out of inodes in a
certain amount of time, it uses simple static thresholds. The alert triggers
at a `warning` level when 5% of available inodes remain, and triggers at a
`critical` level when 3% of available inodes remain.

## Impact

When a node's filesystem becomes full, it has a widespread impact. This issue
can cause any or all of the applications scheduled to that node to experience
anything from degraded performance to becoming fully inoperable. Depending on
the node and filesystem involved, this issue could pose a critical threat to
the stability of the cluster.

## Diagnosis

Refer to the [NodeFilesystemFilesFillingUp][1] runbook.

## Mitigation

Refer to the [NodeFilesystemFilesFillingUp][1] runbook.

[1]: https://github.com/openshift/runbooks/blob/master/alerts/cluster-monitoring-operator/NodeFilesystemFilesFillingUp.md
