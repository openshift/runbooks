# AlertmanagerFailedReload

## Meaning

The alert `AlertmanagerFailedReload` is triggered when the Alertmanager instance
for the cluster monitoring stack has consistently failed to reload its
configuration for a certain period.

## Impact

Alerts for cluster components may not be delivered as expected.

## Diagnosis

Check the logs for the `alertmanager-main` pods in the `openshift-monitoring`
namespace:

```console
$ oc -n openshift-monitoring logs -l 'alertmanager=main'
```

## Mitigation

The resolution depends on the particular issue reported in the logs.
