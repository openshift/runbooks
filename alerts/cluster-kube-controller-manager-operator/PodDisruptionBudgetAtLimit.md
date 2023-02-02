# PodDisruptionBudgetAtLimit

## Meaning

[This alert][PodDisruptionBudgetAtLimit] is fired when the pod disruption
budget is at the minimum disruptions allowed level.
This level is defined by `.spec.minAvailable` or `.spec.maxUnavailable` in
the `PodDisruptionBudget` object.
The number of current healthy pods is equal to the desired healthy pods.

It does not fire when the number of expected pods is 0.


## Impact

the application protected by the pod disruption budget has a sufficient amount
of pods, but is at risk of getting disrupted.

Standard workloads should have at least one pod more than is desired to support
[API-initiated eviction][APIEviction]. Workloads that are at the minimum
disruption allowed level violate this and could block node drain.
This is important for node maintenance and cluster upgrades.

## Diagnosis

Discover the pod disruption budgets that are triggering the alert.

```console
max by(namespace, poddisruptionbudget) (
    kube_poddisruptionbudget_status_current_healthy == kube_poddisruptionbudget_status_desired_healthy
      and on (namespace, poddisruptionbudget) kube_poddisruptionbudget_status_expected_pods > 0
)
```

Look at the [pod disruption budget][SpecifyingPDB] detail.

```console
$ oc get poddisruptionbudgets -o yaml -n $NS $PDB_NAME
```

Look at events for reasons why your extra pods might not be healthy.

```console
$ oc get events -n $NS
```


## Mitigation

Get the selector from the pod disruption budget YAML and see if there are
any extra pods of this application that are not healthy.
You can debug these pods to find a reason why.

```console
$ oc get pods -n $NS --selector="app=myapp"
```

In general ensure you have enough resource for running your pods.
You can also take a look at potential reasons
for [pod disruptions][PodDisruptions].

Finally, take a look at the owner references of these pods to either change
the pod YAML or number of replicas in the parent workload resource to fix
the issue.


[PodDisruptionBudgetAtLimit]: https://github.com/openshift/cluster-kube-controller-manager-operator/blob/20179ecfa3b8c5e766a21c98107f45b84196b914/manifests/0000_90_kube-controller-manager-operator_05_alerts.yaml#L24-L32
[PodDisruptions]: https://kubernetes.io/docs/concepts/workloads/pods/disruptions/
[SpecifyingPDB]: https://kubernetes.io/docs/tasks/run-application/configure-pdb/
[APIEviction]: https://kubernetes.io/docs/concepts/scheduling-eviction/api-eviction/
