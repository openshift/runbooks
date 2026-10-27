# KubeletInstanceUnreachable

## Meaning

The `KubeletInstanceUnreachable` alert is triggered when the monitoring system
has failed to scrape the `/metrics` endpoint of a single Kubelet instance for
more than 15 minutes, either because the monitoring system can't connect to the
Kubelet or because the Kubelet fails to respond in time.

Unlike `KubeletDown`, which fires when Prometheus fails to scrape all Kubelets
across the cluster, this alert targets an individual Kubelet instance that
remains a known scrape target but is unreachable.

## Impact

The affected node is likely unable to respond to configuration changes for pods
and other resources, and some debugging tools are likely not to be functional
for workloads on that node, such as `oc exec` and `oc logs`. Metrics from the
node are also unavailable, which can affect autoscaling and monitoring of the
workloads scheduled there.

## Diagnosis

Identify the affected node. The alert's `node` label points to the Kubelet node
that is unreachable. Review the status of the node and check for recent events
on `Node` or other resources:

```console
$ oc get nodes
$ oc describe node $NODE_NAME
$ oc get events --field-selector 'involvedObject.kind=Node'
$ oc get events
```

You can [access cluster node journal logs][cluster-node-journal-logs] to review
the logs for the Kubelet. If the Kubelet is functional, you can use:

```console
$ oc adm node-logs $NODE_NAME -u kubelet
```

See `oc adm node-logs --help` for alternative ways to select nodes and filter results.

If the Kubelet is not functional and you have SSH access to the node,
use this access to review the logs for the Kubelet:

```console
$ journalctl -b -f -u kubelet.service
```

## Mitigation

The mitigation for this alert depends on the issue causing the Kubelet to become
unreachable. You can begin by checking for networking issues between the
monitoring system and the node, or for node-level configuration issues. If the
node itself is unhealthy or has become unreachable, investigate the node's
health and recover it as appropriate.

[cluster-node-journal-logs]: https://docs.openshift.com/container-platform/latest/support/gathering-cluster-data.html#querying-cluster-node-journal-logs_gathering-cluster-data
