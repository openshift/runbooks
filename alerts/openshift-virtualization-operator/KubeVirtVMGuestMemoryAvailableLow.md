# KubeVirtVMGuestMemoryAvailableLow

## Meaning
The VM guest OS has very low available memory (<3% headroom) for an
extended period with no meaningful swap I/O, approaching OOM conditions.

## Impact
- Very limited memory headroom (risk of imminent OOM)
- Potential application failures and degraded performance

## Diagnosis
- Metrics (Prometheus)
  1) Identify the VMI pod (`vmi_pod` label):
  ```promql
  kubevirt_vmi_info{vm="<vm-name>", namespace="<namespace>",
  phase="running"}
  ```
  2) Headroom (usable/available) < 3% sustained:
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
  3) Swap traffic ~0 over 30m (likely no swap):
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
  4) Major page faults ~0 over 30m:
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
- Immediate actions (in guest)
  - If possible, stop or restart memory-hungry processes to reduce load.
  - If acceptable, drop caches temporarily.
  - Configure/enable swap per policy (or validate it is intentionally disabled).

- Increase VM memory (hotplug if supported; otherwise restart may be required)

To increase VM memory in the VM spec:

  ```bash
  # Edit the VM and adjust spec.template.spec.domain.resources.{requests,limits}.memory
  oc edit vm <vm-name> -n <namespace>
  ```

If a restart is required to apply the change, gracefully stop the VM when
appropriate:

  ```bash
  # Stop the VM (you can start it again from your usual workflow)
  virtctl stop <vm-name> -n <namespace>
  ```

Start the VM after updating the memory:

  ```bash
  virtctl start <vm-name> -n <namespace>
  ```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.