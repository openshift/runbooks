# ClusterOperatorDegraded

## Meaning

The alert `ClusterOperatorDegraded` is triggered by the
[cluster-version-operator](https://github.com/openshift/cluster-version-operator)
(CVO) when a `ClusterOperator` is in the `degraded` state for a certain period.

An Operator reports a `Degraded` state when its current state does not
match the requested state over a period of time, which results in a lower
quality of service. The period of time varies by component, but a `Degraded`
state represents the persistent observation of a condition. A service state
might also be in an `Available` state even when degraded.

For example, a service might request three running pods, but one pod is in a
crash-loop state. In this case, the service is reported as `Available` but
`Degraded` because it might have a lower quality of service.

A component might also be reported as `Progressing` but not `Degraded` because
the change from one state to another does not persist over a long enough
time period to report a `Degraded` state.

A service does not report a `Degraded` state during a normal upgrade. A service
might report `Degraded` in response to a persistent infrastructure failure that
requires administrator intervention--for example, when a control plane host is
unhealthy and has to be replaced. An Operator reports `Degraded` state
if unexpected errors occur over a period of time.
## Impact

This alert indicates that an Operator has encountered an error preventing it
or its operand from working properly. The operand might still be available,
but its intent might not be fulfilled, and therefore an outage might occur.

## Diagnosis

The alert message indicates the Operator for which the alert triggered. The
Operator name is displayed under the `name` label, as shown in the following
example:

```text
 - alertname = ClusterOperatorDegraded
...
 - name = console
...
```

To troubleshoot the issue causing the alert to trigger, use any or all of
the following methods after logging into the cluster:

* Review the status of all Operators to discover if multiple Operators are
in a `Degraded` state:

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

How you resolve the issue causing the `Degraded` state of the Operator varies
depending on the Operator. If the alert is triggered during an upgrade, the
`Degraded` state might recover after some time has passed. If an Operator is
misconfigured, troubleshoot the error by reviewing information about
the Operator in the logs and fix the configuration based on your findings.
