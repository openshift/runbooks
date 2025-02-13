# HighCPUWorkload

## Meaning

This alert fires when a node's CPU utilization exceeds 90% for more than 5
minutes.

## Impact

High CPU utilization can lead to:
- Degraded performance of applications running on the node
- Increased latency in request processing
- Potential service disruptions if CPU usage continues to climb

## Diagnosis

1. Identify the affected node:
   ```bash
   oc get nodes
   ```

2. Check the node's resource usage:
   ```bash
   oc describe node <node-name>
   ```

3. List pods that consume high amounts of CPU:
   ```bash
   oc top pods --all-namespaces --sort-by=cpu
   ```

4. Investigate specific pod details if needed:
   ```bash
   oc describe pod <pod-name> -n <namespace>
   ```

## Mitigation

1. If the issue was caused by a malfunctioning pod:
   - Consider restarting the pod
   - Check pod logs for anomalies
   - Review pod resource limits and requests

2. If the issue is system-wide:
   - Check for system processes that consume high amounts of CPU
   - Consider cordoning the node and migrating workloads
   - Evaluate if node scaling is needed

3. Long-term solutions to avoid the issue:
   - Implement or adjust pod resource limits
   - Consider horizontal pod autoscaling
   - Evaluate cluster capacity and scaling needs

## Additional notes
- Monitor the node after mitigation to ensure CPU usage returns to normal
- Review application logs for potential root causes
- Consider updating resource requests/limits if this is a recurring issue

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.