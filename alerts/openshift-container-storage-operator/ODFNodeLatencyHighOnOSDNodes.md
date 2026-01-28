# ODFNodeLatencyHighOnOSDNodes

## Meaning

ICMP round-trip time (RTT) latency between ODF monitoring probes and
OSD nodes exceeds 10 milliseconds over the last 24 hours. This alert
triggers only on nodes that host Ceph OSD pods, indicating potential
network congestion or issues on the storage network.

## Impact

* Increased latency in Ceph replication and recovery operations.
* Higher client I/O latency for RBD and CephFS workloads.
* Risk of OSDs being marked down if heartbeat timeouts occur.
* Degraded cluster performance and possible client timeouts.


## Diagnosis

1. Check the alert’s instance label to get the node IP.
2. From a monitoring or debug pod, test connectivity:
```bash
ping <node-internal-ip>
```
3. Use mtr or traceroute to analyze path and hops.
4. Verify if the node is under high CPU or network load:
```bash
oc debug node/<node>
top -b -n 1 | head -20
sar -u 1 5
```
5. Check Ceph health and OSD status:
```bash
ceph osd status
ceph -s
```

## Mitigation

1. Isolate traffic: Confirm storage traffic uses a dedicated VLAN or NIC, separate
    from management/tenant traffic.
2. Hardware check: Inspect switch logs, NIC errors (ethtool -S <iface>),
    and NIC firmware.
3. Topology: Ensure OSD nodes are in the same rack/zone or connected via
    low-latency fabric.
4. If latency is transient, monitor; if persistent, engage network or
    infrastructure team.

