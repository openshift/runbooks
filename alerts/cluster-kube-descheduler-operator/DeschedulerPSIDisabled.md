# DeschedulerPSIDisabled

## Meaning

This alert indicates that the Kube Descheduler Operator
is configured with a profile that requires PSI metrics (such as `KubeVirtRelieveAndMigrate`),
but PSI is not enabled on the worker nodes.

Pressure Stall Information (PSI) is a Linux kernel feature that quantifies
system saturation by tracking exactly how long tasks are stalled waiting for
CPU, memory, or I/O resources. The `KubeVirtRelieveAndMigrate` profile uses
these metrics to detect nodes under resource pressure and migrate workloads
to relieve them.

## Impact

The `KubeVirtRelieveAndMigrate` profile requires PSI metrics to be enabled
on all worker nodes. Without this, the descheduler cannot accurately detect
resource pressure or perform migrations to relieve overloaded nodes.

## Diagnosis

To confirm that PSI is disabled on your worker nodes, perform the following checks:

1. **Check for the PSI Interface**
   PSI metrics are exposed via the `/proc/pressure` directory.
   If this directory does not exist, PSI is disabled.

   * Open a debug session on a worker node:
        ```bash
        oc debug node/<node-name>
        ```
   * Once the debug pod starts, check for the existence of the pressure files:
        ```bash
        ls -l /host/proc/pressure/
        ```
   * **Result:** If you receive `No such file or directory`, PSI is **disabled**.
   * **Result:** If you see files like `cpu`, `io`, and `memory`, PSI is **enabled**.

2. **Verify Kernel Arguments**
   Check if the kernel was booted with the necessary arguments.
   * Still in the debug session, run:
       ```bash
       cat /host/proc/cmdline | grep psi
       ```
   * **Result:** If the output is empty or does not contain `psi=1`,
   the kernel argument is missing.

## Mitigation

PSI metrics can be enabled on worker nodes using a Machine Config.

**Note:** Applying this MachineConfig custom resource (CR) will trigger
a rolling reboot of your nodes.

```yaml
apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfig
metadata:
  labels:
    machineconfiguration.openshift.io/role: worker
  name: 99-openshift-machineconfig-worker-psi-karg
spec:
  kernelArguments:
    - psi=1
```

