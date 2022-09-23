# OvnKubernetes Northd Inactive Alert

## Meaning

The alert `OVNKubernetesNorthdInactive` fires when there are no active
instances of Northd within the high availabilty set.

## Impact

ovn-northd is a centralized daemon that translates the logical network flows from
the OVN Northbound Database into the physical datapath flows in the
OVN Southbound database. If there are no active instances of ovn-northd, then
this action will not occur, which will cause a degraded network.
Existing workloads should remain stable, but any attempt at altering the
network control plane will fail.

### Fix alerts before continuing
Resolve any alerts that may cause this alert to fire [Alert Hierarchy](./hierarchy/alerts-hierarchy.svg)
## Diagnosis

First, exec into the northd containers on the master pods with
```shell
`oc exec -it <ovn-master-podname> -c ovn-northd -- bash`
```
and run
```shell
`curl 127.0.0.1:29105/metrics | grep northd`
```
This will show you the cluster metrics associated with northd

Next, use the cli to check if theres an active northd instance with
while still in the northd
```shell
`ovn-appctl -t ovn-northd status` 
```
The result should be Status:active


## Mitigation

Mitigation will depend on what was found in the diagnosis section.

As a general fix, you can try exiting all the ovn-northd procesess
with
```shell
`ovn-appctl -t ovn-northd exit`
```
which should cause the container running northd to restart.
If this does not work you can try restarting the ovn-k master pods.

If the issue persists, reach out to the SDN team on #forum-sdn.
