# ExtremelyHighIndividualControlPlaneMemory

## Meaning

This alert fires when the memory use by an instance within control plane
nodes exceeds 90% of total memory available on that node.  
This is a critical level alert.

## Impact

At such very low free memory, the cluster becomes unstable you can expect slow
responses from kube-apiserver or failing requests specially from etcd.
Moreover, OOM kill is expected, which negatively influences the pod scheduling.

## Diagnosis

To fix this, increase memory of the affected node of control plane nodes.

## Mitigation

Responding to a `HighOverallControlPlaneMemory` alert can mitigate this situation.
`HighOverallControlPlaneMemory` is a warning level alert, which fires
when the total memory usage across all
control-plane nodes crosses 60% for more than 1 hour.

You could scale up the memory of the control plane node if this situation
stays up for long time to prevent the critical alert from firing.
