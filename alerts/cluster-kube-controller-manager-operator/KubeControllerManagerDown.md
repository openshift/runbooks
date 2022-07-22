# KubeControllerManagerDown

## Meaning

[This alert][KubeControllerManagerDown] is fired when KubeControllerManager
has disappeared from Prometheus target discovery.
This means there is no running or properly functioning
instance of [kube-controller-manager][kube-controller-manager].

## Impact

Many features stop working when kube-controller-manager is down.

This includes:
- workload controllers (Deployment, ReplicaSet, DaemonSet, ...)
- resource quotas
- pod disruption budgets
- garbage collection
- certificate signing requests
- service accounts, tokens
- storage
- nodes' statuses and taints
- SCCs
- and more...


## Diagnosis

To see an operator status of kube-controller-manager.

```console
$ oc get clusteroperators.config.openshift.io kube-controller-manager
```

Take a look at `KubeControllerManager`'s `.status.conditions` and
also see what is the current state of each instance of kube-controller-manager
on each node in `.status.nodeStatuses`.

```console
$ oc get -o yaml kubecontrollermanagers.operator.openshift.io cluster
```

See operator events.

```console
$ oc get events -n openshift-kube-controller-manager-operator
```

Look at the operator pod and inspect its logs.

```console
$ oc get pods -n openshift-kube-controller-manager-operator
$ oc logs -n openshift-kube-controller-manager-operator $POD_NAME
```

You can do the same with kube-controller-manager.

See kube-controller-manager events.

```console
$ oc get events -n openshift-kube-controller-manager
```

Look at kube-controller-manager pods and inspect their logs.

```console
$ oc get pods -n openshift-kube-controller-manager
$ oc logs -n openshift-kube-controller-manager $POD_NAME
```


## Mitigation

The resolution depends on the particular issue reported in the statuses,
events and logs.


[KubeControllerManagerDown]: https://github.com/openshift/cluster-kube-controller-manager-operator/blob/20179ecfa3b8c5e766a21c98107f45b84196b914/manifests/0000_90_kube-controller-manager-operator_05_alerts.yaml#L14-L23
[kube-controller-manager]: https://kubernetes.io/docs/reference/command-line-tools-reference/kube-controller-manager/