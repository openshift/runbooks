# MachineConfigControllerPoolAlert

## Meaning

Alerts the user to when the Machine Config
Controller (MCC) detects when nodes are
a part of the master pool and a custom pool.
Causing only the master pool machine configs to
apply.

This alert will fire as a warning when:
- One or more nodes have overlapping
  labels that match both the master pool and a
  custom pool.
- This alert will automatically fire
  immediately if the condition is met.

## Impact
If the master and custom pool apply to a node then
the MCC will make a choice of honoring only the
master pool. This will result
in custom machine configurations not
being applied.

## Diagnosis

If the alert is firing then check
the MCC pod logs in the `openshift-machine-config-operator`
namespace with the following command.

For the following command, replace the
$CONTROLLERPOD variable
with the name of your own
`machine-config-controller` pod name.

```console
oc -n openshift-machine-config-operator logs $CONTROLLERPOD -c machine-config-controller
```
If a node is apart the master pool and a custom pool
then the following will be logged:

```console
Found xxxxxx node that matches selector for custom pool yyyyy, defaulting to xxxxx. This node will not have any custom role configuration as a result. Please review the node to make sure this is intended
```

## Mitigation

Please be aware that while a `master`, `custom`
MCP setup is supported for nodes that it will
activate this alert. Any MCs matching the custom
pool will not be applied to nodes matching the
master pool. If this is not intentional then you
can check the MCP object and it's spec to edit
the selector/labels.

Machine Config Pools will apply a machine
config to nodes that match the spec.nodeSelector.matchLabels.

```console
apiVersion: machineconfiguration.openshift.io/v1
kind: MachineConfigPool
metadata:
  name: <custom-pool>
spec:
  machineConfigSelector:
    matchExpressions:
      - {key: machineconfiguration.openshift.io/role, operator: In, values: [value, value]}
  maxUnavailable: null
  nodeSelector:
    matchLabels:
      key:value
  paused: false
```
