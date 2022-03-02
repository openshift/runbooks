# NoOvnMasterLeader

## Meaning

This alert is triggered when ovn-kubernetes cluster does not have a
leader for more than 10 minute.

## Impact

When ovnkube-master is unable to elect a leader (via kubernetes leader
election API), networking control plane is degraded.
Networking configuration updates applied to the cluster will not be
implemented while there is no OVN Kubernetes leader.
Existing workloads should continue to have connectivity.
OVN-Kubernetes control plane is not functional.

## Diagnosis

### Fix alerts before continuing

Check to ensure the following alerts are not firing and resolved before
continuing as they may cause this alert to fire:

- [NoRunningOvnMaster](./NoRunningOvnMaster)

### OVN-kubernetes master pods

Check if all the ovn-kube masters are running:

    oc get ds -n openshift-ovn-kubernetes ovnkube-master
    oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-master

Check if there is a leader for the ovn-kubernetes cluster:

    oc get cm -n openshift-ovn-kubernetes ovn-kubernetes-master -o json | \
    jq .metadata.annotations

    control-plane.alpha.kubernetes.io/leader: '{"holderIdentity":"ip-10-0-
    131-116.ec2.internal","leaseDurationSeconds":60,"acquireTime":
    "2022-03-07T11:18:21Z","renewTime":"2022-03-07T11:25:32Z",
    "leaderTransitions":5}'

Check the logs for any of the running ovnkube-master to see if there is
leader election happened and if there is an error occurred.

    oc logs ovnkube-master-xxxxx --all-containers | grep elect

## Mitigation

### If the control plane nodes are not running

Follow the steps described in the [disaster and recovery docs](dr_doc)

### If the cluster network operator is reporting error

Follow the condition reported in the operator to fix the operator managed services.

### If one of the ovnkube-master pods is not running

The ovnkube-master container in the ovn kubernetes master pod should run the
leader election if the old leader is down, you may need to check the other
running ovnkube-master pods' logs for more information about why the election
failed.

### If all the ovnkube-master pods are not running

Check the status of the ovnkube-master pods, and follow the
[Pod lifecycle](pod_lifecycle) to see what is blocking the pods to be running.

### If all the ovnkube-master pods are running

Follow the steps above [OVN-kubernetes master pods]

[dr_doc]:(https://docs.openshift.com/container-platform/4.9/backup_and_restore/control_plane_backup_and_restore/disaster_recovery/about-disaster-recovery.html)
[pod_lifecycle]:(https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/)
