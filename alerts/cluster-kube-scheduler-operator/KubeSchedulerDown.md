# KubeSchedulerDown

## Meaning

[This alert][KubeSchedulerDown] is fired when the kube-scheduler pods
has disappeared from Prometheus target discovery.
This means there is no running or properly functioning
instance of [kube-scheduler][kube-scheduler].

## Impact

The cluster is unable to schedule new pods. Newly created pods will stay
in the Pending state until the kube-scheduler is running properly again.

## Diagnosis

To see an operator status of kube-scheduler.

```console
$ oc get clusteroperators.config.openshift.io kube-scheduler
```

Take a look at `KubeScheduler`'s `.status.conditions` and
also see what is the current state of each instance of kube-scheduler
on each node in `.status.nodeStatuses`.

```console
$ oc get -o yaml kubeschedulers.operator.openshift.io cluster
```

See operator events.

```console
$ oc get events -n openshift-kube-scheduler-operator
```

Look at the operator pod and inspect its logs.

```console
$ oc get pods -n openshift-kube-scheduler-operator
$ oc logs -n openshift-kube-scheduler-operator $POD_NAME
```

You can do the same with kube-scheduler.

See kube-scheduler events.

```console
$ oc get events -n openshift-kube-scheduler
```

Look at kube-scheduler pods and inspect their logs.

```console
$ oc get pods -n openshift-kube-scheduler
$ oc logs -n openshift-kube-scheduler $POD_NAME
```


## Mitigation

Since this alert indicates that the control-plane is not functioning
and is unlikely to self-heal, the manual intervention required depend
heavily on the exact scenario.
The status and log insights gained from the Diagnosis section provide
the needed insight. The kube-scheduler pods need to advance
to the Running state in order to mitigate this alert.

[KubeSchedulerDown]: https://github.com/openshift/cluster-kube-scheduler-operator/blob/98d1828fb44fb78abf2e090825404a98cd8a4e22/manifests/0000_90_kube-scheduler-operator_03_servicemonitor.yaml#L47-L56
[kube-scheduler]: https://kubernetes.io/docs/reference/command-line-tools-reference/kube-scheduler/
