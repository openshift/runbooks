# OVNKubernetesNorthdInactive

## Meaning

This alert fires when there are is no active instance of OVN northd within a
specific availability domain. For OCP clusters at versions 4.13 or earlier, the
availability domain is the entire cluster. For OCP clusters at versions 4.14 or
later, the availability domain is a cluster node.

## Impact

ovn-northd is a daemon that translates the logical network flows from the OVN
Northbound Database into the physical datapath flows in the OVN Southbound
database. If there are no active instances of ovn-northd, then this action will
not occur, which will cause a degraded network. Existing workloads may continue
to have connectivity but any additional workloads will not be provisioned. Any
network policy changes will not be implemented on existing workloads. For OCP
clusters at versions 4.13 or earlier the affected domain is the entire cluster.
For OCP clusters at versions 4.14 or later, the affected domain is only the
specific node for which the alert was fired.

### Fix alerts before continuing

Resolve any alerts that may cause this alert to fire: [Alert
hierarchy](./hierarchy/alerts-hierarchy.svg)

## Diagnosis

Investigate the health of the affected ovn-northd processes.

For OCP clusters at versions 4.13 or earlier, the affected ovn-northd processes
run in the northd container of ovnkube-master pods. Find out what those pods are
and exec into them:
```shell
oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-master
oc exec -it <ovnkube-master-podname> -c northd -- bash
```

For OCP clusters at versions 4.14 or later, the affected ovn-northd process runs
in the northd container of the ovnkube-node pod for the affected node. Find out
what pod is that and exec into it:
```shell
oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-node -o wide
oc exec -it <ovnkube-node-podname> -c northd -- bash
```

Then run:
```shell
curl 127.0.0.1:29105/metrics | grep northd
```
This will show you the cluster metrics associated with northd

Next, check if the northd instance is active:
```shell
ovn-appctl -t ovn-northd status
```
The result should be Status:active


## Mitigation

Mitigation will depend on what was found in the diagnosis section.

As a general fix, you can try exiting the affected ovn-northd procesess with
```shell
ovn-appctl -t ovn-northd exit
```
which should cause the container running northd to restart. If this does not
work you can try restarting the pods where the affected ovn-northd procesess are
running.

Contact the incident response team in your organisation if fixing the issue is
not apparent.
