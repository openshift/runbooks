# NoOvnClusterManagerLeader

## Meaning

This [alert][NoOvnClusterManagerLeader] is triggered when ovn-kubernetes
cluster does not have a leader for more than 10 minutes.

## Impact

When ovnkube cluster manager is unable to elect a leader (via kubernetes lease
API), networking control plane is degraded.
Networking configuration updates applied to the cluster will not be
implemented while there is no OVN Kubernetes cluster manager leader.
Existing workloads should continue to have connectivity.
OVN-Kubernetes control plane is not functional.

## Diagnosis

### Fix alerts before continuing

Check to ensure the following alerts are not firing and resolved before
continuing as they may cause this alert to fire:

[Alert hierarchy](./hierarchy/alerts-hierarchy.svg)

### OVN-kubernetes control-plane pods

Check if all the ovn-kube control planes are running:

    oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-control-plane

Check if there is a leader for the ovn-kubernetes cluster:

    oc get lease -n openshift-ovn-kubernetes

    name: ovn-kubernetes-master
    holder: ovnkube-control-plane-8444dff7f9-6qdl4
    age: 89m

`holder` shown above, contains the node name where the leader pod
resides.
Check the logs for any of the running ovnkube-control-plane to see if there is
leader election happened and if there is an error occurred.

    oc logs -n openshift-ovn-kubernetes ovnkube-control-plane-xxxxx --all-containers | grep elect

## Mitigation

### If the control plane nodes are not running

Follow the steps described in the [disaster and recovery documentation][dr_doc].

### If the cluster network operator is reporting error

Follow the condition reported in the operator to fix the operator managed services.

### If one of the ovnkube-control-plane pods is not running

The ovnkube-cluster-manager container in the ovn kubernetes control-plane pod
should run the leader election if the old leader is down, you may need to 
check the other running ovnkube-control-plane pods' logs for more
information about why the election failed.

### If all the ovnkube-control-plane pods are not running

Check the status of the ovnkube-control-plane pods, and follow the
[Pod lifecycle][Pod lifecycle] to see what is blocking the pods to be running.

### If all the ovnkube-control-plane pods are running

Follow the steps above: [OVN-Kubernetes master pods](#ovn-kubernetes-control-plane-pods)

[NoOvnClusterManagerLeader]: https://github.com/openshift/cluster-network-operator/blob/master/bindata/network/ovn-kubernetes/self-hosted/multi-zone-interconnect/alert-rules-control-plane.yaml
[Pod lifecycle]: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/
[dr_doc]: https://docs.openshift.com/container-platform/latest/backup_and_restore/control_plane_backup_and_restore/disaster_recovery/about-disaster-recovery.html
