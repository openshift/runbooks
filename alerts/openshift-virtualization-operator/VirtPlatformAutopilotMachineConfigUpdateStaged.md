# VirtPlatformAutopilotMachineConfigUpdateStaged

## Meaning

Autopilot holds a desired MachineConfig update for at least five minutes while
waiting for a matching MachineConfigPool rollout. The alert selects
`kubevirt_autopilot_machineconfig_update_staged == 1` for `5m`.
The `machineconfig` and `pool` labels identify the staged update and matching
pool. This is an informational alert with no operator health impact.

## Impact

The live MachineConfig retains its previous configuration until the staged
update is released. This is intentional: Autopilot coalesces changes with an
existing rollout to reduce node reboots. Coalescing is best-effort, not a
guarantee of a single rollout: a matching pool permits an apply attempt when
it reports `Updating=True`, which can occur after the rollout is already
underway, or `RenderDegraded=True`, which means MCO cannot render the pool's
configuration at all, so holding the update has no benefit. Nodes that the
rollout has already updated may need to converge again once the staged
configuration is applied.

## Diagnosis

Inspect the named MachineConfigPool and MachineConfig:

```bash
$ oc get machineconfigpool <pool> -o yaml
$ oc get machineconfig <machineconfig> -o yaml
```

Check whether the pool is updating or paused. The informational alert alone
does not mean that the Machine Config Operator is unhealthy.

## Mitigation

No administrator action is required for normal staging. Autopilot releases the
update when a matching pool reports `Updating=True` or `RenderDegraded=True`;
the former may occur after the rollout is underway, and nodes already updated
by that rollout may then converge again. Follow the existing maintenance plan
for a deliberately paused pool; do not start a rollout solely to clear this
informational alert.

If the pool is degraded or a rollout is stuck, investigate the Machine Config
Operator's conditions and alerts. The staged-update alert clears when Autopilot
applies or discards the staged update and removes its metric series.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.