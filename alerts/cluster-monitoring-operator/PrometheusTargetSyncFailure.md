# PrometheusTargetSyncFailure

## Meaning

The `PrometheusTargetSyncFailure` alert triggers when at least one
Prometheus instance has consistently failed to sync its configuration.

## Impact

Metrics and alerts might be missing or inaccurate.

## Diagnosis

1. Determine whether the alert has triggered for the instance of Prometheus used
for default cluster monitoring or for the instance that monitors user-defined
projects by viewing the alert message's `namespace` label: the namespace for
default cluster monitoring is `openshift-monitoring`; the namespace for user
workload monitoring is `openshift-user-workload-monitoring`.

1. Review the logs for the affected Prometheus instance:

    ```console
    $ NAMESPACE='<value of namespace label from alert>'

    $ oc -n $NAMESPACE logs -l 'app.kubernetes.io/name=prometheus'
    level=error ... msg="Creating target failed" ...
    ```

## Mitigation

If the logs indicate a syntax or other configuration error, correct the
corresponding `ServiceMonitor`, `PodMonitor`, `Probe`, or other configuration
resource.
