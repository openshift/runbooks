# GuestVCPUQueueHighCritical

## Meaning
A VirtualMachineInstance (VMI) reported a
**guest CPU run‑queue length greater than 20** runnable or
uninterruptible threads within the last scrape window (120s),
indicating severe CPU contention.

## Impact
* Sustained backlog; high latency and throughput degradation are likely.
* Risk of timeouts, watchdog resets, or I/O amplification.

## Diagnosis
Follow the below steps with extra focus on:
* **Duration** – How long is the queue > 20
* **Host saturation** – if node CPU is also > 90 %, migrate other VMs or
mark the
   node as unschedulable so that no new Pods/VMs are placed thereon it.

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
| Horizon  | Action                                                           |
|----------|------------------------------------------------------------------|
| Immediate| **Prioritise**: live-migrate VM; hot-plug vCPUs; stop or throttle
hot threads.                                            |
| Short term| Raise vCPU limit or split workload across additional VMs.       |
| Long term| Adjust placement rules; add autoscaling tied to run-queue length.|

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.