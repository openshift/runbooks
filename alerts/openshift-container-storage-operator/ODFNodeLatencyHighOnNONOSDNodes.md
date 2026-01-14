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


## Diagnosis

1. From the alert, note the instance (node IP).
2. Confirm the node does not run OSDs:
```bash
oc get pods -n openshift-storage -o wide | grep <node-name>
```
3. Test connectivity:
```bash
ping <node-ip>
mtr <node-ip>
```
4. Check system load and network interface stats on the node:
```bash
oc debug node/<node-name>
sar -n DEV 1 5
ip -s link show <iface>
```
5. Review Ceph monitor logs if the node hosts MONs:
```bash
oc logs -l app=rook-ceph-mon -n openshift-storage
```


## Mitigation

1. Ensure control-plane nodes are not oversubscribed or co-located with noisy workloads.
2. Validate network path between MON/MGR nodes—prefer low-latency, dedicated links.
3. If node is a client (e.g., running applications), verify it’s not on an
   overloaded subnet.
4. Tune kernel network parameters if packet loss or buffer drops are observed.
