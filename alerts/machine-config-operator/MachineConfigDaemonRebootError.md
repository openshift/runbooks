# MCDRebootError

## Meaning

This alert is triggered when the Machine Config Daemon (MCD)
detects an error when trying to reboot a
node. If the MCD is unable to complete
the reboot within 5 minutes, the alert
will fire.

## Impact

If the MCD is unable to succesfully reboot the node,
any pending MachineConfig changes that would
require a reboot would not be propagated,
and the MachineConfig cluster operator would degrade.
## Diagnosis

If a node is stuck and cannot reboot,
you should check the pod logs via the following command,
replacing the `$DAEMONPOD` variable with the actual name
of your own `machine-config-daemon-*` pod on the stuck node.

```console
oc logs -f -n openshift-machine-config-operator $DAEMONPOD -c machine-config-daemon
```

When the MCD needs to reboot a node it will log the
following message via the `update.go` service in the
machine-config-daemon-* pod scheduled on the node.

```console
update.go:2641] Rebooting node
```

It creates a systemd unit via systemd-run that reboots the
nodes via "systemctl reboot".

If the reboot fails at any point then it will log:

```console
"failed to run reboot: %v", err
```

Where %v will be replaced with the exact error.

Then the fmt.Errorf message from the reboot command
will be returned to the `writer.go` service. Which
will then log the following and mark the node
as degraded.

```console
E0219 01:25:20.890006    8301 writer.go:226] Marking Degraded due to:reboot command failed, something is seriously wrong
```

This will increment the `mcd_reboots_failed_total` value by 1.

The MCD will continue to try to reboot the node. If the node is
unable to be rebooted and the `mcd_reboots_failed_total` remains
greater then 0 than the `MCDRebootError` alert will fire.

## Mitigation

When the MCD first logs the error it will log the exact
error it's encountering when trying to reboot the node.
For example:

```console
update.go:2641] failed to run reboot: exec: "systemd-run": executable file not found in $PATH
```

This error indicates that the `systemd-run` file cannot be
found in the /usr/bin/systemd-run $PATH and so the node
cannot reboot succesfully.

The error message will change depending on what is
preventing the reboot.

If a node is failing to reboot. It's best to
start with the journal logs on the node
to troubleshoot why it's failing to reboot.
You can do so with the `journalctl` command.