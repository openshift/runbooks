# ClusterAutoscalerNotSafeToScale

## Meaning

The number of total cores in the cluster has exceeded the maximum number set on
the cluster autoscaler. This is calculated by summing the cpu capacity for all
nodes in the cluster and comparing that number against the maximum cores value
set for the cluster autoscaler (default 320000 cores).

## Possible Causes
* Too many nodes have been created in the cluster.
* Nodes of larger than expected size have joined the cluster.
* Maximum CPU limit on the ClusterAutoscaler is set too low.

## Resolution
This alert is indicating that the cluster autoscaler is unable to continue
scaling out. Depending on your needs and resources this alert may indicate
action is required. If you require more resources in your cluster, a simple
solution is to increase the maximum core count in your ClusterAutoscaler. If you
do not need more resources in your cluster, this condition is non-harmful to the
cluster and the autoscaler will continue to function as normal, with the
exception of creating new nodes. The cluster autoscaler will resume its scale
out functionality once the number of cores in the cluster is fewer than the
maximum.
