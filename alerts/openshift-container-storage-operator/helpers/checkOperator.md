# OCS Operator diagnosis

Checking on the OCS operator status involves checking the operator subscription
status and the operator pod health.

## OCS Operator Subscription Health

Check the ocs-operator subscription status

```bash
  oc get sub $(oc get pods -n openshift-storage | grep -v ocs-operator) -n openshift-storage  -o json | jq .status.conditions
```

Like all operators, the status conditions types are:

**CatalogSourcesUnhealthy, InstallPlanMissing, InstallPlanPending,
InstallPlanFailed**

The status for each type should be False. For example:

```bash
    [
      {
        "lastTransitionTime": "2021-01-26T19:21:37Z",
        "message": "all available catalogsources are healthy",
        "reason": "AllCatalogSourcesHealthy",
        "status": "False",
        "type": "CatalogSourcesUnhealthy"
      }
    ]
```

The output above shows a false status for type CatalogSourcesUnHealthly,
meaning the catalog sources are healthy.

## OCS Operator Pod Health

Check the OCS operator pod status to see if there is an OCS operator upgrading
in progress.

WIP: Find specific status for upgrade (pending?)

To find and view the status of the OCS operator:

```bash
  oc get pod -n openshift-storage | grep ocs-operator OCSOP=$(oc get pod -n openshift-storage -o custom-columns=POD:.metadata.name --no-headers | grep cs-operator)
  echo $OCSOP
  oc get pod/${OCSOP} -n openshift-storage
  oc describe pod/${OCSOP} -n openshift-storage
```

If you determine the OCS operator is in progress, please be patient,
wait 5 minutes and this alert should resolve itself.

If you have waited or see a different error status condition,
please continue troubleshooting.
