# ThanosRuleQueueIsDroppingAlerts

## Meaning

The `ThanosRuleQueueIsDroppingAlerts` alert triggers when the Thanos Ruler queue
is found to be dropping alerting events.

The Thanos Ruler component is deployed only when user-defined monitoring is
enabled. The component enables alerting rules to be deployed as part of
user-defined monitoring. These rules can query the Prometheus instance
responsible for core cluster components and also the Prometheus instance used
for user-defined monitoring.

## Impact

Alerts for user workloads might not be delivered.

## Diagnosis

Review the logs for the Thanos Ruler pods:

```console
$ oc -n openshift-user-workload-monitoring logs -l 'thanos-ruler=user-workload'
...
level=warn ... msg="Alert notification queue full, dropping alerts" numDropped=100
level=warn ... msg="Alert batch larger than queue capacity, dropping alerts" numDropped=100
```

If this alert triggers, it is likely that the user-defined monitoring stack is
firing an extremely large number of alerts. Log into the OpenShift web console
and review the active alerts.

## Mitigation

The default queue capacity for Thanos Ruler is quite high at 10,000 items,
which means that the most likely cause of this issue is a misconfiguration that
causes the user-defined monitoring stack to overload Thanos Ruler with
duplicate or otherwise erroneous alerts. Review all active alerts in the
OpenShift web console and correct any misconfigurations. You can also consider
grouping alerts to mitigate this issue.

