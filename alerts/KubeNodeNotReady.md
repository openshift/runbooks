# KubeNodeNotReady

## Meaning

[This alert][KubeNodeNotReady] is fired when a Kubernetes node is not in `Ready`
state for a certain period. In this case, the node is not able to host any new
pods as described [here][KubeNode].

## Impact

The performance of the cluster deployments is affected, depending on the overall
workload and the type of the node.

## Diagnosis

The notification details should list the node that's not ready. For Example:

```txt
 - alertname = KubeNodeNotReady
...
 - node = node1.example.com
...
```

Login to the cluster. Check the status of that node:

```console
$ oc get node $NODE -o yaml
```

The output should describe why the node isn't ready (e.g.: timeouts reaching the
API or kubelet)
Check the machine for the node:

```console
$ oc get -n openshift-machine-api machine $NODE -o yaml
```

and the events for the machine API:

```console
$ oc get -n openshift-machine-api events
```

If the machine API is not able to replace the node, the machine status and
events should detail why.

## Mitigation

Once the problem was resolved that prevented the machine API from replacing the
node, the instance should be terminated and replaced by the machine API.

[KubeNode]: https://kubernetes.io/docs/concepts/architecture/nodes/#condition
[KubeNodeNotReady]: https://github.com/openshift/cluster-monitoring-operator/blob/aefc8fc5fc61c943dc1ca24b8c151940ae5f8f1c/assets/control-plane/prometheus-rule.yaml#L482-L490
