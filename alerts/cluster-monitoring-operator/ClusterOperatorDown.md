# ClusterOperatorDown

## Meaning

The alert `ClusterOperatorDown` is triggered by
[cluster-version-operator](https://github.com/openshift/cluster-version-operator)
(CVO) when a `ClusterOperator` is not in the `Available` state for a certain
period of time. An operand is `Available` when it is functional in the cluster.

## Impact

This alert indicates that an outage has occurred in your cluster. Investigate
the issue as soon as possible.

## Diagnosis

The alert message provides the name of the Operator that triggered the alert,
as shown in the following example:

```text
 - alertname = ClusterOperatorDown
...
 - name = console
...
```

To troubleshoot the issue causing the alert to trigger, use any or all of
the following methods after logging into the cluster:

* Review the status of all Operators to discover if multiple Operators are
down:

    ```console
    $ oc get clusteroperator
    ```

* Review information about the current status of the Operator:

    ```console
    $ oc get clusteroperator $CLUSTEROPERATOR -ojson | jq .status.conditions
    ```

* Review the associated resources for the Operator:

    ```console
    $ oc get clusteroperator $CLUSTEROPERATOR -ojson | jq .status.relatedObjects
    ```

* Review the logs and other artifacts for the Operator. For example, you can
collect the logs of a specific Operator and store them in a local directory
named `out`:

    ```console
    $ oc adm inspect clusteroperator/$CLUSTEROPERATOR --dest-dir=out
    ```

## Mitigation

How you resolve the issue causing the issue varies depending on the Operator.
If the alert is triggered during an upgrade, the issue might resolve after some
time has passed. Otherwise, troubleshoot the error by reviewing information
about the Operator in the logs and fix the configuration based on your findings.
