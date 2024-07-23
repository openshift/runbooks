# TelemeterClientFailures

## Meaning

The alert `TelemeterClientFailures` is triggered when the Telemeter client fails
to send Telemetry data at a certain rate over a period of time
to Red Hat.

The `telemeter-client` pod running in the `openshift-monitoring`
namespace collects [selected platform metrics](https://docs.openshift.com/container-platform/latest/support/remote_health_monitoring/showing-data-collected-by-remote-health-monitoring.html#showing-data-collected-from-the-cluster_showing-data-collected-by-remote-health-monitoring)
from the `prometheus-k8s` service at
regular intervals using the `/federate` endpoint and ships them
to Red Hat using a custom secured protocol.

## Impact

When the alert fires, Red Hat support and engineering teams don't have a complete
view of the cluster anymore. It may hinder the ability for Red Hat to
proactively detect issues in the cluster.

## Diagnosis

* Review the logs for the pod `telemeter-client`
in the `openshift-monitoring` namespace.

You can review the logs for the `telemeter-client` pod in the
`openshift-monitoring` namespace by running the following command:

```console
oc logs -n openshift-monitoring deployment.apps/telemeter-client -c telemeter-client -f
```

* Open the Observe > Metrics page in the OCP admin console and execute the following
  PromQL expressions to identify where the issue happens.

  * OCP 4.17 and above

    ```console
    sum by(client, status_code) (rate(metricsclient_http_requests_total{status_code!~"200"}[15m])) > 0
    ```

    * The value of the `client` label is `federate_from` when the Telemeter client
      failed to retrieve metrics from Prometheus.
    * The value of the `client` label is `federate_to` when the Telemeter client
      failed to send metrics to Red Hat.

  * OCP 4.16 and below

    * The following query returns result when the Telemeter client failed to retrieve
      metrics from Prometheus.

      ```console
      sum by(client, status_code) (rate(metricsclient_request_retrieve{status_code!~"200"}[15m])) > 0
      ```

    * The following query returns result when the Telemeter client failed to send
      metrics to the Red Hat.

      ```console
      sum by(client, status_code) (rate(metricsclient_request_send{status_code!~"200"}[15m])) > 0
      ```

## Mitigation

The resolution of the issue depends on the origin of the failure.

* The telemeter client fails to retrieve metrics from Prometheus.
  * You need to check the availability of the Prometheus pods in the `openshift-monitoring`
    namespace. If the pods are running, check the logs of the `prometheus` container.

* The telemeter client fails to send metrics to the server.

  * If you use a firewall, make sure that it configured as specified in the
    [OCP documentation](https://docs.openshift.com/container-platform/latest/installing/install_config/configuring-firewall.html).

  * Check the `status_code` label values returned by the PromQL query executed
    at the previous step.
    * 401 and 403 codes indicate a misconfiguration of the client. A typical reason
      is the global [cluster pull secret](https://docs.openshift.com/container-platform/latest/openshift_images/managing_images/using-image-pull-secrets.html#images-update-global-pull-secret_using-image-pull-secrets).
      * Make sure your global cluster pull secret is up to date
    * Status codes between 500 and 599 indicate a problem from the Telemeter
      server side.
      * If you use HTTP proxies and/or firewalls, check their logs.
      * If the error is due to an outage on the Red Hat side and the alert
        doesn't resolve within an error, you can contact the Red Hat support.
