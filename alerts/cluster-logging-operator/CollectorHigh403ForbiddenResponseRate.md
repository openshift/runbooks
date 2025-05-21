# CollectorHigh403ForbiddenResponseRate

## Meaning

The `CollectorHigh403ForbiddenResponseRate` alert indicates that the log collector
is experiencing a sustained high rate (at least 10% over the last 2 minutes, persisting
for 5 minutes) of **HTTP 403 Forbidden** responses when attempting to send log
data to its configured log store.

An **HTTP 403 Forbidden** response indicates that the server understood the request
but refuses to authorize it.

## Impact

Log collection and forwarding to the destination are impaired,
resulting in data loss at the configured log store.

## Diagnosis

1. Examine `CollectorHigh403ForbiddenResponseRate` alert details in OpenShift console:
    - Note the key labels:
        - `app_kubernetes_io_instance`: The name of the collector.
        - `namespace`: Namespace of the collector.
        - `component_id`: The ID of the specific output/sink that's failing. This
        corresponds to an output defined in your `ClusterLogForwarder` (CLF) custom
        resource.
2. Review the CLF Configuration:
    1. Identify the output definition corresponding to the `component_id` obtained
    from the alert details.
    2. Examine the secret or authentication fields defined for the affected output.
3. Validate authentication credentials in the identified secrets:
    - Identify the Kubernetes Secret object referenced by the affected CLF output.
    - Verify the following:
        1. The Secret object exists in the specified namespace.
        2. The expected keys (e.g., token, password, etc.) are
        present and contain valid, non-empty values.
4. Verify log store authorization with the provided credentials.
    - Confirm that the credentials provided in the identified Secret possess the
    necessary write permissions to the target log store.

### Diagnosis for Red Hat Managed Lokistack
1. Verify `clusterRole` (CR) Existence:
    - Confirm the presence of the CR named `logging-collector-logs-writer`
    within the cluster. This CR is provisioned during the installation
    of the `Cluster Logging Operator`.
2. Validate `clusterRoleBinding` (CRB):
    - Ensure that the `logging-collector-logs-writer` CR is correctly bound to
    the `ClusterLogForwarder`'s associated `ServiceAccount` via a CRB.

## Mitigation

1. **Credential Update**: If the issue stems from incorrect or expired authentication
credentials, update the associated OpenShift Secret with the correct and valid credentials.
2. **CLF Configuration**: If a typo or misconfiguration is identified within the
CLF custom resource, directly edit the CLF resource to fix the error.
3. **External Log Store Authorization**: Work with the external log store's
administrator to ensure the provided credentials have the proper write permissions.
This typically involves adjusting roles or access policies.
4. **Red Hat Managed Lokistack**: For Red Hat managed Lokistack deployments, ensure
the `logging-collector-logs-writer` CR exists and is correctly bound to the
ServiceAccount utilized by the CLF via a CRB. If the CR is absent, reinstall the
`Cluster Logging Operator`.
