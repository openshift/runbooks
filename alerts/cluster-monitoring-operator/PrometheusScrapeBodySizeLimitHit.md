# PrometheusScrapeBodySizeLimitHit

## Meaning

The `PrometheusScrapeBodySizeLimitHit` alert triggers when at least one
Prometheus scrape target replies with a response body larger than the
value configured in the `enforcedBodySizeLimit` field in the
`cluster-monitoring-config` config map in the `openshift-monitoring` namespace.

By default, no limit exists on the body size of scraped targets. When a value
is defined for `enforcedBodySizeLimit`, this limit prevents Prometheus from
consuming large amounts of memory if scraped targets return a response of a
size that exceeds the defined limit.

## Impact

Metrics coming from targets responding with a body size that exceeds the
configured size limit are not ingested by Prometheus. The targets are
considered to be down, and they will have their `up` metric set to `0`, which
might also trigger the `TargetDown` alert.

## Diagnosis

You can view the value set for the body size limit in the
`cluster-monitoring-config` config map in the `openshift-monitoring` namespace.
View this value by entering the following command:

```bash
oc get cm -n openshift-monitoring cluster-monitoring-config -o yaml | grep enforcedBodySizeLimit
```

To discover the targets that are exceeding the configured body size limit, open
the OpenShift web console, go to the **Observe** --> **Targets** page, and check
to see which targets are down.

To get more information than is available in the web console, such as details
about discovered labels and the scrape pool, query the Prometheus API endpoint
`/api/v1/targets`. Querying this endpoint will return useful debugging
information for every target, as shown in the following example:

```json
{
  "status": "success",
  "data": {
    "activeTargets": [
      {
        "discoveredLabels": {
          "__address__": "10.128.0.6:8443",
          "__meta_kubernetes_endpoint_address_target_kind": "Pod",
          "__meta_kubernetes_endpoint_address_target_name": "openshift-apiserver-operator-7475b9d64-mdrlc",
          "__meta_kubernetes_endpoint_node_name": "ci-ln-6tfxd7t-72292-d7lf2-master-2",
          ...
          "job": "serviceMonitor/openshift-apiserver-operator/openshift-apiserver-operator/0"
        },
        "labels": {
          "container": "openshift-apiserver-operator",
          "endpoint": "https",
          "instance": "10.128.0.6:8443",
          "job": "metrics",
          "namespace": "openshift-apiserver-operator",
          "pod": "openshift-apiserver-operator-7475b9d64-mdrlc",
          "service": "metrics"
        },
        "scrapePool": "serviceMonitor/openshift-apiserver-operator/openshift-apiserver-operator/0",
        "scrapeUrl": "https://10.128.0.6:8443/metrics",
        "globalUrl": "https://10.128.0.6:8443/metrics",
        "lastError": "",
        "lastScrape": "2022-07-05T13:59:42.924932804Z",
        "lastScrapeDuration": 0.017444282,
        "health": "up",
        "scrapeInterval": "30s",
        "scrapeTimeout": "10s"
      },
      ...
    ],
    "droppedTargets": [
      ...
    ],
  }
}
```

In the `data.activeTargets` field, search for targets in which the value of the
`health` field is not `up`, and check the `lastError` field for confirmation:

1. Get the name of the secret containing the token of the `prometheus-k8s`
   service account by running the following command to check for the name
   `prometheus-k8s-token-[a-z]+` in the line `Tokens`. The following example
   uses the secret name `prometheus-k8s-token-nwtrf`:

    ```bash
    $ oc describe sa prometheus-k8s -n openshift-monitoring
    Name:                prometheus-k8s
    Namespace:           openshift-monitoring
    Labels:              app.kubernetes.io/component=prometheus
                        app.kubernetes.io/instance=k8s
                        app.kubernetes.io/name=prometheus
                        app.kubernetes.io/part-of=openshift-monitoring
                        app.kubernetes.io/version=2.35.0
    Annotations:         serviceaccounts.openshift.io/oauth-redirectreference.prometheus-k8s:
                          {"kind":"OAuthRedirectReference","apiVersion":"v1","reference":{"kind":"Route","name":"prometheus-k8s"}}
    Image pull secrets:  prometheus-k8s-dockercfg-6wl6b
    Mountable secrets:   prometheus-k8s-dockercfg-6wl6b
    Tokens:              prometheus-k8s-token-nwtrf
    Events:              <none>
    ```
2. If no token exists, create a token by entering the following command and
   skip the next step:

    ```bash
    token=$(oc sa new-token prometheus-k8s -n openshift-monitoring)
    ```

3. If a token exists, decode the token from the secret:

    ```bash
    token=$(oc get secret $secret_name_here -n openshift-monitoring -o jsonpath={.data.token} | base64 -d)
    ```

4. Get the route URL of the Prometheus API endpoint:

    ```bash
    host=$(oc get route -n openshift-monitoring prometheus-k8s -o jsonpath={.spec.host})
    ```

5. List all of the scraped targets with a health status different from `up`:

    ```bash
    curl -H "Authorization: Bearer $token" -k https://${host}/api/v1/targets | jq '.data.activeTargets[]|select(.health!="up")'
    ```

6. Review the response body size of the failing target. Enter the following
   command to simulate a scrape of the target's `scrapeUrl` and view the
   response body size:

    ```bash
    oc exec -it prometheus-k8s-0 -n openshift-monitoring  -- curl -k --key /etc/prometheus/secrets/metrics-client-certs/tls.key --cert /etc/prometheus/secrets/metrics-client-certs/tls.crt --cacert /etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt $scrape_url | wc --bytes
    ```

## Mitigation

Your analysis of the issue might reveal that the alert was triggered by one of
two causes:

* The value set for `enforcedBodySizeLimit` is too small.
* A bug exists in the target causing it to report too many metrics.

### Increasing the Body Size Limit

You can increase the body size limit by editing the `cluster-monitoring-config`
config map in the `openshift-monitoring` namespace.

The `prometheusK8s.enforcedBodySizeLimit` field defines this limit. Values for
this field use the [Prometheus size format](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#size).

The following example sets the body size limit to 20MB:

 ```yaml
apiVersion: v1
data:
  config.yaml: |-
    prometheusK8s:
      enforcedBodySizeLimit: 20MB
kind: ConfigMap
metadata:
  name: cluster-monitoring-config
  namespace: openshift-monitoring
 ```

### Bug in the Scraped Target

If you think that the response of the scraped target is too large, you can
contact Red Hat Customer Experience & Engagement.
