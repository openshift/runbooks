# CoreDNSErrorsHigh

## Meaning

CoreDNS is returning SERVFAIL more than 1% of total requests.

## Impact

Warning: Some requests are not resolved.

## Diagnosis

1. Check DNS operator:

    ```shell
    oc get co/dns
    ```

2. Check CoreDNS daemonsets if they are running on all nodes:

    ```shell
    oc get ds -n openshift-dns
    ```

3. Check CoreDNS pods if they can reach the upstream nameservers:

    ```shell
    # CoreDNS pods name start with dns-default as prefix
    oc exec pod/$COREDNS_POD_NAME -c dns -n openshift-dns -- \
        dig -q $DOMAIN_NAME @$IP_OF_UPSTREAM_NAMESERVER
    ```

4. Check the upstream nameservers if they are returning SERVFAIL:

    ```shell
    # Enable Debug logLevel in CoreDNS
    oc patch dnses.operator.openshift.io/default \
        --patch '{"spec":{"logLevel":"Debug"}}' --type=merge
    
    # Follow the log streams in proportion to the number of CoreDNS pods
    oc logs -c dns -l dns.operator.openshift.io/daemonset-dns=default \
        -n openshift-dns --follow --max-log-requests $NUMBER_OF_COREDNS_PODS --timestamps
    ```

## Mitigation

- If there is a connectivity issue between the CoreDNS pods and the workload,
review the [Controlling DNS pod placement](https://docs.openshift.com/container-platform/4.17/networking/networking_operators/dns-operator.html#nw-controlling-dns-pod-placement_dns-operator):

    ```shell
    oc edit dns.operator/default
    ```

    ```yaml
    spec:
      nodePlacement:
        ...output omitted...
    ```

- If there is a connectivity issue between the CoreDNS pods and the upstream
nameservers, review the undercloud connectivity.

- If the upstream nameservers is not healthy to respond to the queries by
the CoreDNS pods, apply a silence to the alert, until these servers are troubleshooted.
