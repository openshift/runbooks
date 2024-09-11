# PrometheusOperatorRejectedResources

## Meaning

The `PrometheusOperatorRejectedResources` alert triggers when the Prometheus
Operator detects and rejects invalid `AlertmanagerConfig`, `PodMonitor`,
`ServiceMonitor`, or `PrometheusRule` objects.

## Impact

The custom resources that trigger the alert are ignored by Prometheus Operator.
As a consequence, they will not be part of the final configuration of the
Prometheus, Alertmanager, or Thanos Ruler components managed by the Prometheus
Operator.

## Diagnosis

### Identify the custom resource type

The first step is to identify the custom resource type and the namespace from
the `resource` and `namespace` labels of the
`PrometheusOperatorRejectedResources` alert. You can find this information by
using the OpenShift web console or the command line interface (CLI).

#### Using the Openshift web console

1. Browse to **Observe** -> **Alerting**.
2. Search for the `PrometheusOperatorRejectedResources` alert.
3. Click the alert to view its details.
4. Scroll down and view the **Labels** field. The resource type is indicated by
the `resource` label.

#### Using the CLI

Get the Alertmanager URL and view a list of alerts that have fired:

```bash
### Retrieve the Alertmanager URL
$ ALERTMANAGER=$(oc get route alertmanager-main -n openshift-monitoring -o jsonpath='{@.spec.host}')

### Get active alerts
$ curl -sk -H "Authorization: Bearer $(oc create token prometheus-k8s -n openshift-monitoring)" \
    https://$ALERTMANAGER/api/v2/alerts?filter=alertname=PrometheusOperatorRejectedResources
```

The `namespace` label can be either `openshift-monitoring` or `openshift-user-workload-monitoring`,
which dictates your course of action:
* When the value is `openshift-monitoring`, this is an issue with the platform
monitoring stack. Please submit a request to the Customer Support.
* When the value is `openshift-user-workload-monitoring`, this is an issue with
a user-defined monitoring resource.


### Identify the resource(s) and reason

To identity which monitoring objects have been rejected and why, use one of the
2 following methods depending on the OCP version.

#### OCP 4.16 and later

The Prometheus operator emits events about invalid resources.

##### Using the Openshift web console

1. Browse to **Home** -> **Events**.
2. Select "All projects" in the Project drop-down list.
3. Select the `resource` label in the Resources drop-down list (for instance,
`ServiceMonitor`).

##### Using the CLI

Check the Kubernetes events related to the `resource` label using the following
command (example for `ServiceMonitor` resources):

```bash
oc get events --field-selector involvedObject.kind=ServiceMonitor --all-namespaces
```

The following is a sample event about a rejected custom resource:

```log
NAMESPACE   LAST SEEN   TYPE      REASON                 OBJECT                       MESSAGE
default     106s        Warning   InvalidConfiguration   servicemonitor/example-app   ServiceMonitor example-app was rejected due to invalid configuration: it accesses file system via bearer token file which Prometheus specification prohibits
```


#### Before OCP 4.16

Check the logs of the Prometheus Operator deployment in the
`openshift-user-workload-monitoring` namespace. The namespace and the name of
the rejected resource will appear in the log entry after the error message.

```bash
oc logs deployment/prometheus-operator -c prometheus-operator -n openshift-user-workload-monitoring
```

The following is a sample error message about a rejected custom resource:

```log
level=warn ts=2023-07-03T20:37:20.740723141Z caller=operator.go:1917 component=prometheusoperator msg="skipping servicemonitor" error="it accesses file system via bearer token file which Prometheus specification prohibits" servicemonitor=quarkus-demo/otel-collector namespace=openshift-user-workload-monitoring prometheus=user-workload
```

## Mitigation

The mitigation depends on which resources are being rejected and why.

### ServiceMonitor and PodMonitor

- Invalid relabeling configuration (for example, a malformed regular expression).
  - Fix the relabeling configuration syntax.
- Invalid TLS configuration.
  - Fix the TLS configuration.
- A scrape interval less than the scrape timeout.
  - Change the scrape timeout or the scrape interval value.
- Invalid secret or configmap key reference.
  - Verify that the secret/configmap object exists and that they key is present
    in the secret/configmap.
- Violation of file system access rules, which can occur when a `ServiceMonitor`
  or `PodMonitor` object references a file to use as a bearer token or references
  a TLS file. These configurations are not allowed in user-defined monitoring.
  - you must create a secret that contains the credential data in the
    same namespace as the `ServiceMonitor` or `PodMonitor` object and use a
    secret key reference in the `ServiceMonitor` or `PodMonitor`
    configuration.

When the alert is triggered by an resource managed by a 3rd-party operator, it
might not be possible to fix the root cause. The resolution will depend on the
status of the operator:

- The operator is a certified Red Hat operator.
  - If the operator is installed in the `openshift-operators` namespace, it
    should be removed and installed in a different namespace because
    `openshift-operators` might contain community operators which don't have
    the same level of support.
  - If the operator is deployed in another namespace than `openshift-operators`
    and its documentation requires adding the
    `openshift.io/cluster-monitoring: "true"` label to this namespace during
    the installation, ensure that the label exists.
  - Otherwise you can exclude the resource from user-defined monitoring by adding
    the `openshift.io/user-monitoring:"false"` label to the resource's namespace
    or the resource itself (the latter requires at least OCP 4.16).
- The operator is a community operator.
  - You can exclude the resource from user-defined monitoring by adding the
    `openshift.io/user-monitoring:"false"` label to the resource's namespace or
    the resource itself (the latter requires at least OCP 4.16).


### AlertmanagerConfig

- Invalid secret or configmap key reference.
  - Verify that the secret/configmap object exists and that they key is present
    in the secret/configmap.
- Invalid receiver or route settings (for example, a missing URL in a Slack action).
  - Fix the improper syntax.
- Configuration option which is not yet available in the Alertmanager version.
  - Update the resource to not use this option.
- Unsupported match rules in inhibition rules.
  - Fix the match rule syntax.

The admission webhook should be able to catch most of these errors. In this
case, the admission webhook might be offline. Please check the
`prometheus-operator-admission-webhook` deployment in the
`openshift-monitoring` namespace.


### PrometheusRule

The resource can be invalid because it contains an invalid expression which
needs to be fixed. The admission webhook should be able to catch the error of
an invalid expression in a `PrometheusRule` object. In this case, the admission
webhook might be offline. Please check the
`prometheus-operator-admission-webhook` deployment in the
`openshift-monitoring` namespace.


## Additional resources

- ["PrometheusOperatorRejectedResources" alert firing continuously in a Red Hat OpenShift Service in RHOCP 4](https://access.redhat.com/solutions/6992399)
