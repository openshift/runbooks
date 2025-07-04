# AlertmanagerFailedReload

## Meaning

The alert `AlertmanagerFailedReload` is triggered when the Alertmanager instance
for the cluster monitoring stack has consistently failed to reload its
configuration for a certain period of time.

## Impact

Alerts for cluster components may not be delivered as expected.

## Diagnosis

1. Determine whether the alert has triggered for the instance of Alertmanager used
   for default cluster monitoring or for the instance that monitors user-defined
   projects by viewing the alert message's `namespace` label:
   - The namespace for default cluster monitoring is `openshift-monitoring`.
   - The namespace for user workload monitoring is `openshift-user-workload-monitoring`.

2. Check the logs for the Alertmanager pods:

```console
$ NAMESPACE='<value of namespace label from alert>'

$ oc -n $NAMESPACE logs -l 'app.kubernetes.io/name=alertmanager' --tail=-1 | \
grep 'Loading configuration file failed.*' \
| sort | uniq -c | sort -n
time=2025-04-18T07:28:00.274Z level=ERROR source=coordinator.go:117 msg="Loading configuration file failed" component=configuration file=/etc/alertmanager/config_out/alertmanager.env.yaml err="address smtp.gmail.com: missing port in address"
```

## Mitigation

The resolution depends on the particular issue reported in the logs.
