# V4SubnetAllocationThresholdExceeded

## Meaning

The `V4SubnetAllocationThresholdExceeded` alert is triggered when more than
80% of subnets for nodes are allocated.

## Impact

This is a warning alert. No immediate impact to the cluster will be observed if
this alert fires and it is a warning to be mindful of your remaining node
subnet allocation. If your remaining subnets are exhausted, then no
further nodes can be added to your cluster.

## Diagnosis

Check the network configuration on the cluster.

    oc get networks.config.openshift.io/cluster -o jsonpath='{.spec.clusterNetwork}'

    [{"cidr":"10.128.0.0/14","hostPrefix":23}]

Calculate the IPv4 subnets capability.

    subnet_capability = 2^[(32 - clusternetwork_netmask) - (32 - hostPrefix)]

It will be 512 if the CIDR netmask is `/14` and hostPrefix is `23`, that means
the cluster can have at most 512 nodes.

Count the number of nodes to compare.

    oc get node --no-headers | wc -l

## Mitigation

### For clusters running prior to 4.13

We do not support adding additional cluster networks for ovn-kubernetes.

User will have to create a new cluster for more worker nodes.

Choosing a larger cluster network CIDR which can hold more subnets could prevent
this happening.

### For clusters running with version 4.13 onward

User can follow the OpenShift document to expand the cluster subnet IP address range.

Refer: https://docs.redhat.com/en/documentation/openshift_container_platform/4.13/html/release_notes/ocp-4-13-release-notes#ocp-4-13-expand-cluster-network-ip-address-range_release-notes
