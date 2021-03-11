# KubePodNotReady

## Meaning

[This alert][KubePodNotReady] is fired in case there are pods which have not
been in `Ready` state for a certain period. This can have different reasons as
described [here][PodLifecycle]: In case of the pod is `Running` but not `Ready`,
the `Readiness` probe is failing. An application-specific error may prevent the
pod from being attached to a service. In case the pod remains in `Pending`, it
can not be deployed to particular namespaces and nodes.

## Impact

This pod is not functional and doesn't receive any traffic. Depending on how
many functional replicas are still available, the impact differs.

## Diagnosis

The notification details should list the pod that's not ready (and the pod's
namespace). E.g.:

```text
 - alertname = KubePodNotReady
...
 - namespace = openshift-logging
 - pod = elasticsearch-cdm-u1gqqbu6-2-868ddd4b45-w224d
...
```

Start by checking the status of the pod:

```console
$ oc get pod -n $NAMESPACE $POD
```

If the pod state is in `Running`, you can check its logs:

```console
$ oc logs -n $NAMESPACE $POD
```

Be aware there may be multiple containers in the pod. A check of all their logs
may be required. If the pod isn't running (for instance, if it's stuck in
ContainerCreating), then try to find out why.

## Mitigation

Try to find the issue in the logs before deleting it.

[KubePodNotReady]: https://github.com/openshift/cluster-monitoring-operator/blob/aefc8fc5fc61c943dc1ca24b8c151940ae5f8f1c/assets/control-plane/prometheus-rule.yaml#L26-L41
[PodLifecycle]: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/