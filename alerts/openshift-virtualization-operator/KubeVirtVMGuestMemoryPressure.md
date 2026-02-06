# KubeVirtVMGuestMemoryPressure

## Meaning
The virtual machine (VM) guest operating system (OS) is under sustained memory
pressure
(low usable memory with elevated major page faults and/or swap I/O). This can
cause
thrashing, swapping, or OOM kills.

## Impact
- Performance degradation due to page thrashing
- Increased swap I/O and latency
- Risk of OOM kills and application instability

## Diagnosis
1) Identify the VMI pod (`vmi_pod` label):
    ```promql
    kubevirt_vmi_info{vm="<vm-name>", namespace="<namespace>",
    phase="running"}
    ```
2) Headroom (usable/available) < 5% indicates pressure:
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
3) Major page faults are elevated:
    ```promql
    sum by (name, namespace) (
      rate(
        kubevirt_vmi_memory_pgmajfault_total{
          name="<vmi_pod>", namespace="<namespace>"
        }[5m]
      )
    )
    ```
4) Swap traffic is elevated (bytes/s):
    ```promql
    sum by (name, namespace) (
      rate(
        kubevirt_vmi_memory_swap_in_traffic_bytes{
          name="<vmi_pod>", namespace="<namespace>"
        }[5m]
      ) +
      rate(
        kubevirt_vmi_memory_swap_out_traffic_bytes{
          name="<vmi_pod>", namespace="<namespace>"
        }[5m]
      )
    )
    ```

## Mitigation
- Short term (in guest OS):
  - Restart or tune memory-heavy processes; reduce workload.
  - Drop caches temporarily.
  - Ensure swap is sized/policy-compliant (or disabled if not desired).

- Increase the VM memory (use a hotplug if supported, otherwise you might need
to restart):
  1) Edit the VM:
        ```bash
      # Edit the VM and adjust spec.template.spec.domain.resources.{requests,limits}.memory
      $ oc edit vm <vm-name> -n <namespace>
        ```

  2) If you need to restart to apply the change, stop the VM when appropriate:
        ```bash
      # Stop the VM (you can start it again from your usual workflow)
      $ virtctl stop <vm-name> -n <namespace>
        ```

  3) After you update the memory, restart the VM:
        ```bash
      $ virtctl start <vm-name> -n <namespace>
        ```

If you cannot resolve the issue, log in to the [Red Hat Customer
Portal](link:https://access.redhat.com)
and open a support case, attaching the artifacts gathered during the diagnosis
procedure.