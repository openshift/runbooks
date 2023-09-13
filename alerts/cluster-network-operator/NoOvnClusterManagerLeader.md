# NoOvnClusterManagerLeader

## Meaning

This alert is triggered when ovn-kubernetes cluster does not have a leader for
more than 10 minutes.

## Impact

When OVN-Kubernetes cluster manager is unable to elect a leader (via kubernetes
lease API), networking control plane is degraded. A subset of networking control
plane functionality is degraded. This includes, but not limited to the following
networking control plane functionality:
* Node resource allocation
* EgressIP assignment/re-assignment and health checks
* EgressService node allocation/re-allocation
* DNS name resolver functionality for EgressFirewall
* OVN secondary networks IPAM

## Diagnosis

### Fix alerts before continuing

Resolve any alerts that may cause this alert to fire: [Alert
hierarchy](./hierarchy/alerts-hierarchy.svg)

### OVN-kubernetes control-plane pods

Check that all pods of the `ovnkube-control-plane` deployment are READY:
```shell
oc get deploy -n openshift-ovn-kubernetes ovnkube-control-plane
oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-control-plane
```

Check if there is a leader for the ovn-kubernetes cluster:
```shell
oc get lease -n openshift-ovn-kubernetes

NAME                    HOLDER                                   AGE
ovn-kubernetes-master   ovnkube-control-plane-85f4486946-jx95r   32m
```

`HOLDER` shown above is the leader pod.

Check the logs of the of the `ovnkube-control-plane` deployment pods to see if
leader election happened or if an error occurred:
```shell
oc logs <podname> -n openshift-ovn-kubernetes --all-containers | grep elect
```

## Mitigation

### If the control plane nodes are not running

Follow the steps described in the [disaster and recovery documentation][dr_doc].

### If the cluster network operator is reporting error

Follow the condition reported in the operator to fix the operator managed
services.

### If one of the ovnkube-control-plane pods is not running

The ovnkube-cluster-manager container in the ovn kubernetes control-plane pod
should run the leader election if the old leader is down, you may need to check
the other running ovnkube-control-plane pods' logs for more information about
why the election failed.

### If all the ovnkube-control-plane pods are not running

Check the status of the ovnkube-control-plane pods, and follow the [Pod
lifecycle][Pod lifecycle] to see what is blocking the pods to be running.

### If all the ovnkube-control-plane pods are running

Follow the steps above: [OVN-Kubernetes control plane
pods](#ovn-kubernetes-control-plane-pods)

[Pod lifecycle]:
    https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/
[dr_doc]:
    https://docs.openshift.com/container-platform/latest/backup_and_restore/control_plane_backup_and_restore/disaster_recovery/about-disaster-recovery.html
