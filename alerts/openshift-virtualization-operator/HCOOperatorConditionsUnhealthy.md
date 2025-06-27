# HCOOperatorConditionsUnhealthy

## Meaning

This alert triggers when the HCO operator conditions or its secondary resources
are in an error or warning state.

## Impact

Resources maintained by the operator might not be functioning correctly.

## Diagnosis

Check the operator conditions:

```bash
oc get HyperConverged kubevirt-hyperconverged -n openshift-cnv -o jsonpath='{.status.conditions}'
```

## Mitigation

Based on the information obtained during the diagnosis procedure, try to
identify the root cause within the operator or any of its secondary resources,
and resolve the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.