# KubeVirtVMGuestMemoryPressure

## Meaning
The VM guest OS is under sustained memory pressure (low usable memory
with elevated major page faults and/or swap I/O), risking thrashing,
swapping, or OOM kills.

## Impact
- Performance degradation due to page thrashing
- Increased swap I/O and latency
- Risk of OOM kills and application instability

## Diagnosis
- Metrics (Prometheus)
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
- Short term (in guest)
  - If possible, Restart or tune memory-heavy processes; reduce workload.
  - If acceptable, drop caches temporarily.
  - Ensure swap is sized/policy-compliant (or disabled if not desired).

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