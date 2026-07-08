# VMNonRecoverableOSPanic

## Meaning

This alert fires when a VM has experienced non-recoverable guest OS panics in
the last 24 hours. The alert is based on the
`kubevirt_vmi_guest_os_panic_total` metric, which tracks panic events for all
panic types (pvpanic, hyper-v, s390, etc.).

A non-recoverable guest OS panic indicates that the guest kernel or operating
system crashed and was unable to recover on its own (e.g. via kdump). The
crash is detected via the pvpanic device (Linux and Windows) or the Hyper-V
enlightenment mechanism (Windows), and reported by QEMU/libvirt to OpenShift
Virtualization.
In the current OpenShift Virtualization default configuration the VM is
destroyed on panic,
but with a `RunStrategy` of `Always` it will be automatically restarted.

## Impact

* The VMI transitions to `Failed` phase when a non-recoverable panic is
detected.
* If the VM has `RunStrategy: Always`, it is automatically restarted.
* Applications on the VM were unavailable during the crash and restart.
* Repeated panics can indicate OS, driver, or workload instability.

## Diagnosis

1. **Confirm the alert labels** (`namespace`, `name`) in Alertmanager or the
   monitoring console and set variables for the following steps:

   ```bash
   export NAMESPACE="<alert namespace label>"
   export VM_NAME="<alert name label>"
   ```

2. **Check VMI phase and events**

   ```bash
   oc get vmi -n $NAMESPACE $VM_NAME -o wide
   oc describe vmi -n $NAMESPACE $VM_NAME
   ```

   Look for `GuestPanicked` or `Stopped` events.

3. **Inspect the guest OS panic metric** (includes panic type and bugcheck
   code):

   ```promql
   kubevirt_vmi_guest_os_panic_total{namespace="$NAMESPACE", name="$VM_NAME"}
   ```

4. **Check the alert expression** to see the number of panics in the last
   24 hours:

   ```promql
   sum by (namespace, name) (increase(kubevirt_vmi_guest_os_panic_total{namespace="$NAMESPACE", name="$VM_NAME"}[24h]))
   ```

5. **Review virt-launcher logs** on the node where the VMI ran:

   ```bash
   POD=$(oc get pod -n $NAMESPACE -l kubevirt.io/domain=$VM_NAME -o name | head -n1)
   oc describe $POD -n $NAMESPACE
   oc logs $POD -n $NAMESPACE -c compute --previous
   oc logs $POD -n $NAMESPACE -c compute
   ```

6. **For Windows guests**, collect crash dumps or event logs from inside the
   guest after the VM is running again. The `bugcheck_code` label on the metric
   provides the Windows BSOD code.

## Mitigation

* **Immediately:** Confirm workload health; restart the VM if it is stuck in
  `Failed` and policy allows.
* **Short term:** Review recent guest changes (updates, drivers, workload
  deploys); reduce load if the panic is resource-related.
* **Long term:** Ensure the VMI spec includes an appropriate panic device
  (`pvpanic` for Linux, `hyperv` for Windows) to enable reliable crash
  detection.

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.