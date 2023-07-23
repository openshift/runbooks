# KubevirtHcoSingleStackIPv6

## Meaning

This alert fires when you install OpenShift Virtualization on a single stack
IPv6 cluster.

## Impact

You cannot create virtual machines.

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

Install OpenShift Virtualization on a single stack IPv4 cluster or on a
dual stack IPv4/IPv6 cluster.
