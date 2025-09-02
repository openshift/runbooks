# VirtualMachineStuckOnNode

## Meaning

This alert fires when a VirtualMachine with an associated VMI
(VirtualMachineInstance) has been stuck in an unhealthy state for more
than 5 minutes on a specific node.

The alert indicates that a VirtualMachine has progressed past initial
scheduling and has an active VMI, but is experiencing runtime issues on
the assigned node. This typically occurs after the VM has been scheduled
to a node but encounters problems during startup, operation, or shutdown
phases.

**Affected States:**
- `Starting` - VMI exists but VM is failing to reach running state
- `Stopping` - VM is attempting to stop but the process is stuck
- `Terminating` - VM is being deleted but the termination process is hanging
- `Error` states - Runtime errors occurring on the node (ErrImagePull,
  ImagePullBackOff, etc.)

## Impact

- **Severity:** Warning
- **User Impact:** VMs are unresponsive or stuck in transition states
- **Business Impact:** Running workloads may be disrupted, affecting
  application performance and availability
- **Node Impact:** Resources may be tied up by stuck VMs, affecting
  other workloads on the same node

## Possible Causes

### Node-Level Issues
- **Node resource exhaustion** (CPU, memory, storage)
- **Container runtime problems** (containerd, CRI-O issues)
- **Insufficient storage** on the node
- **Network connectivity issues** from the node
- **Node entering NotReady** or maintenance mode

### QEMU/KVM Issues
- **QEMU process failures** or hangs
- **KVM acceleration problems** on the node
- **Nested virtualization** configuration issues
- **Hardware compatibility** problems

### Image and Storage Issues
- **Container image pull failures** specific to the node
- **Local image cache corruption**
- **PVC mount failures** on the node
- **Storage backend connectivity** issues from the node
- **Volume attachment timeouts**

### virt-launcher Pod Issues
- **virt-launcher pod** stuck in non-ready state
- **Pod resource limits** being exceeded
- **Security policy violations** (SELinux, AppArmor)
- **Networking problems** within the pod

### libvirt/Domain Issues
- **libvirt daemon** problems on the node
- **Domain definition** conflicts or corruption
- **Migration failures** (if VM was being migrated)
- **Hot-plug operations** that failed and left VM in inconsistent
  state

## Diagnosis

### 1. Check VM and VMI Status
```bash
# Get VM details with node information
oc get vm <vm-name> -n <namespace> -o yaml

# Check VMI status and node assignment
oc get vmi <vm-name> -n <namespace> -o yaml
oc describe vmi <vm-name> -n <namespace>

# Look for related events
oc get events -n <namespace> \
  --field-selector involvedObject.name=<vm-name>
```

### 2. Examine virt-launcher Pod
```bash
# Find the virt-launcher pod for this VM
oc get pods -n <namespace> -l kubevirt.io/domain=<vm-name>

# Check pod status and events
oc describe pod <virt-launcher-pod> -n <namespace>

# Check pod logs for errors
oc logs <virt-launcher-pod> -n <namespace> -c compute
oc logs <virt-launcher-pod> -n <namespace> -c istio-proxy \
  # if using Istio

# Optional: Check resource usage for the virt-launcher pod
oc top pod <virt-launcher-pod> -n <namespace>
```

### 3. Investigate Node Health
```bash
# Check node status and conditions (may require admin
# permissions)
oc describe node <node-name>

# Discover the KubeVirt installation namespace
export NAMESPACE="$(oc get kubevirt -A -o custom-columns="":.metadata.namespace)"

# Check virt-handler on the affected node
oc get pods -n "$NAMESPACE" -o wide | grep <node-name>
oc logs <virt-handler-pod> -n "$NAMESPACE"
```

### 4. Check Storage and Volumes
```bash
# Verify PVC status and mounting
oc get pvc -n <namespace>
oc describe pvc <pvc-name> -n <namespace>

# Check volume attachments on the node
oc get volumeattachment | grep <node-name>

# For DataVolumes, check their status
oc get dv -n <namespace>
oc describe dv <dv-name> -n <namespace>
```

### 5. Verify Image Accessibility from Node
```bash
# Verify image accessibility from the affected node
oc debug node/<node-name> -it --image=busybox

# Inside the debug pod, check which container runtime is used
ps aux | grep -E "(containerd|dockerd|crio)"

# For CRI-O/containerd clusters:
crictl pull <vm-disk-image>

# For Docker-based clusters (less common):
docker pull <vm-disk-image>

# Exit the debug session when done
exit
```

### 6. Exec into the virt‑launcher pod’s compute container and inspect
domains
```bash
oc exec -it <virt-launcher-pod> -n <namespace> -c compute \
  -- virsh list --all | grep <vm-name>
oc exec -it <virt-launcher-pod> -n <namespace> -c compute \
  -- virsh dumpxml <domain-name>

```

## Mitigation

### Pod-Level Issues
1. **Restart virt-launcher pod**:
   ```bash
   oc delete pod <virt-launcher-pod> -n <namespace>
   # The VMI controller will recreate it
   ```

2. **Check resource constraints**:
   ```bash
   oc describe pod <virt-launcher-pod> -n <namespace>
   # Look for resource limit violations
   ```

### Image Issues on Node
1. **Inspect and, if necessary, clear image cache** on the node:
   ```bash
   # SSH to the node or start a debug session on the node:
   oc debug node/<node-name> -it --image=busybox

   # Detect which container runtime is in use
   ps aux | grep -E "(containerd|dockerd|crio)"

   # List cached images first
   # For CRI-O/containerd clusters:
   crictl images
   # For Docker-based clusters:
   docker images

   # Remove only if a corrupted/stale image is suspected
   # For CRI-O/containerd clusters:
   crictl rmi <problematic-image>
   # For Docker-based clusters:
   docker rmi <problematic-image>

   exit
   ```

2. **Force image re-pull**:
   ```bash
   # Delete and recreate the virt-launcher pod
   oc delete pod <virt-launcher-pod> -n <namespace>
   ```

### Storage Issues
1. **Check PVC binding and mounting**:
   ```bash
   oc get pvc -n <namespace>
   # If PVC is stuck, check the storage provisioner
   ```

2. **Resolve volume attachment issues**:
   ```bash
   oc get volumeattachment
   # Delete stuck volume attachments if necessary
   oc delete volumeattachment <attachment-name>
   ```

### Node-Level Issues Resolution
1. **Drain and uncordon node** if it's in a bad state:
   ```bash
   oc drain <node-name> --ignore-daemonsets \
     --delete-emptydir-data
   oc uncordon <node-name>
   ```

2. **Restart node-level components**:
   ```bash
   # Restart virt-handler on the node
   oc delete pod <virt-handler-pod> -n "$NAMESPACE"
   ```

### VM-Level Resolution
1. **Force‑delete the VMI (will trigger new VMI creation)**:
   ```bash
   oc delete vmi <vm-name> -n <namespace> --force \
     --grace-period=0
   ```

2. **Migrate VM to different node**:
   ```bash
   virtctl migrate <vm-name> -n <namespace>
   ```

### Emergency Actions
- **Live migrate** critical VMs away from problematic node
- **Force delete** stuck VMI if safe to do so:
  ```bash
  oc delete vmi <vm-name> -n <namespace> --force \
    --grace-period=0
  ```
- **Cordon node** to prevent new VM scheduling while
  investigating

## Prevention

1. **Node Health Monitoring:**
   - Monitor node resource utilization (CPU, memory,
     storage)
   - Set up alerts for node conditions and taints
   - Regular health checks on container runtime

2. **Resource Management:**
   - Set appropriate resource requests/limits on VMs
   - Monitor PVC and storage utilization
   - Plan for node capacity and VM density

3. **Image Management:**
   - Use image pull policies appropriately (Always,
     IfNotPresent)
   - Pre-pull critical images to nodes
   - Monitor image registry health and connectivity

4. **Networking:**
   - Ensure stable network connectivity between nodes and
     storage
   - Monitor DNS resolution and service discovery
   - Validate network policies don't block required
     traffic

5. **Regular Maintenance:**
   - Keep nodes and OpenShift Virtualization components updated

## Escalation

Escalate to the cluster administrator if:
- Multiple VMs affected on the same node simultaneously
- VMs consistently fail to start or become unresponsive after
  following troubleshooting steps
- Node‑specific issues persist (Like kubelet or kernel panic) and
  require a node reboot
- Suspected OpenShift Virtualization component bugs or regressions
- Unable to access system logs for further diagnosis
- You don't have enough permissions to run the diagnosis and/or
mitigation steps.

## Related Alerts

- `OrphanedVirtualMachineInstances` - May indicate virt-handler
  problems on nodes
- `VirtHandlerDown` - Related to virt-handler pod failures
- `VirtualMachineStuckInUnhealthyState` - For VMs that haven't
  progressed to having VMIs

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support
case, attaching the artifacts gathered during the diagnosis
procedure.