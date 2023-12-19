# PrometheusOperatorRejectedResources

## Meaning

The `PrometheusOperatorRejectedResources` alert triggers when Prometheus Operator
rejects invalid `AlertmanagerConfig`, `PodMonitor`, `ServiceMonitor`, or `PrometheusRule`
objects.

## Impact

The custom resources that trigger the alert are ignored by Prometheus Operator.
As a consequence, they will not be part of the final configuration of the
Prometheus, Alertmanager, or Thanos Ruler components managed by Prometheus Operator.

## Diagnosis

Identify the custom resource type and the namespace from the `resource` and `namespace`
labels of the `PrometheusOperatorRejectedResources` alert. You can find this information
by using the OpenShift web console or the command line interface (CLI).

### Using the Openshift web console

1. Browse to **Observe** -> **Alerting**.
2. Search for the `PrometheusOperatorRejectedResources` alert.
3. Click the alert to view its details.
4. Scroll down and view the **Labels** field. The resource type is indicated by
the `resource` label.

The `namespace` label can be either `openshift-monitoring` or `openshift-user-workload-monitoring`,
which dictates your course of action:
* When the value is `openshift-monitoring`, this is an issue with the platform
monitoring stack. Please submit a request to Customer Support.
* When the value is `openshift-user-workload-monitoring`, this is an issue with
a user-defined monitoring resource.

To identify the rejected monitoring object, check the logs of the Prometheus Operator
deployment in the `openshift-user-workload-monitoring` namespace. The namespace
and the name of the rejected resource will appear in the log entry after the
error message.
```bash
oc logs deployment/prometheus-operator -c prometheus-operator -n openshift-user-workload-monitoring
```

The following is a sample error message about a rejected custom resource:
```log
level=warn ts=2023-07-03T20:37:20.740723141Z caller=operator.go:1917 component=prometheusoperator msg="skipping servicemonitor" error="it accesses file system via bearer token file which Prometheus specification prohibits" servicemonitor=quarkus-demo/otel-collector namespace=openshift-user-workload-monitoring prometheus=user-workload
```

### Diagnose using the CLI

* Get the Alertmanager URL and view a list of alerts that have fired:
```bash
### Gets alertmanager URL
$ ALERT_MANAGER=$(oc get route alertmanager-main -n openshift-monitoring -o jsonpath='{@.spec.host}')

### Gets fired alerts
$ curl -sk -H "Authorization: Bearer $(oc create token prometheus-k8s -n openshift-monitoring)"  https://$ALERT_MANAGER/api/v2/alerts?filter=alertname=PrometheusOperatorRejectedResources
```

## Mitigation

Fix the misconfigured custom resource and reapply it to the cluster.

The root causes for this alert can vary depending on the specific configuration and
deployment of Prometheus Operator.
Possible causes include the following:

- **AlertManagerConfig**
  - Invalid receiver settings--for example, a missing URL in a Slack action.
  - Invalid route settings.
  - Settings that request a feature that is unavailable in the current version.
  - Unsupported match rules in inhibition rules.
- **ServiceMonitor and PodMonitor**
  - An invalid relabeling configuration--for example, a malformed regular expression.
  - An invalid TLS configuration.
  - A scrape interval configured to be longer than the scrape timeout.
  - Information missing in authentication settings.
  - Violation of file system access rules, which can occur when a `ServiceMonitor`
    or `PodMonitor` object references a file to use as a bearer token or references
    a TLS file. These configurations are not allowed in user-defined monitoring.
    Instead, you must create a secret that contains the credential data in the
    same namespace as the `ServiceMonitor` or `PodMonitor` object and use a secret
    key reference in the `ServiceMonitor` or `PodMonitor` configuration.
- **PrometheusRules**
  - A `PrometheusRules` object that contains an invalid expression.

The admission webhook should be able to catch the error of an invalid expression
in a `PrometheusRules` object. If that error shows up in the operator logs,
the admission webhook might be offline. Please check the deployment `prometheus-operator-admission-webhook`.

## Additional resources
- ["PrometheusOperatorRejectedResources" alert firing continuously in a Red Hat OpenShift Service in RHOCP 4](https://access.redhat.com/solutions/6992399)