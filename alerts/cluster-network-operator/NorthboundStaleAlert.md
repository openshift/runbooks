# NorthboundStale

## Meaning

This alert is triggered when ovnkube-controller or northbound database processes
in a specific availability domain are not functioning correctly or if
connectivity between them is broken. For OCP clusters at versions 4.13 or
earlier, the availability domain is the entire cluster. For OCP clusters at
versions 4.14 or later, the availability domain is a cluster node.

## Impact

Existing workloads may continue to have connectivity but any additional
workloads will not be provisioned. Any network policy changes will not be
implemented on existing workloads. For OCP clusters at versions 4.13 or earlier
the affected domain is the entire cluster. For OCP clusters at versions 4.14 or
later, the affected domain is only the specific node for which the alert was
fired.

## Fix alerts before continuing

Resolve any alerts that may cause this alert to fire: [Alert
hierarchy](./hierarchy/alerts-hierarchy.svg)

## Diagnosis

Investigate the health of the affected ovnkube-controller or northbound database
processes that run in the `ovnkube-controller` and `nbdb` containers
repectively.

For OCP clusters at versions 4.13 or earlier, the containers run in
ovnkube-master pods:
```shell
oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-master
```

For OCP clusters at versions 4.14 or later, the containers run in the
ovnkube-node pod of the affected node, one of:
```shell
oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-node -o wide
```

Check the overall status of the affected pods and of its containers:
```shell
oc describe pod/<podname> -n openshift-ovn-kubernetes
```

Check the logs of a specific container:
```shell
oc logs <podname> -n openshift-ovn-kubernetes -c <container>
```
You may need to use `--previous` command with `oc logs` command to get the logs
of the previous execution run of a container. Pay close attention to any log
output starting with "E" for Error.

> NOTE: The checks below only apply for OCP clusters at versions 4.13 or earlier
> where ovnkube-controller remotely connects to a northbound database leader.

Ensure there is an northbound database Leader. To find the database leader, you
can run this script:
```shell
LEADER="not found"; for MASTER_POD in $(oc -n openshift-ovn-kubernetes get pods -l app=ovnkube-master -o=jsonpath='{.items[*].metadata.name}'); do RAFT_ROLE=$(oc exec -n openshift-ovn-kubernetes "${MASTER_POD}" -c nbdb -- bash -c "ovn-appctl -t /var/run/ovn/ovnnb_db.ctl cluster/status OVN_Northbound 2>&1 | grep \"^Role\""); if echo "${RAFT_ROLE}" | grep -q -i leader; then LEADER=$MASTER_POD; break; fi; done; echo "nbdb leader ${LEADER}"
```

If this does not work, exec into each nbdb container on the master pods:
```shell
oc exec -n openshift-ovn-kubernetes -it <podname> -c nbdb --bash
```
and then run:
```shell
ovs-appctl -t /var/run/ovn/ovnnb_db.ctl cluster/status OVN_Northbound
```
You should see a role that will either say leader or follower.

A common cause of database leader issues is that one of the database servers is
unable to participate with the other raft peers due to mismatching cluster ids.
Due to this, they will be unable to elect a database leader.

Check to make sure that the connectivity between ovnkube-master leader and OVN
northbound database leader is healthy.

To determine what node the ovnkube-master leader is on, check the value of
`holderIdentity`:
```shell
oc get lease -n ovn-kubernetes ovn-kubernetes-master -o yaml
```

Then get the logs of the ovnkube-master container on the ovnkube-master pod on
that node:
```shell
oc logs <podname> -n <namespace> -c ovnkube-master
```

You should see a message along the lines of:
```shell
    "msg"="trying to connect" "database"="OVN_Northbound" "endpoint"="tcp:172.18.0.4:6641"
```
This message indicates that the master cannot connect to the database. A
successful connection message will appear in the logs if the master has
connected to the database.

## Mitigation

Mitigation will depend on what was found in the diagnosis section. As a general
fix, you can try restarting the affected pods. Contact the incident response
team in your organisation if fixing the issue is not apparent.
