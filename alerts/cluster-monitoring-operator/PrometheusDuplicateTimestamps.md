# PrometheusDuplicateTimestamps

## Meaning

The `PrometheusDuplicateTimestamps` alert is triggered when there is a constant
increase in dropped samples due to them having identical timestamps.

## Impact

Unwanted samples might be scraped, and desired samples might be dropped.
Consequently, queries may yield incorrect and unexpected results.

## Diagnosis

1. Determine whether the alert has triggered for the instance of Prometheus used
   for default cluster monitoring or for the instance that monitors user-defined
   projects by viewing the alert message's `namespace` label: the namespace for
   default cluster monitoring is `openshift-monitoring` and the namespace for
   user workload monitoring is `openshift-user-workload-monitoring`.

2. Review the logs for the affected Prometheus instance:

   ```shell
   $ NAMESPACE='<value of namespace label from alert>'

   $ oc -n $NAMESPACE logs -l 'app.kubernetes.io/name=prometheus' | \
   grep 'Error on ingesting samples with different value but same timestamp.*' \
   | sort | uniq -c | sort -n
   level=warn ... scrape_pool="the-scrape-pool" target="an-involved-target" \
   msg="Error on ingesting samples with different value but same timestamp"
   ```

   Warning logs similar to the one above should be present.

   In the case where targets are defined via a `ServiceMonitor` or `PodMonitor`,
   the `scrape_pool` label will be in the format
   `serviceMonitor/<namespace>/<service-monitor-name>/<endpoint-id>` or
   `podMonitor/<namespace>/<pod-monitor-name>/<endpoint-id>`.

## Mitigation

If the alert originates from the `openshift-monitoring` namespace, please open a
support case. If not, proceed based on the reviewed logs for the affected
Prometheus instance:

### The logs reveal one of the following issues

- The same target is defined in different scrape pools
- Distinct targets across different scrape pools are producing the same samples

This happens when the target is duplicated, that is, defined multiple times with
identical target labels.

Proceed with the following steps to fix the issue:

1. Use the logs to guide you to the place where the conflicting targets are defined.
2. Remove the duplicated targets or ensure that distinct targets are labeled uniquely.

### The logs only revolve around the same scrape pool

This might mean that the target exposes duplicated samples. In this scenario, even
if the samples have the same value, they are considered duplicates.

Proceed with the following steps to fix the issue:

1. Enable the `debug` log level on the Prometheus instance. See
[a guide to change log level] for monitoring components in OpenShift.
2. Review the new debug logs. The logs should reveal the problematic metrics:

   ```shell
   $ NAMESPACE='<value of namespace label from alert>'

   $ oc -n $NAMESPACE logs -l 'app.kubernetes.io/name=prometheus' | \
   grep 'Duplicate sample for timestamp.*' | sort | uniq -c | sort -n
   level=debug ... msg="Duplicate sample for timestamp" \
   series="a-concerned-series"
   ```

   In this case, you have to fix the metrics exposition in the broken targets to
   resolve the problem. Ensure they do not expose duplicated samples.
3. After resolving the issue, disable the `debug` log level.

[a guide to change log level]: https://docs.openshift.com/container-platform/latest/observability/monitoring/config-map-reference-for-the-cluster-monitoring-operator.html
