# ThanosRuleRuleEvaluationLatencyHigh

## Meaning

The alert fires when Thanos Ruler misses rule evaluations due to slow rule
group processing. The alert is only about the user-defined recording and
alerting rules.

## Impact

Alerts will take longer time to get delivered.

## Diagnosis

Check the logs for the Thanos Ruler pods:

```console
oc -n openshift-user-workload-monitoring logs -l 'thanos-ruler=user-workload' \
-c thanos-ruler
```

Check the logs for the Thanos Querier pods:

```console
oc -n openshift-monitoring logs -l 'app.kubernetes.io/name=thanos-query' \
-c thanos-query
```

This alert will trigger when the rule evaluation takes more time than
the configured interval. It can indicate that your query backend
(i.e Thanos Querier) takes too much time to evaluate the query. If there is
even a single rule that is taking too long (more than the interval for that
group) to evaluate, this alert will fire. This might indicate other problems
like slow StoreAPIs or too complex query expression in rule.

## Mitigation

- The most likely scenario is a misconfiguration causing the user workload
  monitoring stack to overload Thanos Ruler with duplicate or otherwise
  erroneous alerts

- Audit the rule groups that fire the alert to identify expensive queries and
  consider splitting the rule groups into smaller groups if possible

- Verify if resource limits are set on the different monitoring components and
  whether they are throttled

- Check if any of the configured thanos-querier storeAPI endpoints
  have connectivity issues
