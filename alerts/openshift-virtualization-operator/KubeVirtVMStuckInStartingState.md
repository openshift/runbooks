# KubeVirtVMStuckInStartingState
<!-- Edited by apinnick, Nov 2022 -->

## Meaning

This alert fires when a virtual machine (VM) is in a starting state for more
than 5 minutes.

This alert might indicate an issue in the VM configuration, such as a
misconfigured priority class or a missing network device.

## Impact

There is no immediate impact. However, if this alert persists, you should try
to resolve the issue.

## Diagnosis

- Check the virtual machine instance (VMI) details for error conditions:

  ```bash
  $ oc describe vmi <vmi> -n <namespace>
  ```

  Example output:

  ```yaml
  Name:          testvmi-ldgrw
  Namespace:     kubevirt-test-default1
  Labels:        name=testvmi-ldgrw
  Annotations:   kubevirt.io/latest-observed-api-version: v1
                 kubevirt.io/storage-observed-api-version: v1alpha3
  API Version:   kubevirt.io/v1
  Kind:          VirtualMachineInstance
  ...
  Spec:
  ...
    Networks:
      Name:  default
      Pod:
    Priority Class Name:               non-preemtible
    Termination Grace Period Seconds:  0
  Status:
    Conditions:
      Last Probe Time:       2022-10-03T11:08:30Z
      Last Transition Time:  2022-10-03T11:08:30Z
      Message:               virt-launcher pod has not yet been scheduled
      Reason:                PodNotExists
      Status:                False
      Type:                  Ready
      Last Probe Time:       <nil>
      Last Transition Time:  2022-10-03T11:08:30Z
      Message:               failed to create virtual machine pod: pods
      "virt-launcher-testvmi-ldgrw-" is forbidden: no PriorityClass with name
      non-preemtible was found
      Reason:                FailedCreate
      Status:                False
      Type:                  Synchronized
    Guest OS Info:
    Phase:  Pending
    Phase Transition Timestamps:
      Phase:                        Pending
      Phase Transition Timestamp:   2022-10-03T11:08:30Z
    Runtime User:                    0
    Virtual Machine Revision Name:
      revision-start-vm-6f01a94b-3260-4c5a-bbe5-dc98d13e6bea-1
  Events:
    Type     Reason        Age                From                       Message
    ----     ------        ----               ----                       -------
    Warning  FailedCreate  8s (x13 over 28s)  virtualmachine-controller  Error
    creating pod: pods "virt-launcher-testvmi-ldgrw-" is forbidden: no
    PriorityClass with name non-preemtible was found
  ```

## Mitigation

Ensure that the VM is configured correctly and has the required resources.

A `Pending` state indicates that the VM has not yet been scheduled. Check the
following possible causes:

- The `virt-launcher` pod is not scheduled.
- Topology hints for the VMI are not up to date.
- Data volume is not provisioned or ready.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
