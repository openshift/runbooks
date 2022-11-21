# KubePodNotReady

## Meaning

The `KubePodNotReady` alert triggers when a pod has not been in a
`Ready` state for a certain time period. This issue can occur for different
reasons, as described [in the Kubernetes documentation][PodLifecycle]. When a
pod has a status of `Running` but is not in a `Ready` state, the `Readiness`
probe is failing. For example, an application-specific error might be
preventing the pod from being attached to a service. When a pod remains in
`Pending` state, it cannot be deployed to particular namespaces and nodes.

## Impact

The affected pod is not functional and does not receive any traffic. Depending
on how many functional replicas are still available, the severity of the impact
differs.

## Diagnosis

The alert notification message lists the pod that is not ready and the
namespace in which the pod is located, as shown in the following example alert
message:

```text
 - alertname = KubePodNotReady
...
 - namespace = openshift-logging
 - pod = elasticsearch-cdm-u1gqqbu6-2-868ddd4b45-w224d
...
```

To diagnose the cause of the issue, start by reviewing the status of the
affected pod:

```console
$ oc get pod -n $NAMESPACE $POD
```

If the pod is in a `Running` state, review the logs for the pod:

```console
$ oc logs -n $NAMESPACE $POD
```

Be aware that there might be multiple containers in the pod. If so, review the
logs for all of these containers. If the pod is not in a `Running` state--for
instance, if it is stuck in a `ContainerCreating` state--try to find out why.

## Mitigation

The steps you take to fix the issue will depend on the cause that you found when
you examined the logs.

[PodLifecycle]: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/