# PrometheusPossibleNarrowSelectors

## Meaning

The `PrometheusPossibleNarrowSelectors` alert is triggered when PromQL queries
or metric relabel configurations use selectors that could be too restrictive
and might not take into account that values on the `le` label of classic
histograms or the `quantile` label of summaries are ingested as floats in OCP
4.19 and later.

## Impact

Metric relabeling configurations and PromQL queries using narrow selectors
on `quantile` or `le` labels can return incorrect results or stop working as
expected. This applies to queries run directly via the API, in the UI and
dashboards, or within recording rules and alerts.

## Diagnosis

Determine which component triggered the alert. It can originate from
Prometheus, Thanos Querier, or Thanos Ruler. Checking the alert message's
`job` label should help with this.

Also, note the alert message's `namespace` label: the namespace for
default cluster monitoring is `openshift-monitoring` and the namespace for
user workload monitoring is `openshift-user-workload-monitoring`.

## Mitigation

Proceed with the following steps to fix the issue:

1. Enable the `debug` log level on the component. See
[a guide to changing the log level] for monitoring components in OpenShift.
2. Review the new debug logs. The logs should reveal the problematic
selectors:

   ```shell
   $ NAMESPACE='<value of namespace label from alert>'
   $ NAME='<name of the component, see mapping below>'
   # job label (from the alert) -> NAME mapping:
   # prometheus-k8s -> prometheus
   # prometheus-user-workload -> prometheus
   # thanos-ruler -> thanos-ruler
   # thanos-querier -> thanos-query

   $ oc -n $NAMESPACE logs -l app.kubernetes.io/name=$NAME --tail=-1 \
   | grep 'narrow_matcher_label' | sort | uniq -c | sort -n
   ```

    The logs will point to one or both of the following problems:

   ### Narrow selector in query

    In this case, you should expect logs like the following:

    ```text
    ... msg="selector set to explicitly match an integer, but values could be
    floats" component=parser narrow_matcher_label=le integer=1
    matchers="[le=\"1\" __name__=\"foo_bucket\"]"
    # OR
    ... msg="selector set to explicitly match integers only, but values could be
    floats" component=parser narrow_matcher_label=le
    matchers="[le=~\"1|5|10\" __name__=\"foo_bucket\"]"
    ```

    The log should contain the problematic selector to help identify the
    query. We will focus on `le`, but similar logs should appear
    for the `quantile` label.

    For the first log, the selector should be adapted to the fact that in OCP
    4.19 and later, the value of `le` is ingested as `1.0`. If the query needs
    to cover data from both before and after the OCP >= 4.19 upgrade, the
    selector should be `{le=~"1(.0)?"}`. If the query only covers data after
    the upgrade, the selector should be `{le="1.0"}`.

    Even in cases where the query only covers data before the upgrade, it is
    recommended to use selectors that also account for float values.

    For Prometheus, enabling [query logging] can help you learn more about the query

   ### Narrow regex in metric relabel configuration

    In this case, you should expect logs like the following:

    ```text
    ... msg="relabel_config involves 'le' or 'quantile' labels, it may need to
    be adjusted to account for float values" component=relabel
    source_labels="__name__, le" regex=^(?s:foo_bucket;0.5|5|50)$
    ```

    The log should contain the problematic selector, which would help identify
    the concerned metric relabel configuration. We will focus on `le`, but
    similar logs should show up for the `quantile` label.

    In this case, because the `5` and `50` values no longer exist (assuming
    no relabeling is present to bring the integer values back, which is not
    recommended) the regex over `le` should be adapted to match `5.0` and
    `50.0` instead: `(0.5|5|50)(\.0)?`

3. After resolving the issue, disable the `debug` log level.

---

If you cannot resolve the issue, or if examination concluded that the alert
only fired for false positives, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.

[a guide to changing the log level]: https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/monitoring/configuring-core-platform-monitoring#setting-log-levels-for-monitoring-components_storing-and-recording-data

[query logging]: https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/monitoring/configuring-core-platform-monitoring#setting-query-log-file-for-prometheus_storing-and-recording-data
