# GuestVCPUQueueHighWarning

## Meaning
A VirtualMachineInstance (VMI) reported a
**guest CPU run‑queue length greater than 10** runnable or
uninterruptible threads within the most‑recent scrape window (120s).
The run‑queue length is derived from `guest_load_1m – vCPU_count`.

## Impact
* Moderate CPU contention inside the guest;
  latency may spike but workload still progresses.
* Early signal that the VM might need additional vCPUs or
  that a short‑lived process is causing bursts.

## Diagnosis
1. **Confirm queue length**
   ```promql
   kubevirt_vmi_guest_vcpu_queue{namespace="$NS",name="$VM"}
   ```
2. **Check host CPU usage**
   ```promql
   rate(kubevirt_vmi_cpu_usage_seconds_total{namespace="$NS",name="$VM"}[2m])
   ```
3. **Inspect guest processes**
   `virtctl console <vm>` → `top -H` or `pidstat -u 1`
4. **Verify vCPU allocation**
   ```bash
   oc get vmi $VM -ojsonpath='{.spec.domain.cpu}'
   ```

## Mitigation
| Horizon | Action |
|---------|--------|
| Immediate | Optionally live‑migrate the VM to a quieter node or throttle
noisy processes.   |
| Short term | Hot‑plug / increase vCPUs; tune application thread pools.
                |
| Long term  | Implement horizontal scaling (HPA/KEDA, VMReplicaSet); review
placement rules. |

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.