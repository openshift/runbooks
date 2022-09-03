# OVNKubernetesOvnKubeMasterLostConnectionToDatabases

## Meaning

This alert is triggered when OVN-Kubernetes master control plane lose
connection to OVN SouthBound/NorthBound databases for 5 minutes
or longer.

## Impact
Networking control plane is not functional. Networking configuration updates
will not be applied to the cluster.
Without a functional networking control plane, existing workloads may continue
to be partially functional,
but new workloads will not be functional.
Updates required for functioning Kubernetes services will not be performed.

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

2. OVN nortbound database is not functioning.

   Check to see if the `NorthboundStaleAlert` is firing. If the nbdb container
   on ovn-kubernetes master is not running, check the container logs and
   proceed from there.

## Mitigation

The appropriate mitigation will be very different depending on the cause of the
error discovered in the diagnosis.
As a general fix, you can try restarting the ovn-k master pods.

If the issue persists, reach out to the SDN team on #forum-sdn.