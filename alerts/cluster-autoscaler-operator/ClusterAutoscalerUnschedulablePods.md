# ClusterAutoscalerUnschedulablePods

## Meaning

The cluster autoscaler is unable to scale up and is alerting that there are
unschedulable pods because of this condition.

## Possible Causes
* The autoscaler is unable to create new machines due to replica limits on the
  MachineAutoscalers.
* The autoscaler is unable to create new machines due to maximum node, CPU, or
  RAM limits on the ClusterAutoscaler.
* Kubernetes is waiting for new nodes to become ready before scheduling pods to
  them.

## Resolution
In many cases this alert is normal and expected depending on the configuration
of the autoscaler. You should check the replica limits in the MachineAutoscaler
resources to ensure they are large enough. You should also check the maximum
totals nodes, CPU, and RAM limits in the ClusterAutoscaler resource to ensure
they are valid.

In rare cases it is possible that the cloud provider is taking longer than 20
minutes to create new nodes. This should be investigated with the cloud provider
and their specific process for node creation.
