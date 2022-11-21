# ThanosRuleRuleEvaluationLatencyHigh

## Meaning

The `ThanosRuleRuleEvaluationLatencyHigh` alert triggers when Thanos Ruler
misses rule evaluations due to slow rule group processing.

This alert triggers only for user-defined recording and alerting rules.

## Impact

The delivery of alerts will be delayed.

## Diagnosis

Review the logs for the Thanos Ruler pods:

```console
oc -n openshift-user-workload-monitoring logs -l 'thanos-ruler=user-workload' \
-c thanos-ruler
```

Review the logs for the Thanos Querier pods:

```console
oc -n openshift-monitoring logs -l 'app.kubernetes.io/name=thanos-query' \
-c thanos-query
```

This alert triggers when rule evaluation takes longer than the configured
interval.

If the alert triggers, it might indicate that Thanos Querier is taking too much
time to evaluate the query. This alert will trigger if rule evaluation for even
a single rule is taking too long--that is, longer than the interval for that
group.

If the alert triggers, it might also mean that other problems exist, such as
StoreAPIs are responding slowly or a query expression in a rule is too complex.

## Mitigation

- Check for a misconfiguration that causes the user workload monitoring stack
  to overload Thanos Ruler with duplicate or otherwise erroneous alerts.

- Audit the rule groups that fire the alert to identify expensive queries and
  consider splitting these rule groups into smaller groups if possible.

- Verify whether resource limits are set on any monitoring components
  and whether any components are throttled.

- Check whether any of the configured `thanos-querier` storeAPI endpoints
  have connectivity issues.
