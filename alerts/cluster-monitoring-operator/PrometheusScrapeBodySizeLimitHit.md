# PrometheusScrapeBodySizeLimitHit

## Meaning

The alert `PrometheusScrapeBodySizeLimitHit` is triggered when at least one of
Prometheus' scrape targets replies with a response body larger than the
configured `body_size_limit`.
By default, there is no limit on the body size of the scraped targets.
When set, this limit prevents Prometheus from consuming excessive amounts of memory
when scraped targets return a response that is deemed too large.

## Impact

Metrics coming from targets responding with a body exceeding the configured
size limit aren't ingested by Prometheus. The targets will be considered as being
down and they will have their `up` metric set to 0 which may also trigger the
`TargetDown` alert.

## Diagnosis

We can find the value of the configured body size limit in the configmap `cluster-monitoring-config`
in the namespace `openshift-monitoring`.

To check the value in configmap, use this command:
```bash
oc get cm -n openshift-monitoring cluster-monitoring-config -o yaml | grep enforcedBodySizeLimit
```

To find out which targets are exceeding the body size limit, we can refer
to the Openshift Web UI, go to the Observe > Targets page, and check which
targets are down.

If we need more information than the UI provides, such as discovered labels
and scrape pool, we can query the Prometheus API endpoint `/api/v1/targets`.

For every target, we have some useful information for debugging, as
shown in the example below:

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

In the field `data.activeTargets` search for those that have `health` field
different than `up` and check for the `lastError` field for confirmation:

1. Get the name of the secret containing the token of service account `prometheus-k8s`.
We use the following command and check for the name  `prometheus-k8s-token-[a-z]+`
in the line `Tokens`. In the example below, the secret name is `prometheus-k8s-token-nwtrf`.

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
2. a. If there is no token available, we can create a token using this command
below and skip step 2.b.
```bash
token=$(oc sa new-token prometheus-k8s -n openshift-monitoring)
```


2. b. If the token exists, decode the token from the secret.

```bash
token=$(oc get secret $secret_name_here -n openshift-monitoring -o jsonpath={.data.token} | base64 -d)
```

3. Get the route URL of the Prometheus API endpoint.

```bash
host=$(oc get route -n openshift-monitoring prometheus-k8s -o jsonpath={.spec.host})
```

4. List all the scraped targets with health status different than "up".

```bash
curl -H "Authorization: Bearer $token" -k https://${host}/api/v1/targets | jq '.data.activeTargets[]|select(.health!="up")'
```

5. Check the failling target's response body size. We simulate a scrape of its
`scrapeUrl` using the command below and check for the body size.

```bash
oc exec -it prometheus-k8s-0 -n openshift-monitoring  -- curl -k --key /etc/prometheus/secrets/metrics-client-certs/tls.key --cert /etc/prometheus/secrets/metrics-client-certs/tls.crt --cacert /etc/prometheus/configmaps/serving-certs-ca-bundle/service-ca.crt $scrape_url | wc --bytes
```

## Mitigation

According to the result of the investigation, the cause can be either the limit is
too small or there is a bug in the target causing it to report too many metrics.

### Increasing the Body Size Limit

We can increase the body size limit by editing the configmap
`cluster-monitoring-config` in the namespace `openshift-monitoring`. The field
`prometheusK8s.enforcedBodySizeLimit` defines this limit, accepting [the same size
 format as Prometheus uses](https://prometheus.io/docs/prometheus/latest/configuration/configuration/#size).

 The example below sets the body size limit to 20MB.

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

If you believe that the response of the scraped target is too large,
you can contact Red Hat Customer Experience & Engagement.
