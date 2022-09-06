# OVNKubernetesControllerDisconnectedSouthboundDatabase

## Meaning

The `OVNKubernetesControllerDisconnectedSouthboundDatabase` alert is triggered
when the OVN controller is not connected to OVN southbound database
for more than 5 minutes.

## Impact

Networking control plane is degraded on the node. Existing workloads
on the node may continue to have connectivity but any networking configuration
update will not be applied.

## Diagnosis

### Fix alerts before continuing

Check to ensure the following alerts are not firing and resolved before
continuing as they may cause this alert to fire:

- [NoRunningOvnMaster](./NoRunningOvnMaster.md)

### OVN-kubernetes master pods

Find ovnkube-master pods.

```shell
pods=($(oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-master \
  -o jsonpath={..metadata.name}))
```

Check the container statuses in ovnkube-master pods.

```shell
for pod in ${pods}
do
  echo "${pod}:\n"
  oc describe pods -n openshift-ovn-kubernetes ${pod}
done
```

Check the sbdb container logs in ovnkube-master pods.

```shell
for pod in ${pods}
do
  echo "${pod}:\n"
  oc logs -n openshift-ovn-kubernetes ${pod} sbdb
done
```

### OVN-kubernetes node pod

Find the ovnkube-node pod running on the affected node.

```shell
pod=$(oc get pod -n openshift-ovn-kubernetes -l app=ovnkube-node \
  -o jsonpath={..metadata.name} --field-selector spec.nodeName=<node>)
```

Check the container statuses in the ovnkube-node pod.

```shell
oc describe po -n openshift-ovn-kubernetes ${pod}
```

Check the southbound database connection status.

```shell
oc exec -n openshift-ovn-kubernetes ${pod} -c ovn-controller -- ovn-appctl connection-status
```

Check the logs of the ovnkube-node container looking for OVN controller error logs.

```shell
oc logs -n openshift-ovn-kubernetes ${pod} -c ovnkube-node
```

Check the logs of the ovn-controller container.

```shell
oc logs -n openshift-ovn-kubernetes ${pod} -c ovn-controller
```

Using tcpdump on the affected node verify the traffic flow to southbound database.

```shell
oc debug node/<node> -- tcpdump -i <primary_interface> tcp and port 9642
```

## Mitigation

Mitigation will depend on what was found in the diagnosis section.

If there is no traffic flowing between southbound database and the ovn-controller
it can mean that there are underlying issues in the infrastructure.