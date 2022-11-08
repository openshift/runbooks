# AlertmanagerFailedToSendAlerts

## Meaning

The alert `AlertmanagerFailedToSendAlerts` is triggered when any of the
Alertmanager instances in the cluster monitoring stack has repeatedly
failed to send notifications to an integration.

## Impact

Some alert notifications are not delivered.

## Diagnosis

Review the logs for the pod in the namespace indicated in the alert message.

For example, the following sample alert message refers to the
`alertmanager-main-1` pod in the `openshift-monitoring` namespace:

> Alertmanager openshift-monitoring/alertmanager-main-1 failed to send 75%
> of notifications to webhook.

You can review the logs for the `alertmanager-main-1` pod in the
`openshift-monitoring` namespace by running the following command:

```console
$ oc -n openshift-monitoring logs alertmanager-main-1
```

## Mitigation

The resolution depends on the particular issue reported in the logs.
