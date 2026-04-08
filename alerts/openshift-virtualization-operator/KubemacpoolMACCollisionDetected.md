# KubemacpoolMACCollisionDetected

## Meaning

Multiple running workloads are using the same MAC address.

## Impact

MAC collisions cause network issues: packet loss, ARP table conflicts, and
traffic being delivered to the wrong targets.

## Diagnosis

1. Query the `kmp_mac_collisions` metric to see which MACs are colliding
   (value > 1 means collision):



   **Note:** You can use the Openshift metrics explorer available at
'https://{OPENSHIFT_BASE_URL}/monitoring/query-browser'.

   ```promql
   $ kmp_mac_collisions > 1
   ```

2. For each colliding MAC, find the VMIs involved:

   ```bash
   $ export MAC=<MAC_ADDRESS>
   $ oc get vmi -A -o jsonpath='{range .items[*]}{.metadata.namespace}{"\t"}{.metadata.name}{"\t"}{.status.interfaces[*].mac}{"\n"}{end}' | grep -i "$MAC"
   ```

4. For each colliding MAC, find Pods with Multus secondary interfaces involved:

   ```bash
   $ export MAC=<MAC_ADDRESS>
   $ oc get pod -A -o json | jq -r --arg mac "$MAC" '.items[] | .metadata as $m | (.metadata.annotations["k8s.v1.cni.cncf.io/network-status"] // "[]" | fromjson)[] | select(.mac // "" | test($mac; "i")) | "\($m.namespace)\t\($m.name)\t\(.mac)"'
   ```

## Mitigation

Remove the collision by deleting or reconfiguring one of the colliding
workloads to use a different MAC address.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.