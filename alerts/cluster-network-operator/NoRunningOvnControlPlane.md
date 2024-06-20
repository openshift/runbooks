# NoRunningOvnControlPlane

## Meaning

This alert is triggered when there are no OVN-Kubernetes control plane pods
[Running][PodRunning] for `5m`.

## Impact

When OVN-Kubernetes control plane is not running, networking control plane is
degraded. A subset of networking control plane functionality is degraded. This
includes, but not limited to the following networking control plane
functionality:
* Node resource allocation
* EgressIP assignment/re-assignment and health checks
* EgressService node allocation/re-allocation
* DNS name resolver functionality for EgressFirewall
* OVN secondary networks IPAM

### Fix alerts before continuing

Resolve any alerts that may cause this alert to fire: [Alert
hierarchy](./hierarchy/alerts-hierarchy.svg)

## Diagnosis

### Control plane issue

This can occur when multiple control plane nodes are powered off or are unable
to connect with each other via the network. Check that all control plane nodes
are powered on and that network connections between each machine are functional:
```shell
oc get node -l node-role.kubernetes.io/control-plane=""
```

### Cluster network operator

Cluster network operator (CNO) manages the CRUD operations for OVN-Kubernetes
daemonset. Verify the CNO is running:
```shell
oc -n openshift-network-operator get pods -l name=network-operator
```

Verify the [CNO](https://github.com/openshift/cluster-network-operator/) is
functioning without error by checking the operator Status:
```shell
oc get co network
```

If the network is degraded, you can see the full error message by describing the
object. Pay attention to any error message reported by Status Conditions of type
Degraded:
```shell
oc describe co network
```

A successful reconcile for this daemonset looks like this in the CNO logs:
```shell
I0611 11:13:35.048771       1 log.go:245] Apply / Create of (apps/v1, Kind=Deployment) openshift-ovn-kubernetes/ovnkube-control-plane was successful
```

### OVN-Kubernetes control-plane pod

Check that all pods of the `ovnkube-control-plane` deployment are READY:
```shell
oc get deploy -n openshift-ovn-kubernetes ovnkube-control-plane
oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-control-plane
```

If one of the `ovnkube-control-plane` pods is not READY, check the overall
status of the pod and that of specific containers:
```shell
oc describe pod/<podname> -n openshift-ovn-kubernetes
```

After understanding which container is not starting successfully, gather the
runtime logs from that container:
```shell
oc logs <podname> -n openshift-ovn-kubernetes -c <container>
```
You may need to use `--previous` command with `oc logs` command to get the logs
of the previous execution run of a container. Pay close attention to any log
output starting with "E" for Error.

## Mitigation

The appropriate mitigation will be very different depending on the cause of the
error discovered in the diagnosis. Investigate the issue using the steps
outlined in diagnosis and contact the incident response team in your
organisation if fixing the issue is not apparent.

[PodRunning]:
    https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase
