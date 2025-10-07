# GuestFilesystemAlmostOutOfSpace

## Meaning

This alert triggers when a file system in a virtual machine (VM) is running out
of available disk space.

The alert has the following severity levels:

- **Warning**: Triggers when file-system usage is between 85% and 95% for 5
minutes.
- **Critical**: Triggers when file-system usage exceeds 95%.

The alert provides details about the specific file system (`disk_name`), its
mount point (`mount_point`), the VMI name, and namespace.

## Impact

When a guest file system runs out of space, the VM may experience:

- Application failures or crashes due to inability to write data
- System instability or inability to perform updates
- Potential data loss if applications cannot write to disk
- Degraded performance as the system struggles with low disk space
- The VM may become unresponsive or enter a failed state

Critical space exhaustion (>95%) significantly increases the risk of immediate
application or system failure.

## Diagnosis

1. Check the file-system usage metrics for the VMI:

   ```bash
   # Get the domain name of the VMI
   $ oc exec -it <pod-name> -n <namespace> -- virsh list
   ```

   ```bash
   # Query filesystem info via QEMU Guest Agent
   $ oc exec -it <pod-name> -n <namespace> -- virsh qemu-agent-command <domain-name> '{"execute": "guest-get-fsinfo"}'
   ```

2. You can also connect to the VMI by using the `virtctl` console to inspect
disk usage in the guest OS:

   ```bash
   $ virtctl console <vmi-name> -n <namespace>
   ```

   In the guest OS, run commands to evaluate file-system usage. For example, in
   Linux guests:

   ```bash
   $ df -h
   $ du -sh /* | sort -h
   ```

## Mitigation

### Immediate Actions

1. **Free up space in the guest file system**:

   - Remove temporary files, logs, or caches
   - Clean up old application data or archives
   - Remove unnecessary packages or applications

2. **Expand the disk size** if the underlying PVC supports volume expansion:

   a. Check if the storage class supports volume expansion:

      ```bash
      $ oc get storageclass <storage-class-name> -o yaml | grep allowVolumeExpansion
      ```

   b. Modify the PVC to request a larger size:

      ```bash
      $ oc edit pvc <pvc-name> -n <namespace>
      # Update the spec.resources.requests.storage value to a larger size
      ```

   c. Restart the VM to apply the changes:

      ```bash
      $ virtctl restart <vm-name> -n <namespace>
      ```

      This step causes the disk.img to be resized to match the new PVC size.
      While this is primarily needed for filesystem volumes, it's recommended
      to always restart for consistency.

### Long-term Solutions

1. **Set up log rotation and cleanup policies** within guest operating systems.

2. **Use appropriate volume sizes** when provisioning VMs based on expected
workload requirements.

3. **Consider using dynamic storage provisioning** with StorageClasses that
support volume expansion.

4. **Implement automated cleanup scripts** or maintenance tasks within VMs to
manage disk space proactively.

5. **Review and optimize application storage patterns** to minimize unnecessary
disk usage.

If the issue persists or the file system continues to fill up rapidly after
cleanup, investigate the root cause such as application bugs, excessive logging,
or unexpected data growth.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.