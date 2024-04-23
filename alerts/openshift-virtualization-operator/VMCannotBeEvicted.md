# VMCannotBeEvicted

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

If possible, resolve the issue preventing the VM from migrating.
Configure the VM to shut down during node drains or pod eviction by setting
`evictionStrategy: None` in the VM manifest.
