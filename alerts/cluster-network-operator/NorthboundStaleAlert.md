# Northbound Stale Alert

## Meaning

The alert `NorthboundStale` will be triggered if the ovn-kube master process is
not functioning, if the northbound database is not functioning, or if
connectivity between the ovnkube master and the database is broken.

This alert will be triggered if `NoRunningOvnMaster` or if `NoOvnMasterLeader`
is firing, so check those alerts before proceeding here.

## Impact

Networking control plane is degraded and Networking Configuration updates applied
to the cluster will not be applied.

Existing workloads should continue to have connectivity, but the OVN-Kubernetes
control plane and/or the OVN southbound database may not be functional.

## Diagnosis

Investigate the causes that can trigger this alert.

1. Is the ovnkube-master process running, i.e. is the container running.

   If it is not, check the logs and proceed from there.

   You can check the container names on the ovn-kube master nodes with this command:

   ```shell
   oc get pods -n=openshift-ovn-kubernetes -o jsonpath='{range .items[*]}{"\n"}{.metadata.name}{":\t"}{range .spec.containers[*]}{.name}{", "}{end}{end}' |sort | grep master
   ```

   You should also run
   `oc get pods -n=openshift-ovn-kubernetes | grep master`
   to ensure that all containers are ready.

   To see the pod logs you can run `oc describe pod/<podname> -n <namespace>`

   To see the logs of a specific container within the pod you can run
`oc logs <podname> -n <namespace> -c <containerName>`

2. Is OVN Northbound Database functioning.

   Check to see if the northbound database containers are running without errors.

   If there are no errors, ensure there is an OVN Northbound Database Leader.

   If you have made it through the debug steps in step one, then you should already
   know that the containers are healthy on the master pod/s, and you should now
   ensure that there is a northbound database leader.

   To find the database leader, you can run this script.

   ```shell
   LEADER="not found"; for MASTER_POD in $(oc -n openshift-ovn-kubernetes get pods -l app=ovnkube-master -o=jsonpath='{.items[*].metadata.name}'); do RAFT_ROLE=$(oc exec -n openshift-ovn-kubernetes "${MASTER_POD}" -c nbdb -- bash -c "ovn-appctl -t /var/run/ovn/ovnnb_db.ctl cluster/status OVN_Northbound 2>&1 | grep \"^Role\""); if echo "${RAFT_ROLE}" | grep -q -i leader; then LEADER=$MASTER_POD; break; fi; done; echo "nbdb leader ${LEADER}"
   ```

   If this does not work, exec into each nbdb container on the master pods with:
   `oc exec -n openshift-ovn-kubernetes -it <ovnkube-master podname> -c nbdb -- bash`
   and then run:
   `ovs-appctl -t /var/run/ovn/ovnnb_db.ctl cluster/status OVN_Northbound`
   You should see a role that will either say leader or follower.

   A common cause of database leader issues is that one of the database servers
   is unable to participate with the other raft peers due to mismatching cluster
   ids. Due to this, they will be unable to elect a database leader.

   Try restarting the ovn-kube master pods to resolve this issue.

3. Lastly, check to make sure that the connectivity between ovnkube-master leader
   and OVN northbound database leader is healthy.

   To determine what node the ovnkube-master leader is on, check the value of
   `holderIdentity`:

   ```shell
   oc get lease -n ovn-kubernetes ovn-kubernetes-master -o yaml
   ```

   Then get the logs of the ovnkube-master container on the ovnkube-master pod on
   that node with
   `oc logs <podname> -n <namespace> -c ovnkube-master`

   You should see a message along the lines of:
   `"msg"="trying to connect"
"database"="OVN_Northbound" "endpoint"="tcp:172.18.0.4:6641"`.

   This message indicates that the master cannot connect to the database. A
   successful connection message will appear in the logs if the master has
   connected to the database.

## Mitigation

Mitigation will depend on what was found in the diagnosis section.
As a general fix, you can try restarting the ovn-k master pods.
If the issue persists, reach out to the SDN team on #forum-sdn.
