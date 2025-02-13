# HAControlPlaneDown

## Meaning

A control plane node has been detected as not ready for more than 5 minutes.

## Impact

When a control plane node is down, it affects the high availability and
redundancy of the OpenShift Container Platform control plane. This can
negatively impact:
- API server availability
- Controller manager operations
- Scheduler functionality
- etcd cluster health (if etcd is co-located)

## Diagnosis

1. Check the status of all control plane nodes:
   ```bash
   oc get nodes -l node-role.kubernetes.io/control-plane=''
   ```

2. Get detailed information about the affected node:
   ```bash
   oc describe node <node-name>
   ```

3. Review system logs on the affected node:
   ```bash
   ssh <node-address>
   ```

   ```bash
   journalctl -xeu kubelet
   ```

## Mitigation

1. Check node resources:
   - Verify CPU, memory, and disk usage
      ```bash
      # Check the node's CPU and memory resource usage
      oc top node <node-name>
      ```

      ```bash
      # Check node status conditions for DiskPressure status
      oc get node <node-name> -o yaml | jq '.status.conditions[] | select(.type == "DiskPressure")'
      ```
   - Clear disk space if necessary
   - Restart the kubelet if resource issues are resolved

2. If the node is unreachable:
   - Verify network connectivity
   - Check physical/virtual machine status
   - Ensure the node has power and is running

3. If the kubelet is generating errors:
   ```bash
   systemctl status kubelet
   ```

   ```bash
   systemctl restart kubelet
   ```

4. If the node cannot be recovered:
   - If possible, safely drain the node
      ```bash
      oc drain <node-name> --ignore-daemonsets --delete-emptydir-data
      ```
   - Investigate hardware/infrastructure issues
   - Consider replacing the node if necessary

## Additional notes
- Maintain at least three control plane nodes for high availability
- Monitor etcd cluster health if the affected node runs etcd
- Document any infrastructure-specific recovery procedures

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.