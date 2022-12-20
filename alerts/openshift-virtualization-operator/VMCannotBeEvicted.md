# VMCannotBeEvicted
<!-- Edited by apinnick, Nov 2022-->

## Meaning

This alert fires when the eviction strategy of a virtual machine (VM) is set
to `LiveMigration` but the VM is not migratable.

## Impact

Non-migratable VMs prevent node eviction. This condition affects operations
such as node drain and updates.

## Diagnosis

1. Check the VMI configuration to determine whether the value of
`evictionStrategy` is `LiveMigrate`:

   ```bash
   $ oc get vmis -o yaml
   ```

2. Check for a `False` status in the `LIVE-MIGRATABLE` column to identify VMIs
that are not migratable:

   ```bash
   $ oc get vmis -o wide
   ```

3. Obtain the details of the VMI and check `spec.conditions` to identify the
issue:

   ```bash
   $ oc get vmi <vmi> -o yaml
   ```

   Example output:

   ```yaml
   status:
     conditions:
     - lastProbeTime: null
       lastTransitionTime: null
       message: cannot migrate VMI which does not use masquerade to connect
       to the pod network
       reason: InterfaceNotLiveMigratable
       status: "False"
       type: LiveMigratable
   ```

## Mitigation

Set the `evictionStrategy` of the VMI to `shutdown` or resolve the issue that
prevents the VMI from migrating.
