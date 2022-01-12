# AlertmanagerClusterFailedToSendAlerts

## Meaning

The alert `AlertmanagerClusterFailedToSendAlerts` is triggered when **all**
Alertmanager instances in a cluster has consistently failed to send
notifications to an integration.

## Impact

A fraction of notifications to the integration are not delivered.

## Diagnosis

Check the logs of the `alertmanager-main` pods in the `openshift-monitoring`
namespace:


```console
$ oc -n openshift-monitoring logs -l 'alertmanager=main'
```

Below is a non-exhaustive list of reasons for this alert to be fired.
- Unreachable endpoint
- Misconfigured address or credentials

## Mitigation

The resolution depends on the particular issue reported in the logs.
