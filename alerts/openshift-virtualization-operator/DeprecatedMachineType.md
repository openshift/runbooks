# DeprecatedMachineType

## Meaning
This alert triggers when one or more Virtual Machines (VMs) are using machine
types that have been marked as no longer supported.

## Impact

**Running VMs**
- Continue running but are using an unsupported machine type.
- Newer nodes might not support the machine type, so the next live migration may
  fail or become unresponsive.
- This can block node maintenance and disrupt high availability.

**Stopped VMs**
- No immediate issue while stopped.
- When restarted, they might fail to schedule or become unresponsive if newer
  nodes no longer support the machine type.
- This might prevent workloads from coming back online after maintenance
  or cluster upgrades.

## Diagnosis
The alert detects VMs using unsupported machine types.

**Identify affected VMs**
Use the alert description to locate VM names, namespaces, and nodes (if
running).

**Root Cause:**
The VM's `spec.template.spec.domain.machine.type` field is set to a type
that has been marked as unsupported. This can happen due to:

- VMs created before the removal of support for a machine type.
- VM templates not updated after cluster upgrades.
- Manual VM creation using old machine type references.

## Mitigation
Update affected VMs to use a supported machine type. You can do one of the
following:

- Edit VM definitions individually by modifying the
  `spec.template.spec.domain.machine.type` field.
- For a smoother and cleaner update of multiple VMs, use the
  `kubevirt-api-lifecycle-automation` tool to transition all deprecated VMs
  in one operation. This ensures consistent, automated migration and reduces
  manual errors or downtime during cluster upgrades. For details, see [Updating
multiple

VMs](https://docs.redhat.com/en/documentation/openshift_container_platform/latest/html/virtualization/managing-vms#virt-updating-multiple-vms_virt-edit-vms).

**Important:** Plan and apply these updates before performing cluster
upgrades to avoid VM restart failures or compatibility issues.

Alternative: remove the machine field to let OpenShift Virtualization default
to the newest
supported type:

```bash
# Removes spec.template.spec.domain.machine so the mutating webhook defaults it
oc patch vm <vm-name> -n <namespace> --type='json' \
  -p='[{"op":"remove","path":"/spec/template/spec/domain/machine"}]'
```

Note: The automation tool performs these JSON patches
in bulk across multiple VMs.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.