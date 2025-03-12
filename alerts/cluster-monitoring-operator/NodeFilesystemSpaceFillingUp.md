# NodeFilesystemSpaceFillingUp

## Meaning

This alert indicates that a filesystem on an OpenShift worker node (RHCOS) is
running out of space and is predicted to be full within a specific timeframe:

- Warning: predicted to run out of space within 24 hours and current usage > 60%
- Critical: predicted to run out of space within 4 hours and current usage > 85%

## Impact

- Potential service disruptions on the OpenShift worker node
- Pods may fail to create or run due to insufficient storage
- Image pulls may fail if `/var/lib/containers` is affected
- Node may become unschedulable if critical system partitions fill up
- etcd may crash if `/var/lib/etcd` fills up (on control plane nodes)

## Diagnosis

1. Identify which node and filesystem is affected through the alert labels:

   ```shell
   oc get nodes <NODE_NAME> -o wide
   ```

2. For critical system areas, check via debug pod:

   ```shell
   oc debug node/<NODE_NAME> -- chroot /host df -h
   ```

3. Identify large directories on the affected node:

   ```shell
   oc debug node/<NODE_NAME> -- \
     chroot /host du -h --max-depth=1 <mountpoint> | sort -hr
   ```

4. Check container-specific storage:

   ```shell
   oc debug node/<NODE_NAME> -- \
     chroot /host du -h --max-depth=1 \
     /var/lib/containers /var/lib/kubelet | sort -hr
   ```

5. Check for any large log files:

   ```shell
   oc debug node/<NODE_NAME> -- chroot /host find /var/log -type f -size +100M
   ```

## Mitigation

1. For container storage issues:

   - Clean up unused containers and images on the node:

   ```shell
   oc debug node/<NODE_NAME> -- chroot /host crictl rmi --prune
   ```

   - Remove unused pods:

   ```shell
   oc debug node/<NODE_NAME> -- \
     chroot /host crictl rm \
     $(crictl ps -a -q --state exited)
   ```

2. For log file issues:

   - RHCOS uses the systemd journal for logs, so you may need to clean up
     journal storage:

   ```shell
   oc debug node/<NODE_NAME> -- chroot /host journalctl --vacuum-time=1d
   ```

3. Long-term solutions:

   - Adjust MachineConfig to increase partition sizes
   - Consider using the Container Storage Interface (CSI) for pods with high
     storage needs
   - Configure proper log rotation for containers in the OpenShift logging stack
   - Increase storage capacity of critical volumes in the cluster infrastructure
   - Set appropriate resource quotas and limits for namespaces

4. For temporary relief, mark node as unschedulable to prevent new pods:

   ```shell
   oc adm cordon <NODE_NAME>
   ```

## Additional Notes

- RHCOS uses an immutable filesystem model for the OS components
- Modifications to system directories require proper MachineConfig resources
- The `/var` directory is one of the few writable areas for system services
- Container storage paths are most commonly filled due to accumulated images
- Regular monitoring of storage trends should be implemented across the cluster
