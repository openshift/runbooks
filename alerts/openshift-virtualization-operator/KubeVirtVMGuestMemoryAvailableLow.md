# KubeVirtVMGuestMemoryAvailableLow

## Meaning
The virtual machine (VM) guest operating system (OS) has very low available
memory
(<3% headroom) for an extended period with no meaningful swap I/O, approaching
OOM
conditions.

## Impact
- Very limited memory headroom (risk of imminent OOM)
- Potential application failures and degraded performance

## Diagnosis
 1) Identify the VMI pod (`vmi_pod` label):
     ```promql
     kubevirt_vmi_info{vm="<vm-name>", namespace="<namespace>",
     phase="running"}
     ```
 2) Headroom (usable/available) is < 3% sustained:
     ```promql
     sum by (name, namespace) (
       kubevirt_vmi_memory_usable_bytes{
         name="<vmi_pod>", namespace="<namespace>"
       }
     ) /
     sum by (name, namespace) (
       kubevirt_vmi_memory_available_bytes{
         name="<vmi_pod>", namespace="<namespace>"
       }
     )
     ```
 3) Swap traffic ~0 is over 30m (likely no swap):
     ```promql
     sum by (name, namespace) (
       rate(
         kubevirt_vmi_memory_swap_in_traffic_bytes{
           name="<vmi_pod>", namespace="<namespace>"
         }[30m]
       ) +
       rate(
         kubevirt_vmi_memory_swap_out_traffic_bytes{
           name="<vmi_pod>", namespace="<namespace>"
         }[30m]
       )
     )
     ```
 4) Major page faults are ~0 over 30m:
     ```promql
     sum by (name, namespace) (
       rate(
         kubevirt_vmi_memory_pgmajfault_total{
           name="<vmi_pod>", namespace="<namespace>"
         }[30m]
       )
     )
     ```

## Mitigation
- Immediate actions (in guest OS):
  - Stop or restart memory-hungry processes to reduce load.
  - Drop caches temporarily.
  - Enable and configure swap per policy (or verify that it was intentionally
disabled).

- Increase the VM memory (use a hotplug if supported, otherwise you might need
to restart):

1) Edit the VM:
    ```bash
    # Edit the VM and adjust spec.template.spec.domain.resources.{requests,limits}.memory
    oc edit vm <vm-name> -n <namespace>
    ```

2) If you need to restart to apply the change, stop the VM when appropriate:
    ```bash
    # Stop the VM (you can start it again from your usual workflow)
    virtctl stop <vm-name> -n <namespace>
    ```

3) After you update the memory, restart the VM:
    ```bash
    virtctl start <vm-name> -n <namespace>
    ```

If you cannot resolve the issue, log in to the [Red Hat Customer
Portal](link:https://access.redhat.com)
and open a support case, attaching the artifacts gathered during the diagnosis
procedure.