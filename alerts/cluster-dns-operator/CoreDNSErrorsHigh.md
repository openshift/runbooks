# CoreDNSErrorsHigh

## Meaning

CoreDNS is returning SERVFAIL more than 1% of total requests.

## Impact

Warning: Some requests are not resolved.

## Diagnosis

1. Check dns's operator:

    ```shell
    $ oc get co dns
    ```

2. Check dns's daemonsets if they are running on all nodes:

    ```shell
    $ oc get ds -n openshift-dns
    ```

3. Check dns' pods if they can reach upstream nameservers:

    ```shell
    $ oc exec -c dns dns-default-xxxxx -it -n openshift-dns -- dig www.example.com @$IP_OF_UPSTREAM_NAMESERVER
    ```

4. Check dns' upstream nameservers if they are returning SERVFAIL:

    ```shell
    $ oc logs -c dns -l dns.operator.openshift.io/daemonset-dns=default -n openshift-dns 
    ```

## Mitigation

If there is a connectivity issue between the coredns pods and the workload, review the [Controlling DNS pod placement](https://docs.openshift.com/container-platform/4.17/networking/networking_operators/dns-operator.html#nw-controlling-dns-pod-placement_dns-operator):

```shell
$ oc edit dns.operator/default
```

```yaml
spec:
  nodePlacement:
      ...output omitted...
```

If there is a connectivity issue between the coredns pods and the upstream nameserver, review the undercloud connectivity.

If the upstream nameservers is not healthy to respond to the queries by the coredns pods, apply a silence to the alert, until these servers are troubleshooted.
