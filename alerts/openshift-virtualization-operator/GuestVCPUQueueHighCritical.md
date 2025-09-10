# GuestVCPUQueueHighCritical

## Meaning
A VirtualMachineInstance (VMI) reported a
**guest CPU run‑queue length greater than 20** runnable or
uninterruptible threads within the last scrape window (120s).
This indicates severe CPU contention.

## Impact
* Sustained backlog; high latency and throughput degradation are likely.
* Risk of timeouts, watchdog resets, or I/O amplification.

## Diagnosis

1. **Confirm queue length**
   ```promql
   kubevirt_vmi_guest_vcpu_queue{namespace="$NS",name="$VM"}
   ```
   The longer the queue is over 20, the more significant issues it might cause.

2. **Check host CPU usage**
   ```promql
   rate(kubevirt_vmi_cpu_usage_seconds_total{namespace="$NS",name="$VM"}[2m])
   ``
   If the CPU usage is more than 90%, migrate other VMs or
   mark the node as unschedulable, so that no new pods/VMs are placed on it.

3. **Inspect guest processes**
   `virtctl console <vm>` → `top -H` or `pidstat -u 1`

4. **Verify vCPU allocation**
   ```bash
   oc get vmi $VM -ojsonpath='{.spec.domain.cpu}'
   ```

## Mitigation

* **Immediately:** Live-migrate VM, hot-plug vCPUs, stop or throttle hot
threads.
* **Short term:** Raise vCPU limit or split workload across additional VMs.
* **Long term:** Adjust placement rules; add autoscaling tied to run-queue
length.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.