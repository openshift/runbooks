# OperatorConditionsUnhealthy [Deprecated]

This alert is deprecated. You can safely ignore or silence it.
y resources
are in an error or warning state.

## Impact

Resources maintained by the operator might not be functioning correctly.

## Diagnosis

Check the operator conditions:

```bash
oc get <CR> <CR_OBJECT> -n <namespace> -o jsonpath='{.status.conditions}'
```

For example:

```bash
oc get HyperConverged kubevirt-hyperconverged -n kubevirt -o jsonpath='{.status.conditions}'
```

## Mitigation

Based on the information obtained during the diagnosis procedure, try to
identify the root cause within the operator or any of its secondary resources,
and resolve the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.