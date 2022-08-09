# Southbound Stale Alert

## Meaning

The alert `SouthboundStale` is triggered when OVN southbound DB has not been
written to for 2 minutes or longer by OVN northd.
Therfore, any networking control plane changes are not being updated
in the Southbound DB for cluster resources to consume.
Make sure that the `NorthBoundStale`, `NoOvnMasterLeader`,`NoRunningOvnMaster`
alerts are not firing.
If they are, triage them before continuing here.

## Impact

Networking control plane is degraded.
Networking Configuration updates applied to the cluster will not be applied.
Existing workloads should continue to have connectivity.

## Diagnosis

There are a few scenarios that can cause this alert to trigger.

1. OVN SDBD is not functional.

   Check to see if the SBDB is running without errors
and has a RAFT cluster leader.

   Check the containers names on the ovn-kube master pods with this command:

   ```shell
   oc get pods -n=openshift-ovn-kubernetes -o jsonpath='{range .items[*]}{"\n"}{.metadata.name}{":\t"}{range .spec.containers[*]}{.name}{", "}{end}{end}' |sort | grep master
   ```

   You should also run `oc get pods -n=openshift-ovn-kubernetes | grep master`
   to ensure that all containers on your master pods are running.

   To see the pod events you can run
`oc describe pod/<podname> -n <namespace>`

   To see the logs of a specific container within the pod you can run
   `oc logs <podname> -n <namespace> -c <containerName>`

   Find the RAFT leader.

   You can run this script to do this.
   ```shell
   LEADER="not found"; for MASTER_POD in $(oc -n openshift-ovn-kubernetes get pods -l app=ovnkube-master -o=jsonpath='{.items[*].metadata.name}'); do RAFT_ROLE=$(oc exec -n openshift-ovn-kubernetes "${MASTER_POD}" -c sbdb -- bash -c "ovn-appctl -t /var/run/ovn/ovnsb_db.ctl cluster/status OVN_Southbound 2>&1 | grep \"^Role\""); if echo "${RAFT_ROLE}" | grep -q -i leader; then LEADER=$MASTER_POD; break; fi; done; echo "sbdb leader ${LEADER}"
   ```

   If this does not work, exec into each sbdb container on the master pods with
   `oc exec -n openshift-ovn-kubernetes -it <ovnkube-master podname> -c sbdb -- bash`
   and then run:
   `ovs-appctl -t /var/run/ovn/ovnsb_db.ctl cluster/status OVN_Southbound`

   You should see a role that will either say leader or follower.

   A common cause of database leader issues is that one of the database
   servers is unable to participate with the other raft peers due to
   mismatching cluster ids. Due to this, they will be unable to elect
   a database leader.

   Try restarting the ovn-kube master pods to resolve this issue.

2. OVN northd is not functional.

   If northd is down, then the logical flows in the northbound database will not
   be translated to the logical datapath flows in the southbound database.
   If you have made it through the first step, you will already know if the
   container is running or not, so you should check the logs to see if there
   are any errors.

3. OVN nortbound database is not functioning.

   Check to see if the `NorthboundStaleAlert` is firing. If the nbdb container
   on ovn-kubernetes master is not running, check the container logs and
   proceed from there.

4. OVN northd cannot connect to one or both of the NBDB/SBDB leader or OVN
   northd is overloaded.

   You can check in the logs of the northd container on the ovn-kube master
   pod with the active instance of ovn-northd. To determine this you must
   exec into the ovnkube-master container on each ovnkube-master pod,
   then run this command:
   `curl 127.0.0.1:29105/metrics | grep ovn_northd_status`

   This will return the status of northd on that ovnkube-master pod. Once you
   find the active instance of ovn-northd, you can check the logs of the
   northd container on that ovnkube-master pod.

   If northd is overloaded, there will be logs in the container along the
   lines of `dropped x number of log messages due to excessive rate` or a
   message that contains `(99% CPU usage)` or some other high percentage
   CPU usage.

   You can also check the cpu usage of the container by logging into your
   openshift cluster console. In the Observe section on the sidebar, click
   metrics, then run this query
   `container_cpu_usage_seconds_total{pod="$PODNAME",
   container="$CONTAINERNAME"}`

## Mitigation

Mitigation will depend on what was found in the diagnosis section.

As a general fix, you can try restarting the ovn-k master pods.

If the issue persists, reach out to the SDN team on #forum-sdn.
