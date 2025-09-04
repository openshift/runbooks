# DeprecatedMachineType

## Meaning
This alert fires when one or more Virtual Machines (VMs)
are using machine types that have been marked as deprecated (no
longer supported).

## Impact

**Running VMs**
- Continue running but are using a deprecated machine type.
- Newer nodes may not support this type, so the next live migration may
  fail or get stuck.
- This can block node maintenance and disrupt high availability.

**Stopped VMs**
- No immediate issue while stopped.
- When restarted, they may fail to schedule or get stuck if newer nodes
  no longer support the deprecated machine type.
- This can prevent workloads from coming back online after maintenance
  or cluster upgrades.

## Diagnosis
The alert detects VMs using deprecated machine types.

**Identify affected VMs**
   Use the alert description to locate VM names, namespaces, and nodes (if
running).

**Root Cause:**
The VM's `spec.template.spec.domain.machine.type` field is set to a type
that has been marked as deprecated. This can happen due to:

- VMs created before a machine type deprecation.
- VM templates not updated after cluster upgrades.
- Manual VM creation using old machine type references.

## Mitigation
Update affected VMs to use a supported machine type. You can:

- Edit VM definitions individually by modifying the
  `spec.template.spec.domain.machine.type` field.
- Or, for a smoother and cleaner update of multiple VMs, use the
  `kubevirt-api-lifecycle-automation` tool to transition all deprecated VMs
  in one operation. This ensures consistent, automated migration and reduces
  manual errors or downtime during cluster upgrades. See the documentation
  for details: [Updating multiple
VMs](https://docs.redhat.com/en/documentation/openshift_container_platform/4.18/html/virtualization/managing-vms#virt-updating-multiple-vms_virt-edit-vms)

**Important:** Plan and apply these updates before performing cluster
upgrades to avoid VM restart failures or compatibility issues.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.