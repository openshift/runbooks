# ClusterOperatorDegraded

## Meaning

The alert `ClusterOperatorDegraded` is fired by
[cluster-version-operator](https://github.com/openshift/cluster-version-operator)(CVO)
when a `ClusterOperator` is in `degraded` state for a certain period. An
operator reports `Degraded` when its current state does not match its desired
state over a period resulting in a lower quality of service. The time may vary
by component, but a `Degraded` state represents the persistent observation of a
condition. A service state may be `Available` even when degraded. For example,
your service may desire three running pods, but one pod is in a crash-looping.
The service is `Available` but `Degraded` because it may have a lower quality of
service. A component may be `Progressing` but not `Degraded` because the
transition from one state to another does not persist over a long enough period
to report `Degraded`. A service should not report `Degraded` during a normal
upgrade. A service may report `Degraded` in response to a persistent
infrastructure failure that requires administrator intervention. For example,
when a control plane host is unhealthy and has to be replaced. An operator
should report `Degraded` if unexpected errors occur over a period, but the
expectation is that all unexpected errors are handled as operators mature.

## Impact

An operator has encountered an error that is preventing it or its operand from
working properly. The operand may still be available, but its intent may not be
fulfilled. If this is true, it means that the operand is at risk of an outage or
improper configuration.

## Diagnosis

The alert would convey exactly which operator the alert was fired for. The
operator name will be displayed under the `name` label. For example:

```text
 - alertname = ClusterOperatorDegraded
...
 - name = console
...
```

First, log in to the cluster. Multiple operators could be degraded at the same
time. Check the status of all operators to know whether there are more
`Degraded`:

```console
$ oc get clusteroperator
```

Typically, the status will give some hint about the operator state.

```console
$ oc get clusteroperator $CLUSTEROPERATOR -ojson | jq .status.conditions
```

Further on, if you would like to go through the associated resources for that
particular operator, you can use the command:

```console
$ oc get clusteroperator $CLUSTEROPERATOR -ojson | jq .status.relatedObjects
```

Collect logs and artifacts for a given operator. As an example, you can collect
the logs of a specific operator and store them in a local directory named `out`
by using the following command:

```console
$ oc adm inspect clusteroperator/$CLUSTEROPERATOR --dest-dir=out
```

## Mitigation

The resolution steps would vary and depend on the particular `ClusterOperator`
that is in consideration. If there is an upgrade going on, then the `Degraded`
state may recover itself in some time. If a `ClusterOperator` is misconfigured,
then try to find the error in the collected logs and fix the configuration.
