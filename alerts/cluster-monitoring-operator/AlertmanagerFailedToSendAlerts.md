# AlertmanagerFailedToSendAlerts

## Meaning

The alert `AlertmanagerFailedToSendAlerts` is triggered when any of the
Alertmanager instances of the cluster monitoring stack has consistently
failed to send notifications to an integration.

## Impact

Some of the notifcations are not delivered.

## Diagnosis

Check the logs of the `alertmanager-main` pod in `openshift-monitoring` namespace
in the alert description. E.g. If the Alert description says

> Alertmanager openshift-monitoring/alertmanager-main-1 failed to send 75%
> of notifications to webhook.

```console
$ oc -n openshift-monitoring logs alertmanager-main-1
```

## Mitigation

The resolution depends on the particular issue reported in the logs.
