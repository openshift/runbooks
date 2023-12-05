# ODFRBDClientBlocked

## Meaning

This alert indicates that an RBD client might be blocked by Ceph on a specific
node within your Kubernetes cluster. The blocklisting occurs when the
`ocs_rbd_client_blocklisted metric` reports a value of 1 for the node.
Additionally, there are pods in a CreateContainerError state on the same node.
The blocklisting can potentially result in the filesystem for the Persistent
Volume Claims (PVCs) using RBD becoming read-only.
It is crucial to investigate this alert to prevent any disruption to your
storage cluster.

## Impact

High. It is crucial to investigate this alert to prevent any disruption to your
storage cluster.
This may cause the filesystem for the PVCs to be in a read-only state.

## Diagnosis

The blocklisting of an RBD client can occur due to several factors, such as
network or cluster slowness. In certain cases, the exclusive lock contention
among three contending clients (workload, mirror daemon, and manager/scheduler)
 can lead to the blocklist.

## Mitigation

Taint the blocklisted node: In Kubernetes, consider tainting the node that is
blocklisted to trigger the eviction of pods to another node. This approach
relies on the assumption that the unmounting/unmapping process progresses
gracefully. Once the pods have been successfully evicted, the blocklisted node
can be untainted, allowing the blocklist to be cleared. The pods can then be
moved back to the untainted node.

Reboot the blocklisted node: If tainting the node and evicting the pods do not
resolve the blocklisting issue, a reboot of the blocklisted node can be
attempted. This step may help alleviate any underlying issues causing the
blocklist and restore normal functionality.

Please note that investigating and resolving the blocklist issue promptly is
essential to avoid any further impact on the storage cluster.
