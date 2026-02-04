# VirtualMachineStuckInUnhealthyState

## Meaning

This alert triggers when a virtual machine (VM) is in an unhealthy state
for more than 10 minutes and does not have an associated VMI
(VirtualMachineInstance).

The alert indicates that a VM is experiencing early-stage lifecycle issues
before a VMI can be successfully created. This typically occurs during the
initial phases of VM startup when OpenShift Virtualization is trying to
provision resources,
pull images, or schedule the workload.

**Affected states:**
- `Provisioning` - The VM is preparing resources (DataVolumes, PVCs)
- `Starting` - VM is attempting to start but no VMI exists yet
- `Terminating` - VM is being deleted but without an active VMI
- `Error` - Various scheduling, image, or resource allocation errors

## Impact

- **Severity:** Warning
- **User impact:** VMs cannot start or are stuck in error states
- **Business impact:** Workloads cannot be deployed, which affects application
  availability

## Possible causes

### Resource-related issues
- Insufficient cluster resources (CPU, memory, storage)
- Missing or misconfigured storage classes
- PVC provisioning failures
- DataVolume creation/import failures

### Image and registry issues
- Container image pull failures for containerDisk volumes
- Registry authentication problems
- Network connectivity issues to image registries
- Missing or corrupted VM disk images

### Scheduling and node issues
- No schedulable nodes available (all nodes cordoned/unschedulable)
- Insufficient resources like KVM/GPU on available nodes
- A mismatch between requested and available CPU models
- Node selector constraints cannot be satisfied
- Taints and tolerations preventing scheduling

### Configuration issues
- Invalid VM specifications (malformed YAML, unsupported features)
- Missing required Secrets or ConfigMaps
- Incorrect resource requests/limits
- Network configuration errors

## Diagnosis

1. Check VM status and events
    ```bash
    # Get VM details and status
    $ oc get vm <vm-name> -n <namespace> -o yaml

    # Check VM events for error messages
    $ oc describe vm <vm-name> -n <namespace>

    # Look for related events in the namespace
    $ oc get events -n <namespace> --sort-by='.lastTimestamp'
    ```

2. Verify resource availability
    ```bash
    # Check node resources and schedulability
    $ oc get nodes -o wide
    $ oc describe nodes

    # Check storage classes and provisioners
    $ oc get storageclass
    $ oc get pv,pvc -n <namespace>

    # For DataVolumes (if using)
    $ oc get datavolume -n <namespace>
    $ oc describe datavolume <dv-name> -n <namespace>
    ```

3. Check image availability (for containerDisk)
    ```bash
    # If using containerDisk, verify image accessibility from the affected node
    # Start a debug session on the node hosting the VM (or a representative node)
    $ oc debug node/<node-name> -it --image=busybox

    # Inside the debug pod, check which container runtime is used
    $ ps aux | grep -E "(containerd|dockerd|crio)"

    # For CRI-O/containerd clusters use crictl to pull the image
    $ crictl pull <vm-disk-image>

    # For Docker-based clusters (less common)
    $ docker pull <vm-disk-image>

    # Exit the debug session when done
    $ exit

    # Check image pull secrets if required
    $ oc get secrets -n <namespace>
    ```

4. Verify OpenShift Virtualization configuration
    ```bash
    # Discover the KubeVirt installation namespace
    $ export NAMESPACE="$(oc get kubevirt -A -o custom-columns="":.metadata.namespace)"

    # Check KubeVirt CR conditions (expect Available=True)
    $ oc get kubevirt -n "$NAMESPACE" \
      -o jsonpath='{range .items[*].status.conditions[*]}{.type}={.status}{"\n"}{end}'

    # Or check a single CR named 'kubevirt'
    $ oc get kubevirt kubevirt -n "$NAMESPACE" \
      -o jsonpath='{.status.conditions[?(@.type=="Available")].status}'

    # Verify virt-controller is running
    $ oc get pods -n "$NAMESPACE" \
      -l kubevirt.io=virt-controller

    # Check virt-controller logs for errors
    # Replace <virt-controller-pod> with a pod name from the list above
    $ oc logs -n "$NAMESPACE" <virt-controller-pod>

    # Verify virt-handler is running
    $ oc get pods -n "$NAMESPACE" \
      -l kubevirt.io=virt-handler -o wide

    # Check virt-handler logs for errors (daemonset uses per-node pods)
    # Replace <virt-handler-pod> with a pod name from the list above
    $ oc logs -n "$NAMESPACE" <virt-handler-pod>
    ```

5. Review VM specification
Inspect the following details in the VM's spec to catch common
misconfigurations:

- Disks and volumes (in spec.template.spec.domain.devices and volumes):
  - A bootable disk is defined using *bootOrder: 1*
  - Each disk name matches a volume name
  - Volume sources are valid: PVC, DataVolume, containerDisk, secret, or
    configMap

- Resources (in spec.template.spec.domain.resources):
  - Resource requests and limits are set and do not exceed node capacity

- Scheduling (in spec.template.spec):
  - nodeSelector, affinity, and tolerations are not overly restrictive

- Image pull configuration (in spec.template.spec.imagePullSecrets):
  - imagePullSecrets are configured if using a private image registry

- Power strategy (in spec.runStrategy or spec.running):
  - Only one of spec.runStrategy or spec.running is set, and it matches the
    desired behavior

## Mitigation

### Resource issues
- **Scale up the cluster** if the issue is insufficient resources
- **Create missing storage classes** or configure default storage
- **Resolve PVC/DataVolume failures**:
   ```bash
   $ oc get pvc -n <namespace>
   $ oc describe pvc <pvc-name> -n <namespace>
   ```

### Image issues
- **Verify image accessibility**:
   ```bash
   # Validate from the node
   $ oc debug node/<node-name> -it --image=busybox

   # Inside the debug pod, detect runtime and pull
   $ ps aux | grep -E "(containerd|dockerd|crio)"

   # For CRI-O/containerd clusters:
   $ crictl pull <image-name>

   # For Docker-based clusters (less common):
   $ docker pull <image-name>

   $ exit
   ```
- **Configure image pull secrets** if needed:
   ```bash
   $ oc create secret docker-registry <secret-name> \
     --docker-server=<registry-url> \
     --docker-username=<username> \
     --docker-password=<password>
   ```

### Scheduling issues
- **Review VM scheduling constraints** and relax if too restrictive:
   - nodeSelector, affinity, and tolerations
   - Required CPU model, host devices, or features

- **Verify that node taints and tolerations** allow scheduling:
   - Ensure the VM tolerates node taints that apply to target nodes

- **Ensure that nodes have required capabilities**:
   - KVM availability, CPU features, GPU, SR-IOV, or storage access

- If nodes were intentionally cordoned for maintenance, **uncordon** when
   appropriate:
   ```bash
   $ oc uncordon <node-name>
   ```

### Configuration issues resolution
- **Fix VM specification errors** based on the *oc describe* output:
   ```bash
   # Edit VM specification directly
   $ oc edit vm <vm-name> -n <namespace>

   # Or patch specific fields
   $ oc patch vm <vm-name> -n <namespace> --type='merge' \
     -p='{"spec":{"template":{"spec":{"domain":{"resources": \
       {"requests":{"memory":"2Gi"}}}}}}}'
   ```
- **Create missing secrets/configmaps**:
   ```bash
   $ oc create secret generic <secret-name> \
     --from-literal=key=value
   ```
- **Adjust resource requests** if they exceed node capacity

### Emergency workarounds
- **Restart the VM** to apply specification changes:
  ```bash
  # Restart the VM to pick up spec changes
  $ virtctl restart <vm-name> -n <namespace>
  ```
- **Scale down non-critical workloads** temporarily if resource
  constraints exist
- **Change storage class** if PVC provisioning fails:
  ```bash
  # Check current storage class status
  $ oc get storageclass
  $ oc describe storageclass <current-storage-class>

  # Look for PVC provisioning errors
  $ oc describe pvc <pvc-name> -n <namespace>

  # If seeing "no volume provisioner" or similar errors,
  # specify a working storage class in VM spec:
  # spec.dataVolumeTemplates[].spec.pvc.storageClassName:
  #   <working-class>
  ```

## Prevention

- **Resource planning:**
   - Monitor cluster resource utilization
   - Set appropriate VM guest resources in the VM domain guest spec.
   - Plan storage capacity and provisioning

- **Image management:**
   - Use local image registries where possible to reduce
     latency
   - Configure DataVolume import methods appropriately:
     * **Pod import method**: Images pulled to temporary pods (default)
     * **Node import method**: Images pulled directly to nodes
       (requires pre-pulling)
   - Pre-pull critical containerDisk images to nodes only if
     using node import method

- **Monitoring:**
   - Set up alerts for cluster resource exhaustion
   - Monitor storage provisioner health
   - Track VM startup success rates

- **Testing:**
   - Validate VM templates in development environments
   - Test VM deployments after cluster changes
   - Regularly verify image accessibility

## Escalation

Escalate to the cluster administrator if:
- Multiple VMs are affected simultaneously (possible cluster-wide issue)
- Issue persists after following resolution steps
- A malfunction of OpenShift Virtualization components is suspected
- You are unable to access system logs for further diagnosis
- You do not have enough permissions to run the diagnosis or mitigation steps

## Related alerts

- `VirtControllerDown` - May indicate controller issues preventing
  VM processing
- `LowKVMNodesCount` - Related to insufficient KVM-capable nodes
- `KubeVirtNoAvailableNodesToRunVMs` - Indicates no nodes available
  for VM scheduling

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support
case, attaching the artifacts gathered during the diagnosis
procedure.