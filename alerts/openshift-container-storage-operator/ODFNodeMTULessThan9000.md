# ODFNodeMTULessThan9000

## Meaning

At least one physical or relevant network interface on an ODF node has an
MTU (Maximum Transmission Unit) less than 9000 bytes, violating ODF best
practices for storage networks.

## Impact

* Suboptimal Ceph network performance due to increased packet overhead.
* Higher CPU utilization on OSD nodes from processing more packets.
* Potential for packet fragmentation if mixed MTU sizes exist in the path.
* Reduced throughput during rebalancing or recovery.


## Diagnosis

1. List all nodes in the storage cluster:
```bash
oc get nodes -l cluster.ocs.openshift.io/openshift-storage=''
```
2. For each node, check interface MTUs:
```bash
oc debug node/<node-name>
ip link show
# Look for interfaces like eth0, ens*, eno*, etc. (exclude veth, docker, cali)
```
3. Alternatively, use Prometheus:
```promql
node_network_mtu_bytes{device!~"^(veth|docker|flannel|cali|tun|tap).*"} < 9000
```
4. Verify MTU consistency across all nodes and all switches in the storage fabric.

## Mitigation

1. Ensure the node network interfaces are configured for 9000 bytes
2. Ensure switches in between the nodes support 9000 bytes on their ports.