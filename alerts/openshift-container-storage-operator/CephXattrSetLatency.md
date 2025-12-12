# CephXattrSetLatency

## Meaning

This alert indicates that the Ceph Metadata Server (MDS) is experiencing high
latency when setting extended attributes (xattr) on files. The alert fires when
the average setxattr operation latency exceeds 30 milliseconds over a 5-minute
period.

**What are extended attributes (xattrs)?**

Extended attributes are named key/value metadata entries stored alongside
filesystem objects (inodes, directories, symlinks) in CephFS. They follow POSIX
conventions with namespace prefixes:

- **user.*** — Application-specific metadata
- **security.*** — SELinux labels and security contexts
- **system.*** — POSIX ACLs and system metadata
- **trusted.*** — Administrative attributes (requires CAP_SYS_ADMIN)

**What does setxattr do?**

The MDS performs setxattr operations on behalf of clients to write or update
extended attributes. This involves:

- Permission and capability checks
- Namespace validation
- In-memory metadata update
- Journal entry creation and durable commit

## Impact

**Severity:** Warning

High xattr set latency can cause:

- Slow file system operations, particularly for applications that rely heavily
  on extended attributes (e.g., SELinux, NFSv4 ACLs, backup tools)
- Degraded overall CephFS performance
- Application timeouts or failures when interacting with the file system
- Increased latency for file creation and modification operations
- Stalled workflows for operations like save, chmod, backup/restore

**Typical latency ranges:**

| Condition | Expected Latency |
| --------- | ---------------- |
| Light-load LAN with SSDs | 1–10 ms |
| Busy clusters or HDD-backed metadata | 10–100+ ms |
| Under contention, failover, or WAN | 100 ms to seconds |

## Diagnosis

### Step 1: Check MDS Status and Health

Access the Ceph tools pod and check the overall cluster and MDS health:

```bash
oc rsh -n openshift-storage $(oc get pods -n openshift-storage -l app=rook-ceph-tools -o name)
```

Run the following commands:

```bash
ceph status
ceph mds stat
ceph fs status
```

Look for any warnings related to slow metadata IOs or MDS health issues.

### Step 2: Check MDS Performance Metrics

Examine the MDS operations in flight:

```bash
ceph daemon mds.<mds-name> dump_ops_in_flight
```

To find the active MDS name:

```bash
ceph fs status -f json-pretty | jq -r '.mdsmap[] | select(.state=="active") | .name'
```

### Step 3: Check MDS CPU and Memory Usage

Using the OpenShift console, go to Workloads -> Pods and select the MDS pod
(e.g., `rook-ceph-mds-ocs-storagecluster-cephfilesystem-*`). Click on the
Metrics tab to review CPU and memory usage.

Alternatively, check MDS resource usage:

```bash
oc adm top pod -n openshift-storage -l app=rook-ceph-mds
```

### Step 4: Check for Network Issues

Network latency between MDS and OSDs can cause slow metadata operations.
The client-to-MDS RPC round-trip typically adds 0.5–5 ms on LAN, but can be
significantly higher on loaded or WAN links.

Follow the steps in the
[Check Ceph Network Connectivity SOP](helpers/networkConnectivity.md)
to verify network health.

### Step 5: Check OSD Performance

Slow OSD operations can cascade into MDS latency issues. The metadata
write-to-durable-store step depends on underlying storage performance:

- NVMe: ~0.1–1 ms
- SSD: ~1–5 ms
- HDD: significantly higher

Check for slow OSD operations:

```bash
ceph health detail
ceph osd perf
```

### Step 6: Check for Lock Contention

If an inode is locked or requires cross-MDS coordination (cap flushing,
referrals, recovery), latency can increase by 10s to 100s of milliseconds
in pathological cases.

Check for blocked operations:

```bash
ceph daemon mds.<mds-name> dump_blocked_ops
```

## Mitigation

### Recommended Actions

1. **Increase MDS CPU Resources:**

   If the MDS CPU usage is consistently high, increase the allocated CPU.
   MDS is largely single-threaded, so higher clock speed CPUs are more
   effective than additional cores:

   ```bash
   oc patch -n openshift-storage storagecluster ocs-storagecluster \
       --type merge \
       --patch '{"spec": {"resources": {"mds": {"limits": {"cpu": "8"}, "requests": {"cpu": "8"}}}}}'
   ```
   **Note:** If the above step doesn't resolve the issue,
   that is the CPU usage remains high even after the above change,
   request the next higher power of two (16 CPUs), and repeat as needed
   (32, 64, etc.).

2. **Increase MDS Cache Memory:**

   If the MDS cache is under pressure, increase the memory allocation:

   ```bash
   oc patch -n openshift-storage storagecluster ocs-storagecluster \
       --type merge \
       --patch '{"spec": {"resources": {"mds": {"limits": {"memory": "8Gi"}, "requests": {"memory": "8Gi"}}}}}'
   ```

   **Note:** ODF sets `mds_cache_memory_limit` to half of the MDS pod memory
   request/limit. Setting memory to 8GB results in a 4GB cache limit.

3. **Scale Out with Multiple Active MDS:**

   For high metadata workloads, consider running multiple active MDS instances
   to reduce lock contention and distribute metadata operations:

   ```bash
   oc patch -n openshift-storage storagecluster ocs-storagecluster \
       --type merge \
       --patch '{"spec": {"managedResources": {"cephFilesystems": {"activeMetadataServers": 2}}}}'
   ```

   Always increase `activeMetadataServers` by 1. This is effective when
   metadata load is distributed across multiple directories/PVs.

4. **Use Faster Metadata Storage:**

   If the metadata pool is backed by HDDs, consider migrating to SSD or NVMe
   storage for improved journal/WAL commit latency.

5. **Address Network Issues:**

   If network connectivity issues are identified, escalate to the network
   or infrastructure team. Optimize network by lowering RTT and ensuring
   adequate bandwidth. See
   [Check Ceph Network Connectivity SOP](helpers/networkConnectivity.md).

6. **Address Underlying OSD Issues:**

   If OSDs are slow, investigate and resolve OSD performance problems first.
   Slow OSDs directly impact MDS performance. Refer to
   [CephOSDSlowOps runbook](CephOSDSlowOps.md) for guidance.

7. **Restart MDS (if stuck operations detected):**

   If operations appear stuck due to internal issues, restarting the MDS may
   help:

   ```bash
   oc delete pod -n openshift-storage -l app=rook-ceph-mds
   ```

   The pod will be automatically recreated by the operator.

If the issue persists after taking the above actions, please contact Red Hat
Support for further assistance.

## Additional Resources

- [Red Hat Ceph Storage Troubleshooting Guide](https://docs.redhat.com/en/documentation/red_hat_ceph_storage/9/html/troubleshooting_guide/index)
