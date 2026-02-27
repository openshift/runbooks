# KubemacpoolMACCollisionDetected

## Meaning

Multiple running workloads are using the same MAC address.

## Impact

MAC collisions cause network issues: packet loss, ARP table conflicts, and
traffic being delivered to the wrong targets.

## Diagnosis

1. Set the `KMP_NAMESPACE` environment variable:

   ```bash
   $ export KMP_NAMESPACE="$(oc get pod -A --no-headers -l \
      control-plane=mac-controller-manager | awk '{print $1}')"
   ```

2. Query the `kmp_mac_collisions` metric to see which MACs are colliding
   (value > 1 means collision):

   ```bash
   $ oc exec -n $KMP_NAMESPACE deployment/kubemacpool-mac-controller-manager \
      -c manager -- curl -s http://localhost:8080/metrics | grep kmp_mac_collisions
   ```

3. For each colliding MAC, find the VMIs involved:

   ```bash
   $ oc get vmi -A -o jsonpath='{range .items[*]}{.metadata.namespace}{"\t"}{.metadata.name}{"\t"}{.status.interfaces[*].mac}{"\n"}{end}' | grep -i "<MAC_ADDRESS>"
   ```

## Mitigation

Remove the collision by deleting or reconfiguring one of the colliding
workloads to use a different MAC address.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.