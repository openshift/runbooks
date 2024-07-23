# TelemeterClientFailures

## Meaning

The alert `TelemeterClientFailures` is triggered when the Telemeter client tries
and fails to send Telemetry data at a certain rate over a period of time
to Red Hat.

The `telemeter-client-xxxxxxxxx-xxxxx` pod running in the `openshift-monitoring`
namespace collects [selected platform metrics](https://docs.openshift.com/container-platform/latest/support/remote_health_monitoring/showing-data-collected-by-remote-health-monitoring.html#showing-data-collected-from-the-cluster_showing-data-collected-by-remote-health-monitoring)
from the `prometheus-k8s` service at
regular intervals using the `/federate` endpoint and ships them periodically
to Red Hat using a custom secured protocol.

## Impact

When the alert fires, Red Hat support and engineering teams don't have a complete
view of the cluster anymore. It may hinder the ability for Red Hat to
proactively detect issues in the cluster.

## Diagnosis

Review the logs for the pod `telemeter-client-xxxxxxxxx-xxxxx`
in the `openshift-monitoring` namespace

You can review the logs for the `telemeter-client-xxxxxxxxx-xxxxx` pod in the
`openshift-monitoring` namespace by running the following command:

```console
$ oc logs -n openshift-monitoring deployment.apps/telemeter-client -c telemeter-client | less
```

## Mitigation

The resolution depends on the particular issue reported in the logs.

* Telemeter client fails to pull metrics from Prometheus.
  * You'd need to check the availability of Prometheus.

* Telemeter client fails to push to the server.

  * It can be due to Red Hat Telemeter server outage.

    ```console
    2024-05-05T09:33:28.347041825+08:00 level=error caller=forwarder.go:276 ts=2024-05-05T01:33:28.346019008Z
    component=forwarder/worker msg="unable to forward results" err="gateway server reported unexpected error code: 503:
    <html>\r\n  <head>\r\n    <meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\r\n\r\n    <style
    type=\"text/css\">\r\n      body {\r\n        font-family: \"Helvetica Neue\", Helvetica, Arial, sans-serif;\r\n       
    line-height: 1.66666667;\r\n        font-size: 16px;\r\n        color: #333;\r\n        background-color: #fff;\r\n       
    margin: 2em 1em;\r\n      }\r\n      h1 {\r\n        font-size: 28px;\r\n        font-weight: 400;\r\n      }\r\n      p
    {\r\n        margin: 0 0 10px;\r\n      }\r\n      .alert.alert-info {\r\n        background-color: #F0F0F0;\r\n        margin-top: 30px;\r\n        padding: 30px;\r\n      }\r\n      .alert p {\r\n        padding-left: 35px;\r\n      
    \r\n      ul {\r\n        padding-left: 51px;\r\n        position: relative;\r\n      }\r\n      li {\r\n       
    font-size: 14px;\r\n        margin-bottom: 1em;\r\n      }\r\n      p.info {\r\n        position: relative;\r\n       
    font-size: 20px;\r\n      }\r\n      p.info:before, p.info:after {\r\n        content: \"\";\r\n        left: 0
    \r\n        position: absolute;\r\n        top: 0;\r\n      }\r\n"

    2024-05-05T09:34:42.700742128+08:00 level=error caller=forwarder.go:276 ts=2024-05-05T01:34:42.700567307Z
    component=forwarder/worker msg="unable to forward results" err="gateway server reported unexpected error code: 502:
    <html><body><h1>502 Bad Gateway</h1>\nThe server returned an invalid or incomplete response.\n</body></html>\n"
    ```

  * The brief outage on telemetry service can cause some clusters not being able
    to report telemetry to Red Hat. Except for this, the outage should have no impact.
    TelemeterClientFailures will resolve as soon as telemetry outage is resolved.

* It can be an issue with the [cluster pull secret](https://docs.openshift.com/container-platform/latest/openshift_images/managing_images/using-image-pull-secrets.html#images-update-global-pull-secret_using-image-pull-secrets)
  * Make sure your cluster pull secret is up to date
