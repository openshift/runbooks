# AlertmanagerClusterFailedToSendAlerts

## Meaning

The alert `AlertmanagerClusterFailedToSendAlerts` is triggered when all
Alertmanager instances in a cluster have consistently failed to send
notifications to an integration.

## Impact

Some notifications are not delivered to the integration.

## Diagnosis

Review the logs of the `alertmanager-main` pods in the `openshift-monitoring`
namespace by running:

```console
$ oc -n openshift-monitoring logs -l 'alertmanager=main'
```

The following reasons might cause this alert to fire:

- The endpoint is not reachable.
- The endpoint's URL or credentials are misconfigured.

## Mitigation

How you resolve the problem causing the alert to fire depends on the particular
issue reported in the logs.
