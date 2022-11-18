# KubeNodeNotReady

## Meaning

The `KubeNodeNotReady` alert triggers when a node is not in a `Ready` state
over a certain period of time. If this alert triggers, the node cannot host any
new pods, as described [here][KubeNode].

## Impact

The issue that triggers this alert degrades the performance of the cluster
deployments. The severity of the degradation depends on the overall workload
and the type of node.

## Diagnosis

The alert notification message includes the affected node, as shown in the
following example:

```txt
 - alertname = KubeNodeNotReady
...
 - node = node1.example.com
...
```

* Log in to the cluster. Review the status of the node indicated in the alert
message:

    ```console
    $ oc get node $NODE -o yaml
    ```

    The output of this command describes why the node is not ready. For example
    network issues could be causing timeouts when trying to reach the API or
    kubelet.

* Check the machine for the node:

    ```console
    $ oc get -n openshift-machine-api machine $NODE -o yaml
    ```

* Check the events for the machine API:

    ```console
    $ oc get -n openshift-machine-api events
    ```

    If the machine API is not able to replace the node, the machine status and
    events list will provide the details.

## Mitigation

After you resolve the problem that prevented the machine API from replacing the
node, the instance is terminated and replaced by the machine API, but only if
`MachineHealthChecks` are enabled for the nodes. Otherwise, a manual restart is
required.

[KubeNode]: https://kubernetes.io/docs/concepts/architecture/nodes/#condition
