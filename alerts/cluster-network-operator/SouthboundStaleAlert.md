# Southbound Stale Alert

## Meaning

The [alert][SouthboundStale] is triggered when OVN southbound DB has not been
written to for 2 minutes or longer by OVN northd.
Therfore, any networking control plane changes are not being updated
in the Southbound DB for cluster resources to consume.
Make sure that the `NorthBoundStale`, `NoOvnMasterLeader`,`NoRunningOvnMaster`
alerts are not firing.
If they are, triage them before continuing here.

## Impact

Networking control plane is degraded.
Networking Configuration updates applied to the effected node will not be applied.
Existing workloads should continue to have connectivity.

### Fix alerts before continuing

Resolve any alerts that may cause this alert to fire:
[Alert hierarchy](./hierarchy/alerts-hierarchy.svg)

## Diagnosis

There are a few scenarios that can cause this alert to trigger.

1. OVN SDBD is not functional.

   Check to see if the SBDB is running without errors

   Check the containers names on the ovn-kube node pods with this command:

   ```shell
   oc get pods -n=openshift-ovn-kubernetes -o jsonpath='{range .items[*]}{"\n"}{.metadata.name}{":\t"}{range .spec.containers[*]}{.name}{", "}{end}{end}' |sort | grep node
   ```

   You should also run `oc get pods -n=openshift-ovn-kubernetes | grep controller`
   to ensure that all containers on your master pods are running.

   To see the pod events you can run
   `oc describe pod/<podname> -n <namespace>`

   To see the logs of a specific container within the pod you can run
   `oc logs <podname> -n <namespace> -c <containerName>`


2. OVN northd is not functional.

   If northd is down, then the logical flows in the northbound database will not
   be translated to the logical datapath flows in the southbound database.
   If you have made it through the first step, you will already know if the
   container is running or not, so you should check the logs of the northd
   container on the node to see if there are any errors.

3. OVN nortbound database is not functioning.

   Check to see if the `NorthboundStaleAlert` is firing. If the nbdb container
   on ovn-kubernetes master is not running, check the container logs and
   proceed from there.

   You can also check the cpu usage of the container by logging into your
   openshift cluster console. In the Observe section on the sidebar, click
   metrics, then run this query
   `container_cpu_usage_seconds_total{pod="$PODNAME",
   container="$CONTAINERNAME"}`

## Mitigation

Mitigation will depend on what was found in the diagnosis section.

As a general fix, you can try restarting the ovn-k master pods.

If the issue persists, reach out to the SDN team on #forum-sdn.

[SouthboundStale]: https://github.com/openshift/cluster-network-operator/blob/master/bindata/network/ovn-kubernetes/self-hosted/multi-zone-interconnect/alert-rules-control-plane.yaml
