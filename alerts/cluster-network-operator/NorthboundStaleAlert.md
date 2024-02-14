# Northbound Stale Alert

## Meaning

The [alert][NorthboundStale] will be triggered if the ovn-kube control-plane
process is not functioning, if the northbound database on the node is not
functioning, or if connectivity between the node and the database is broken.

## Impact

Networking control plane on that node is degraded and Networking Configuration
updates applied to the node will not be applied.

Existing workloads should continue to have connectivity, but the OVN-Kubernetes
control plane and/or the OVN southbound database may not be functional.

## Fix alerts before continuing

Resolve any alerts that may cause this alert to fire:
[Alert hierarchy](./hierarchy/alerts-hierarchy.svg)

## Diagnosis

Investigate the causes that can trigger this alert.

1. Is the ovn-controller process running, i.e. is the container running.

   If it is not, check the logs and proceed from there.

   You can check the container names on the ovn-kube node pods with this command:

   ```shell
   oc get pods -n=openshift-ovn-kubernetes -o jsonpath='{range .items[*]}{"\n"}{.metadata.name}{":\t"}{range .spec.containers[*]}{.name}{", "}{end}{end}' |sort 
   ```

   You should also run
   `oc get pods -n=openshift-ovn-kubernetes | grep node`
   to ensure that all containers on the node pods are ready.

   To see the pod logs you can run `oc describe pod/<podname> -n <namespace>`

   To see the logs of a specific container within the pod you can run
   `oc logs <podname> -n <namespace> -c <containerName>`

2. Is OVN Northbound Database functioning.

   Check to see if the northbound database containers are running without errors.


   If you have made it through the debug steps in step one, then you should already
   know that the containers are healthy on the node and control plane pods.

## Mitigation

Mitigation will depend on what was found in the diagnosis section.
As a general fix, you can try restarting the ovn-k node pod that is effected.
If the issue persists, reach out to the SDN team on #forum-sdn.

[NorthBoundStale]: https://github.com/openshift/cluster-network-operator/blob/master/bindata/network/ovn-kubernetes/self-hosted/multi-zone-interconnect/alert-rules-control-plane.yaml