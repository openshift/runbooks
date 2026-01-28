# ODFNodeLatencyHighOnNONOSDNodes

## Meaning

ICMP RTT latency to non-OSD ODF nodes (e.g., MON, MGR, MDS, or client nodes)
exceeds 100 milliseconds over the last 24 hours. These nodes participate in
Ceph control plane or client access but do not store data.

## Impact

* Delayed Ceph monitor elections or quorum instability.
* Slower metadata operations in CephFS.
* Increased latency for CSI controller operations.
* Potential timeouts in ODF operator reconciliation.
* Not support if it is a permanent configuration.


## Diagnosis

1. From the alert, note the instance (node IP).
2. Test connectivity:
```bash
ping <node-ip>
mtr <node-ip>
```
3. Check system load and network interface stats on the node:
```bash
oc debug node/<node-name>
sar -n DEV 1 5
ip -s link show <iface>
```
4. Review Ceph monitor logs if the node hosts MONs:
```bash
oc logs -l app=rook-ceph-mon -n openshift-storage
```
5. switch network monitoring to see if any ports are too busy.

## Mitigation

1. Ensure control-plane nodes are not oversubscribed or co-located with noisy workloads.
2. Validate network path between MON/MGR nodes—prefer low-latency, dedicated links.
3. If node is a client (e.g., running applications), verify it’s not on an
   overloaded subnet.
4. Tune kernel network parameters if packet loss or buffer drops are observed.
