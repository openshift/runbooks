# ODFNodeNICBandwidthSaturation

## Meaning

A network interface on an ODF node is operating at >90% of its reported
link speed, indicating potential bandwidth saturation.

## Impact

* Network congestion leading to packet drops or latency spikes.
* Slowed Ceph replication, backfill, and recovery.
* Client I/O timeouts or stalls.
* Possible Ceph OSD evictions due to heartbeat failures.

## Diagnosis

1. From alert, note instance and device.
2. Check current utilization:
```bash
oc debug node/<node>
sar -n DEV 1 5
```
3. Use Prometheus to graph:
```promql
rate(node_network_receive_bytes_total{instance="<ip>", device="<dev>"}[5m]) * 8
rate(node_network_transmit_bytes_total{...}) * 8
```
4. Determine if traffic is Ceph-related (e.g., during rebalance) or external.

## Mitigation

1. Short term: Throttle non-essential traffic on the node.
    * Taint the OSD node to prevent scheduling of non-storage workloads.
    * Drain existing non-essential pods from the node.
2. Long term:
    * Upgrade to higher-speed NICs (e.g., 25GbE → 100GbE).
    * Use multiple bonded interfaces with LACP.
    * Separate storage and client traffic using VLANs or dedicated NICs.
3. Tune Ceph osd_max_backfills, osd_recovery_max_active to reduce
   recovery bandwidth.
4. Enable NIC offload features (TSO, GRO) if disabled.
