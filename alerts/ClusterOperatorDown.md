# ClusterOperatorDown

## Meaning

The alert `ClusterOperatorDown` is fired by
[cluster-version-operator](https://github.com/openshift/cluster-version-operator)
(CVO) when a `ClusterOperator` is not in the `Available` state for a certain
period. An operand that is functional in the cluster is `Available`.

## Impact

There is an outage that needs to be checked at the earliest.

## Diagnosis

The alert would convey exactly which operator the alert is for. The message will
contain the name. For example:

```text
 - alertname = ClusterOperatorDown
...
 - name = console
...
```

First, log in to the cluster. Multiple operators could be down at the same time.
Check the status of all operators to know whether more are not `Available`:

```console
$ oc get clusteroperator
```

Typically, the status will give some hint about its state. Investigate further
for the operator that is not `Available` by using the following command:

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

The resolution steps would vary and depend on the particular operator that is in
consideration. If there is an upgrade going on, then this issue may resolve
itself after some time. Otherwise, try to find an error in the logs.
