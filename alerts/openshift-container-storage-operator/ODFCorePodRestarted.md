# ODFCorePodRestarted

## Meaning

A core ODF pod (OSD, MON, MGR, ODF operator, or metrics exporter) has
restarted at least once in the last 24 hours while the Ceph cluster is active.

## Impact

* If OSDs are restarted frequently or do not start up within 5 minutes,
  the cluster might decide to rebalance the data onto other more reliable
  disks. If this happens, the cluster will temporarily be slightly less
  performant.
* Operator restart delays configuration changes or health checks.
* May indicate underlying instability (resource pressure, bugs, or node issues).

## Diagnosis

1. Identify pod from alert (pod, namespace).
2. [pod debug](helpers/podDebug.md)

## Mitigation

1. If OOMKilled: Increase memory limits for the container.
2. If CrashLoopBackOff: Check for configuration errors or version incompatibilities.
3. If node-related: Cordon and drain the node; replace if faulty.
4. Ensure HA: MONs should be ≥3; OSDs should be distributed.
5. Update: If due to a known bug, upgrade ODF to a fixed version.
