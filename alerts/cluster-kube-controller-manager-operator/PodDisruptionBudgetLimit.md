# PodDisruptionBudgetLimit

## Meaning

[This alert][PodDisruptionBudgetLimit] is fired when the pod disruption budget
is below the minimum disruptions allowed level and is not satisfied.
This level is defined by `.spec.minAvailable` or `.spec.maxUnavailable` in
the `PodDisruptionBudget` object.
The number of current healthy pods is less than the desired healthy pods.


## Impact

The application protected by the pod disruption budget has an insufficient
amount of pods.
This means the application is either not running at all or running with
suboptimal number of pods.
This can have various impact depending on the type and importance of such
application.

Standard workloads should have at least one pod more than is desired to support
[API-initiated eviction][APIEviction]. Workloads that are below the minimum
disruption allowed level violate this and could block node drain.
This is important for node maintenance and cluster upgrades.


## Diagnosis

Discover the pod disruption budgets that are triggering the alert.

```console
max by(namespace, poddisruptionbudget) (
    kube_poddisruptionbudget_status_current_healthy < kube_poddisruptionbudget_status_desired_healthy
)
```

Look at the [pod disruption budget][SpecifyingPDB] detail to see how many pods
are healthy and how many are desired.

```console
$ oc get poddisruptionbudgets -o yaml -n $NS $PDB_NAME
```

Look at events for reasons why your pods are not healthy.

```console
$ oc get events -n $NS
```


## Mitigation

Get the selector from the pod disruption budget YAML and debug the pods
of the application to find a reason why they are not healthy.

```console
$ oc get pods -n $NS --selector="app=myapp"
```

In general ensure you have enough resource for running your pods.
You can also take a look at potential reasons
for [pod disruptions][PodDisruptions].

Finally, take a look at the owner references of these pods to either change
the pod YAML or number of replicas in the parent workload resource to fix
the issue.


[PodDisruptionBudgetLimit]: https://github.com/openshift/cluster-kube-controller-manager-operator/blob/20179ecfa3b8c5e766a21c98107f45b84196b914/manifests/0000_90_kube-controller-manager-operator_05_alerts.yaml#L33-L41
[PodDisruptions]: https://kubernetes.io/docs/concepts/workloads/pods/disruptions/
[SpecifyingPDB]: https://kubernetes.io/docs/tasks/run-application/configure-pdb/
[APIEviction]: https://kubernetes.io/docs/concepts/scheduling-eviction/api-eviction/
