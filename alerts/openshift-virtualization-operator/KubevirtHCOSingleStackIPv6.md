# KubevirtHCOSingleStackIPv6

## Meaning

This alert fires when you install OpenShift Virtualization on a single stack
IPv6 cluster.

## Impact

You cannot create virtual machines.

## Diagnosis

- Obtain the network details for the cluster by running the following command:
  ```shell
  $ oc describe network cluster
  ```
  The network section contains an IPv6 CIDR.

## Mitigation

Install OpenShift Virtualization on a single stack IPv4 cluster or a dual stack
IPv4/IPv6 cluster.
