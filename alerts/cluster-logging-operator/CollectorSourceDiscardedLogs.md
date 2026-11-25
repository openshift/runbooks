# CollectorSourceDiscardedLogs

## Meaning

The `CollectorSourceDiscardedLogs` alert indicates that the log collector is
discarding log lines at the source level. This typically occurs when log lines
exceed the configured `max_line_bytes` limit (default 3MB for audit logs).

The alert fires when either:
- `vector_component_discarded_events_total` increases for a source component, or
- `vector_component_errors_total` increases with error codes specific to oversized
  lines (`reading_line_from_file` or `reading_line_from_kubernetes_log`).

## Impact

- **Log data loss:** Oversized log lines are silently dropped and never forwarded
  to the configured log store.
- **Incomplete audit trail:** If audit logs are affected, security-relevant events
  may be missing from the log store.
- **Silent failure:** Without this alert, discarded logs produce no visible error
  in the forwarding pipeline — the logs simply disappear.

## Diagnosis

1. **Identify the affected source from alert labels:**
    - `namespace`: Namespace of the collector.
    - `app_kubernetes_io_instance`: Name of the `ClusterLogForwarder`.
    - `component_id`: The specific source input (e.g., `input_audit_host`,
      `input_infrastructure_container`).
    - `component_type`: Source type (`file` or `kubernetes_logs`).

2. **Check collector logs for discard warnings:**

   ```bash
   oc logs <POD_NAME> -n <NAMESPACE> | grep -i "max_line_bytes\|line too long\|discarding"
   ```

   Look for messages indicating which files contain oversized lines.

3. **Inspect the current discard metrics:**

   ```bash
   oc exec <POD_NAME> -n <NAMESPACE> -- \
     curl -s http://localhost:24231/metrics | grep -i discard
   ```

4. **Check the configured `max_line_bytes` value:**

   ```bash
   oc exec <POD_NAME> -n <NAMESPACE> -- \
     grep max_line_bytes /etc/vector/vector.toml
   ```

   The default for audit log sources is 3145728 (3MB).

5. **Identify the source of oversized log lines:**

   For audit logs, check if any API server or OAuth server is producing
   unusually large audit events:

   ```bash
   oc exec <POD_NAME> -n <NAMESPACE> -- \
     wc -L /var/log/kube-apiserver/audit.log
   ```

   This shows the length of the longest line in the file.

## Mitigation

### Investigate the Root Cause

Oversized log lines are often caused by:
- API requests or responses with large payloads being logged at a verbose audit
  level.
- Applications writing excessively long single-line log messages (e.g.,
  base64-encoded blobs, serialized objects).

### Adjust Audit Policy

If audit logs are affected, consider adjusting the audit policy to reduce the
verbosity for specific API groups or resources that produce large events:
- Use `Metadata` level instead of `Request` or `RequestResponse` for
  high-volume or large-payload resources.

### Monitor After Changes

After making changes, verify the alert resolves by checking that the discard
metrics stop increasing:

```bash
oc exec <POD_NAME> -n <NAMESPACE> -- \
  curl -s http://localhost:24231/metrics | grep vector_component_discarded_events_total
```

## Notes

- The `max_line_bytes` setting is configured per source in the collector's
  Vector configuration. It is not directly user-configurable through the
  `ClusterLogForwarder` API.
- The default value of 3MB (3145728 bytes) is set for audit log sources to
  accommodate large Kubernetes audit events.
- Container log sources (`kubernetes_logs`) use a separate
  `max_merged_line_bytes` setting for handling partial log line merging.
