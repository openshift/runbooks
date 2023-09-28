# NoRunningOvnControlPlane

## Meaning

This [alert][NoRunningOvnControlPlane] is triggered when there are no
OVN-Kubernetes control plane pods
[Running][PodRunning]
This is a critical-level alert if no OVN-Kubernetes control plane pods
are not running for `10m`.

## Impact
Networking control plane is not functional. Networking configuration updates
will not be applied to the cluster.
Without a functional networking control plane, existing workloads may continue
to be partially functional,
but new workloads will not be functional.
Updates required for functioning Kubernetes services will not be performed.

### Fix alerts before continuing

Resolve any alerts that may cause this alert to fire:
[Alert hierarchy](./hierarchy/alerts-hierarchy.svg)

## Diagnosis
### Control plane issue

This can occur multiple control plane nodes are powered off or are unable to
connect each other via the network. Check that all control plane nodes are
powered and that network connections between each machine are functional.

    oc get node -l node-role.kubernetes.io/master=""

### Cluster network operator
Cluster Network operator (CNO) manages the CRUD operations for OVN-Kubernetes
daemonset.
Verify the CNO is running:

    oc -n openshift-network-operator get pods -l name=network-operator

Verify the [CNO](https://github.com/openshift/cluster-network-operator/) is
functioning without error by checking the operator Status:

    oc get co network

If the network is degraded, you can see the full error message by describing the
object. Pay attention to any error message reported by Status Conditions of type
Degraded:

    oc describe co network


### OVN-Kubernetes control-plane pod
Verify the ovnkube-control-plane pods are running

    oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-control-plane

Verify the nodes the control-plane pods are running on are up

    oc get nodes -l node-role.kubernetes.io/master="" -o name | wc -l

If one of the `ovnkube-control-plane` is not running
then understand which container is failing in the OVN-Kubernetes
control plane pods by describing one of the failing pods with `oc describe pod ...`.
After understanding which container is not starting successfully, gather the
runtime logs from that container.
You may need to use `--previous` command with `oc logs` command to get the
logs of the previous execution run of a container.

Pay close attention to any log output starting with "E" for Error:

    oc -n openshift-ovn-kubernetes logs $OVNKUBE-CONTROL-PLANE-POD-NAME
    --all-containers=true | grep "^E"

## Mitigation

The appropriate mitigation will be very different depending on the cause of the
error discovered in the diagnosis.
Investigate the issue using the steps outlined in diagnosis and contact the
incident response team in your organisation if fixing the issue is not apparent.

[NoRunningOvnControlPlane]: https://github.com/openshift/cluster-network-operator/blob/master/bindata/network/ovn-kubernetes/self-hosted/multi-zone-interconnect/alert-rules-control-plane.yaml
[PodRunning]: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#pod-phase

### Legacy information (pre 4.14/ovn ic)

The alert firing here would be NoRunningOvnMaster, and the trigger
would be if there were no running master pods.

The impact would be the same as detailed above.

On clusters running 4.13 or below, you can verify that the
desired number of daemonsets is equal to the number of k8s
control plane nodes
    oc get ds -n openshift-ovn-kubernetes ovnkube-master
    oc get nodes -l node-role.kubernetes.io/master="" -o name | wc -l
If _READY_ count from the daemonset `ovnkube-master` is not equal to
_DESIRED_ then understand which container is failing in the OVN-Kubernetes
master pod by describing one of the failing pods with `oc describe pod ...`.
