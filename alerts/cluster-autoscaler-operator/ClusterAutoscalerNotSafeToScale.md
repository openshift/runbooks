# ClusterAutoscalerNotSafeToScale

## Meaning

The cluster autoscaler has detected that the number of unready nodes is too high
and it is not safe to continute scaling operations. It makes this determination
by checking that the number of ready nodes is greater than the minimum ready count
(default of 3) and the ratio of unready to ready nodes is less than the maximum
unready node percentage (default of 45%). If either of those conditions are not
true then the cluster autoscaler will enter an unsafe to scale state until the
conditions change.

## Possible Causes
* The cluster has too many nodes in an unready state.
* A large number of new nodes have been created and are taking longer than 15
  minutes to join the cluster.

## Resolution
This alert is indicating an issue with nodes not reaching a ready state. You
should investigate the logs associated with your cloud provider controllers and
the Machine API resources to discover the root cause. For more information on
why nodes, or machines, might not become ready please see the
[Machine API FAQ](https://github.com/openshift/machine-api-operator/blob/master/FAQ.md).
