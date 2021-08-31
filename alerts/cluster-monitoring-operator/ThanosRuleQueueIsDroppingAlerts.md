# ThanosRuleQueueIsDroppingAlerts

## Meaning

The [Thanos Ruler][1] is deployed when user workload monitoring is enabled. It
allows alerting rules deployed as part of user workload monitoring to query both
the the Prometheus instance responsible for cluster components as well the user
workload Prometheus instance. This alert is triggered when the Thanos Ruler
queue is found to be dropping alerting events.

## Impact

Alerts for user workloads may not be delivered.

## Diagnosis

Check the logs for the Thanos Rules pods:

```console
$ oc -n openshift-user-workload-monitoring logs -l 'thanos-ruler=user-workload'
...
level=warn ... msg="Alert notification queue full, dropping alerts" numDropped=100
level=warn ... msg="Alert batch larger than queue capacity, dropping alerts" numDropped=100
```

If this alert has triggered, you likely have an extremely large number of alerts
flowing from the user workload monitoring stack. Log into the OpenShift web
console and check the active alerts.

## Mitigation

The default queue capacity for Thanos Ruler is quite high at 10,000 items,
meaning the most likely scenario is a misconfiguration causing the user workload
monitoring stack to overload Thanos Ruler with duplicate or otherwise erroneous
alerts. Check the active alerts in the OpenShift web console, and correct any
misconfigurations or consider grouping alerts.

[1]: https://thanos.io/v0.22/components/rule.md
