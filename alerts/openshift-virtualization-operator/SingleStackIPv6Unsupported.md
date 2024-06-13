# SingleStackIPv6Unsupported

## Meaning

This alert fires when user tries to install OpenShift Virtualization on a single
stack IPv6 cluster.

OpenShift Virtualization is not yet supported on an OpenShift cluster configured
with single stack IPv6. It's progress is being tracked on [this
issue](https://issues.redhat.com/browse/CNV-28924).

## Impact

OpenShift Virtualization Operator can't be installed on a single stack IPv6
cluster, and hence creation virtual machines on top of such a cluster is not
possible.

## Diagnosis

- Check the cluster network configuration by running the following command:
  ```shell
  $ oc get network.config cluster -o yaml
  ```
  The output displays only an IPv6 CIDR for the cluster network.

  Example output:
  ```text
  apiVersion: config.openshift.io/v1
  kind: Network
  metadata:
    name: cluster
  spec:
    clusterNetwork:
    - cidr: fd02::/48
      hostPrefix: 64
  ```

## Mitigation

It is recommended to use single stack IPv4 or a dual stack IPv4/IPv6 networking
to use OpenShift Virtualization .Refer the
[documentation](https://docs.openshift.com/container-platform/latest/networking/ovn_kubernetes_network_provider/converting-to-dual-stack.html).