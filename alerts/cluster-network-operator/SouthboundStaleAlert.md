# SouthboundStale

## Meaning

This alert is triggered when northd or southbound database processes in a
specific availability domain are not functioning correctly or if connectivity
between them is broken. For OCP clusters at versions 4.13 or earlier, the
availability domain is the entire cluster. For OCP clusters at versions 4.14 or
later, the availability domain is a cluster node.

## Impact

Existing workloads may continue to have connectivity but any additional
workloads will not be provisioned. Any network policy changes will not be
implemented on existing workloads. For OCP clusters at versions 4.13 or earlier
the affected domain is the entire cluster. For OCP clusters at versions 4.14 or
later, the affected domain is only the specific node for which the alert was
fired.

### Fix alerts before continuing

Resolve any alerts that may cause this alert to fire: [Alert
hierarchy](./hierarchy/alerts-hierarchy.svg)

## Diagnosis

Investigate the health of the affected northd or southbound database processes
that run in the `northd` and `sbdb` containers repectively.

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

If northd is overloaded, there will be logs in the `northd` container along the
lines of `dropped x number of log messages due to excessive rate` or a message
that contains `(99% CPU usage)` or some other high percentage CPU usage.

You can also check the cpu usage of the container by logging into your openshift
cluster console. In the Observe section on the sidebar, click metrics, then
run this query `container_cpu_usage_seconds_total{pod="$PODNAME",
container="$CONTAINERNAME"}`

> NOTE: The checks below only apply for OCP clusters at versions 4.13 or earlier
> where ovn-northd remotely connects to a southbound database leader.

Ensure there is an sourthbound database Leader. To find the database leader, you
can run this script:
```shell
LEADER="not found"; for MASTER_POD in $(oc -n openshift-ovn-kubernetes get pods -l app=ovnkube-master -o=jsonpath='{.items[*].metadata.name}'); do RAFT_ROLE=$(oc exec -n openshift-ovn-kubernetes "${MASTER_POD}" -c sbdb -- bash -c "ovn-appctl -t /var/run/ovn/ovnsb_db.ctl cluster/status OVN_Southbound 2>&1 | grep \"^Role\""); if echo "${RAFT_ROLE}" | grep -q -i leader; then LEADER=$MASTER_POD; break; fi; done; echo "sbdb leader ${LEADER}"
```

If this does not work, exec into each sbdb container on the master pods:
```shell
oc exec -n openshift-ovn-kubernetes -it <podname> -c sbdb -- bash
```
and then run:
```shell
ovs-appctl -t /var/run/ovn/ovnsb_db.ctl cluster/status OVN_Southbound
```
You should see a role that will either say leader or follower.

A common cause of database leader issues is that one of the database servers is
unable to participate with the other raft peers due to mismatching cluster ids.
Due to this, they will be unable to elect a database leader.

## Mitigation

Mitigation will depend on what was found in the diagnosis section. As a general
fix, you can try restarting the affected pods. Contact the incident response
team in your organisation if fixing the issue is not apparent.
