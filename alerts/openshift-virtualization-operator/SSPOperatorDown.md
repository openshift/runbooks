# SSPOperatorDown [Deprecated]

This alert is deprecated. You can safely ignore or silence it.
(SSP) Operator
pods are down.

The SSP Operator is responsible for deploying and reconciling the common
templates and the Template Validator.

## Impact

Dependent components might not be deployed. Changes in the components might not
be reconciled. As a result, the common templates and/or the Template Validator
might not be updated or reset if they fail.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get deployment -A | grep ssp-operator | awk '{print $1}')"
   ```

2. Check the status of the `ssp-operator` pods.

   ```bash
   $ oc -n $NAMESPACE get pods -l control-plane=ssp-operator
   ```

3. Obtain the details of the `ssp-operator` pods:

   ```bash
   $ oc -n $NAMESPACE describe pods -l control-plane=ssp-operator
   ```

4. Check the `ssp-operator` logs for error messages:

   ```bash
   $ oc -n $NAMESPACE logs --tail=-1 -l control-plane=ssp-operator
   ```

## Mitigation

Try to identify the root cause and resolve the issue.
If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.

**Note:** Starting from 4.14, this runbook will no longer be supported. For a
supported runbook, please see [SSPDown
Runbook](http://kubevirt.io/monitoring/runbooks/SSPDown.html).