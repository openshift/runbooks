# KubeletDown

## Meaning

The `KubeletDown` alert is triggered when the monitoring system has not been
able to reach any of the Kubelets in a cluster for more than 15 minutes.

## Impact

This alert represents a critical threat to the cluster's stability. Excluding
the possibility of a network issue preventing the monitoring system from
scraping Kubelet metrics, multiple nodes in the cluster are likely unable to
respond to configuration changes for pods and other resources, and some
debugging tools are likely not to be functional, such as `oc exec` and
`oc logs`.

## Diagnosis

Review the status of the nodes and check for recent events on `Node` or other
resources:

```console
$ oc get nodes
$ oc describe node $NODE_NAME
$ oc get events --field-selector 'involvedObject.kind=Node'
$ oc get events
```

If you have SSH access to the nodes, use this access to review the logs for the
Kubelet:

```console
$ journalctl -b -f -u kubelet.service
```

## Mitigation

The mitigation for this alert depends on the issue causing the Kubelets to
become unresponsive. You can begin by checking for general networking issues or
for node-level configuration issues.
