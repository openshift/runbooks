# KubeAggregatedAPIErrors

## Meaning

The `KubeAggregatedAPIErrors` alert is triggered when multiple calls to the
aggregated OpenShift API fail over a certain period of time.

## Impact

Aggregated API errors can result in the unavailability of some OpenShift
services.

## Diagnosis

The alert message contains information about the affected API and the scope of
the impact, as shown in the following sample:

```text
 - alertname = KubeAggregatedAPIErrors
 - name = v1.packages.operators.coreos.com
 - namespace = default
...
 - message = Kubernetes aggregated API v1.packages.operators.coreos.com/default has reported errors. It has appeared unavailable 5 times averaged over the past 10m.
```

## Mitigation

Troubleshoot and fix the issue or issues causing the aggregated API errors by
checking the availability status for each API and by verifying the
authentication certificates for the aggregated API.

### Check the availability status for each API

At least four aggregated APIs exist in an OpenShift cluster:

* the API for the `openshift-apiserver` namespace
* the API for the `prometheus-adapter` in the namespace `openshift-monitoring`
* the API for the the `package-server` service in the
`openshift-operator-lifecycle-manager` namespace
* the API for the `openshift-oauth-apiserver` namespace

1. Check the availability of all APIs. To get a list of `APIServices` and their
backing aggregated APIs, use the following command:

    ```console
    $ oc get apiservice
    ```

    The `SERVICE` column in the returned data shows the aggregated API name.
    Normally, the availability status for every listed API will be shown as
    `True`. If the status is `False`, it means that requests for that API
    service, API server pods, or resources belonging to that `apiGroup` have
    failed many times in the past few minutes.

2. Fetch the pods that serve the unavailable API. For example, for
`openshift-apiserver/api` use the following command:

    ```console
    $ oc get pods -n openshift-apiserver
    ```

    If the status is not shown as `Running`, review the logs for more details.
    Because these pods are controlled by a deployment, they can be restarted
    when they do not respond to requests.

### Verify the authentication certificates for the aggregated API

1. Verify that the certificates have not expired and are still valid:

    ```console
    $ oc get configmaps -n kube-system extension-apiserver-authentication
    ```

    If required, you can save these certificates to a file and use the following
    command to check the expiration dates for each certificate file:

    ```console
    $ openssl x509 -noout -enddate -in {myfile_with_certs.crt}
    ```

    The aggregated APIs use these certificates to validate requests. If
    they are expired, see [the OpenShift documentation][cert] for information
    about how to add a new certificate.

[cert]: https://docs.openshift.com/container-platform/latest/security/certificates/api-server.html
