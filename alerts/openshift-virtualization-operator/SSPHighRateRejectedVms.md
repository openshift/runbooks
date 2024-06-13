# SSPHighRateRejectedVms

## Meaning

This alert fires when a user or script attempts to create or modify a large
number of virtual machines (VMs), using an invalid configuration.

## Impact

The VMs are not created or modified. As a result, the environment might not
behave as expected.

## Diagnosis

1. Export the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get deployment -A | grep ssp-operator | awk '{print $1}')"
   ```

2. Check the `virt-template-validator` logs for errors that might indicate the
cause:

   ```bash
   $ oc -n $NAMESPACE logs --tail=-1 -l name=virt-template-validator
   ```

   Example output:

   ```text
   {"component":"kubevirt-template-validator","level":"info","msg":"evalution
   summary for ubuntu-3166wmdbbfkroku0:\nminimal-required-memory applied: FAIL,
   value 1073741824 is lower than minimum [2147483648]\n\nsucceeded=false",
   "pos":"admission.go:25","timestamp":"2021-09-28T17:59:10.934470Z"}
   ```

## Mitigation

Try to identify the root cause and resolve the issue.
If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.