# ClusterLogForwarderNotReady

## Meaning

The `ClusterLogForwarderNotReady` alert indicates that a ClusterLogForwarder
has been in a not ready state for more than 1 minute. This typically indicates
a validation error in the ClusterLogForwarder spec.

## Impact

When this alert fires, the following impacts may occur:

- **Log data loss:** Logs are not being collected or forwarded to their
configured destinations while the ClusterLogForwarder remains not ready.
- **No log visibility:** Downstream systems (e.g., LokiStack, Elasticsearch,
CloudWatch) stop receiving new log data, impacting dashboards, alerting, and
troubleshooting capabilities.

## Diagnosis

1. **Check the ClusterLogForwarder status conditions:**

   ```bash
   oc get clusterlogforwarder -A
   oc describe clusterlogforwarder <NAME> -n <NAMESPACE>
   ```

   Look at the `Status.Conditions` section for the Ready condition. The `Reason`
   and `Message` fields indicate why the forwarder is not ready. Common reasons
   include:

   - `ValidationFailure` — the ClusterLogForwarder spec has configuration errors

2. **Check input, output, pipeline, and filter conditions:**

   ```bash
   oc get clusterlogforwarder <NAME> -n <NAMESPACE> -o jsonpath='{.status}' | jq .
   ```

   Review `inputConditions`, `outputConditions`, `pipelineConditions`, and
   `filterConditions` for specific validation errors on individual components.

3. **Check if the collector pods are running:**

   ```bash
   oc get pods -n <NAMESPACE> -l app.kubernetes.io/component=collector
   ```

   If no pods exist, the collector was not deployed (likely due to a validation
   failure). If pods exist but are not running, check their status and logs.

4. **Review collector pod logs:**

   ```bash
   oc logs <POD_NAME> -n <NAMESPACE>
   ```

5. **Check the operator logs for reconciliation errors:**

   ```bash
   oc logs -n openshift-logging -l name=cluster-logging-operator
   ```

   Look for errors related to the ClusterLogForwarder reconciliation.

## Mitigation

### Fix Validation Errors

If the Ready condition reason is `ValidationFailure`:

1. Review the ClusterLogForwarder spec for configuration errors:

   ```bash
   oc get clusterlogforwarder <NAME> -n <NAMESPACE> -o yaml
   ```

2. Check the individual component conditions (inputs, outputs, pipelines, filters)
   to identify which component failed validation.

3. Correct the configuration errors in the ClusterLogForwarder spec. The operator
   will automatically re-reconcile and deploy the collector once the spec is valid.

### Restart the Operator

If the issue persists and no clear configuration error is found:

```bash
oc delete pod -n openshift-logging -l name=cluster-logging-operator
```

The operator will restart and re-reconcile all ClusterLogForwarder resources.

## Notes

- The `log_forwarder_ready` metric exposes the Ready condition with labels
  `resource_namespace`, `resource_name`, and `status` (True/False/Unknown).
- The alert uses `status="False"` to detect the not ready state. A status of
  `Unknown` (e.g., during operator restart) does not trigger this alert.
