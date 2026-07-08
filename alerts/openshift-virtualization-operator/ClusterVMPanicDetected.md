# ClusterVMPanicDetected

## Meaning

This alert fires when one or more VMs across the cluster have experienced
non-recoverable guest OS panics in the last 24 hours. This may indicate a
cluster-wide infrastructure issue such as a faulty node image, a broken
driver, or a shared storage problem affecting VMs.

The alert is based on the `kubevirt_vmi_guest_os_panic_total` metric, which
tracks panic events detected via the pvpanic device (Linux and Windows) or
the Hyper-V enlightenment mechanism (Windows).

## Impact

* One or more VMs are crashing across the cluster.
* Applications running on affected VMs were unavailable during crashes.
* If VMs have `RunStrategy: Always`, they restart automatically but may
  continue crash-looping.
* If multiple VMs are affected, the issue likely relates to shared
  infrastructure rather than individual workloads.

## Diagnosis

1. **Identify all affected VMs** by querying the panic metric:

   ```promql
   sum by (namespace, name) (increase(kubevirt_vmi_guest_os_panic_total[24h])) > 0
   ```

2. **Look for common patterns** across affected VMs:

   ```promql
   kubevirt_vmi_guest_os_panic_total
   ```

   Check whether panics share the same `type` (e.g., all `pvpanic` or all
   `hyper-v`) or the same `bugcheck_code`.

3. **Check if affected VMs share a common node:**

   ```bash
   oc get vmi -A -o wide | grep -E "<vm-name-1>|<vm-name-2>|..."
   ```

   If all affected VMs run on the same node, the issue is likely
   node-specific (hardware, kernel, driver).

4. **Check node health and events:**

   ```bash
   oc describe node <node-name>
   oc get events -A --field-selector involvedObject.kind=Node
   ```

5. **Review virt-launcher logs** for any of the affected VMs:

   ```bash
   POD=$(oc get pod -n <namespace> -l kubevirt.io/domain=<vm-name> -o name | head -n1)
   oc logs $POD -n <namespace> -c compute --previous
   ```

## Mitigation

* **Immediately:** Identify whether panics are concentrated on specific
  nodes. If so, cordon the affected node(s) to prevent new VMs from
  scheduling there.
* **Short term:** Check for recent cluster-wide changes (node OS updates,
  driver updates, storage changes) that coincide with the panics.
* **Long term:** Investigate the root cause (faulty hardware, driver
  incompatibility, storage issue) and apply fixes across the affected
  infrastructure.

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.