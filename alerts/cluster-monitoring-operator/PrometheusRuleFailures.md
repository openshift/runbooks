# PrometheusRuleFailures

## Meaning

The `PrometheusRuleFailures` alert triggers when there has been
a constant increase in failed evaluations of Prometheus rules for more
than 15 minutes.

## Impact

Recorded metrics and alerts might be missing or inaccurate.

## Diagnosis

1. Determine whether the alert has triggered for the instance of Prometheus used
for default cluster monitoring or for the instance that monitors user-defined
projects by viewing the alert message's `namespace` label: the namespace for
default cluster monitoring is `openshift-monitoring`; the namespace for user
workload monitoring is `openshift-user-workload-monitoring`.

1. Review the logs for the affected Prometheus instance:

    ```console
    $ NAMESPACE='<value of namespace label from alert>'

    $ oc -n $NAMESPACE logs -l 'app.kubernetes.io/name=prometheus' | \
    grep -o 'Evaluating rule failed.*' | sort | uniq -c | sort -n
    level=error ... msg="Evaluating rule failed." ...
    ```

Note that you can also evaluate the rule expression in the OpenShift web
console.

## Mitigation

If the logs indicate a syntax or other configuration error, troubleshoot the
issue:

- If a `PrometheusRule` is included with OpenShift, open a support case so that
a bug can be logged and the expression fixed.
- If a `PrometheusRule` is not included with OpenShift, then correct the
corresponding resource.
