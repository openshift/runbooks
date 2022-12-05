# OutdatedVirtualMachineInstanceWorkloads
<!-- Edited by apinnick Nov 2022-->

## Meaning

This alert fires when running virtual machine instances (VMIs) in
outdated `virt-launcher` pods are detected 24 hours after the OpenShift
Virtualization control plane has been updated.

## Impact

Outdated VMIs might not have access to new OpenShift Virtualization
features.

Outdated VMIs will not receive the security fixes associated with
the `virt-launcher` pod update.

## Diagnosis

1. Identify the outdated VMIs:

   ```bash
   $ oc get vmi -l kubevirt.io/outdatedLauncherImage --all-namespaces
   ```

2. Check the `KubeVirt` custom resource (CR) to determine whether
`workloadUpdateMethods` is configured in the `workloadUpdateStrategy`
stanza:

   ```bash
   $ oc get kubevirt kubevirt --all-namespaces -o yaml
   ```

3. Check each outdated VMI to determine whether it is live-migratable:

   ```bash
   $ oc get vmi <vmi> -o yaml
   ```

   Example output:

   ```yaml
   apiVersion: kubevirt.io/v1
   kind: VirtualMachineInstance
   ...
     status:
       conditions:
       - lastProbeTime: null
         lastTransitionTime: null
         message: cannot migrate VMI which does not use masquerade
         to connect to the pod network
         reason: InterfaceNotLiveMigratable
         status: "False"
         type: LiveMigratable
   ```

## Mitigation

### Configuring automated workload updates

Update the `HyperConverged` CR to enable automatic workload updates.

### Stopping a VM associated with a non-live-migratable VMI

- If a VMI is not live-migratable and if `runStrategy: always` is
set in the corresponding `VirtualMachine` object, you can update the
VMI by manually stopping the virtual machine (VM):

  ```bash
  $ virctl stop --namespace <namespace> <vm>
  ```

A new VMI spins up immediately in an updated `virt-launcher` pod to
replace the stopped VMI. This is the equivalent of a restart action.

Note: Manually stopping a _live-migratable_ VM is destructive and
not recommended because it interrupts the workload.

### Migrating a live-migratable VMI

If a VMI is live-migratable, you can update it by creating a `VirtualMachineInstanceMigration`
object that targets a specific running VMI. The VMI is migrated into
an updated `virt-launcher` pod.

1. Create a `VirtualMachineInstanceMigration` manifest and save it
as `migration.yaml`:

   ```yaml
   apiVersion: kubevirt.io/v1
   kind: VirtualMachineInstanceMigration
   metadata:
     name: <migration_name>
     namespace: <namespace>
   spec:
     vmiName: <vmi_name>
   ```

2. Create a `VirtualMachineInstanceMigration` object to trigger the
migration:

   ```bash
   $ oc create -f migration.yaml
   ```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the Diagnosis procedure.
