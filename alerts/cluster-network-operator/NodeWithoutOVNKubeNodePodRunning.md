# NodeWithoutOVNKubeNodePodRunning

## Meaning

The `NodeWithoutOVNKubeNodePodRunning` alert is triggered when one or more Linux
nodes do not have a running OVNkube-node pod for a period of time.

## Impact

This is a warning alert. Existing workloads on the node may continue to have
connectivity but any additional workloads will not be provisioned on the node.
Any network policy changes will not be implemented on existing workloads on the
node.

### Fix alerts before continuing

Resolve any alerts that may cause this alert to fire:
[Alert hierarchy](./hierarchy/alerts-hierarchy.svg)

## Diagnosis

Check the nodes which should have the ovnkube-node running.

    oc get node -l kubernetes.io/os!=windows

Check the expected running replicas of ovnkube-node.

    oc get daemonset ovnkube-node -n openshift-ovn-kubernetes

Check the ovnkube-node pods status on the nodes.

    oc get po -n openshift-ovn-kubernetes -l app=ovnkube-node -o wide

Describe the pod if there is non-running ovnkube-node pod.

    oc describe po -n openshift-ovn-kubernetes <ovnkube-node-name>

Check the pod logs for the failing ovnkube-node pods

    oc logs <ovnkube-node-name> -n openshift-ovn-kubernetes --all-containers

## Mitigation

Mitigation for this alert is not possible to understand in advance.

If you are seeing that any of the ovnkube-node pods is not in Running status,
you can try to delete the pod and let it being recreated by the daemonset
controller.

    oc delete po <ovnkube-node> -n openshift-ovn-kubernetes

If the issue cannot be fixed by recreating the pod, reboot of the affected node
might be an option to refresh the full stack (include OVS on the node).
