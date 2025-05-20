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

1. Examine `CollectorHigh403ForbiddenResponseRate` alert details in the OpenShift
console:
    - Note the key labels:
        - `app_kubernetes_io_instance`: The name of the collector.
        - `namespace`: Namespace of the collector.
        - `component_id`: The ID of the specific sink in `Vector` that's failing
        (i.e `output_lokistack_otlp_application`, `output_my_splunk`).
2. Review the `ClusterLogForwarder` (CLF) configuration:
    1. Identify the output definition corresponding to the `component_id` obtained
    from the alert details.
        - See naming schemas for [non-Lokistack](#non-lokistack-outputs) and
        [Lokistack](#lokistack-outputs) outputs.
    2. Examine the secret or authentication fields defined for the affected output.
3. Validate authentication credentials in the identified secrets:
    - Identify the Kubernetes `Secret` object referenced by the affected CLF output.
    - Verify the following:
        1. The `Secret` object exists in the specified namespace.
        2. The expected keys (e.g., token, password, etc.) are
        present and contain valid, non-empty values.
4. Verify log store authorization with the provided credentials.
    - Confirm that the credentials provided in the identified Secret possess the
    necessary permissions (i.e write) to the target log store.

### Diagnosis for Red Hat Managed Lokistack

#### Verify `ClusterRole`'s Existence
Confirm the presence of the `ClusterRole` named `logging-collector-logs-writer`
within the cluster. This `ClusterRole` is provisioned during the deployment
of the `Cluster Logging Operator`.
The `ClusterRole` is not managed by any operator.

##### Query for the `ClusterRole`
```bash
$ oc get clusterrole logging-collector-logs-writer
```
#### Validate `ClusterRoleBinding`
Ensure the CLF's associated `ServiceAccount` is correctly
bound to the `logging-collector-logs-writer` `ClusterRole`.

##### Query for the `ServiceAccount` associated with the CLF
```bash
$ oc get clusterlogforwarder <CR-NAME> -o jsonpath='{.spec.serviceAccount}' -n <NAMESPACE>
```

##### Query for the `ClusterRoles` that the `ServiceAccount` is bound to
```bash
$ oc get clusterrolebinding -o json | jq -r 'try .items[]? | select(.subjects[].name=="<SERVICEACCOUNT>").roleRef.name'
```

## Mitigation

### Credential Update

If the issue stems from incorrect or expired authentication
credentials, update the associated OpenShift `Secret` with the correct and valid
credentials.

### CLF Configuration
If a typo or misconfiguration is identified within the
CLF custom resource, directly edit the resource to fix the error.

### External Log Store Authorization
Work with the external log store's
administrator to ensure the provided credentials have the proper permissions
(i.e write).
This typically involves adjusting roles or access policies.

### Red Hat Managed Lokistack
For Red Hat managed Lokistack deployments, ensure
the `logging-collector-logs-writer` `ClusterRole` exists and that the `ServiceAccount`
is bound to the `ClusterRole`.

#### Check existence of `logging-collector-logs-writer` `ClusterRole`
```bash
$ oc get clusterrole logging-collector-logs-writer
```

If `ClusterRole` is absent there are two options:

1. Redeploy the `Cluster Logging Operator`.

2. Manually create the `logging-collector-logs-writer` `ClusterRole` using the
[upstream CLO ClusterRole manifest](https://github.com/openshift/cluster-logging-operator/blob/master/bundle/manifests/logging-collector-logs-writer_rbac.authorization.k8s.io_v1_clusterrole.yaml).

```bash
$ oc create -f bundle/manifests/logging-collector-logs-writer_rbac.authorization.k8s.io_v1_clusterrole.yaml
```

#### Binding the `ServiceAccount` to the `logging-collecor-logs-writer` `ClusterRole`

1. Query for the `ServiceAccount` associated with the CLF:
```bash
$ oc get clusterlogforwarder <CR-NAME> -o jsonpath='{.spec.serviceAccount}' -n <NAMESPACE>
```

2. Inspect which `ClusterRoles` the `ServiceAccount` is bound to:
```bash
$ oc get clusterrolebinding -o json | jq -r 'try .items[]? | select(.subjects[].name=="<SERVICEACCOUNT>").roleRef.name'
```

If the `ServiceAccount` is not bound to the `logging-collecor-logs-writer` `ClusterRole`:

3. Bind `ServiceAccount` to `ClusterRole`
```bash
$ oc adm policy add-cluster-role-to-user logging-collector-logs-writer -z <SA-NAME> -n <SA-NAMESPACE>
```

# Notes

## Naming schemas

### Non-Lokistack Outputs
The naming schema is `output_<OUTPUT_NAME>`,
where the `<OUTPUT_NAME>` corresponds to the CLF's named output, with all
punctuations replaced by underscores
(e.g., `output_my_splunk` will correspond to a CLF output named `my-splunk`).
### Lokistack Outputs
The naming schema is `output_<OUTPUT_NAME>_<INPUT_TYPE>`,
where the `<OUTPUT_NAME>` corresponds to the CLF's named Lokistack output,
with all punctuations replaced by underscores and `<INPUT_TYPE>` is the tenant.
(e.g `output_my_lokistack_application` will correspond to a CLF output named
`my-lokistack` receiving `application` logs.)