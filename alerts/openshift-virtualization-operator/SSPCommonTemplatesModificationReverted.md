# SSPCommonTemplatesModificationReverted

## Meaning

This alert fires when the Scheduling, Scale, and Performance (SSP) Operator
reverts changes to common templates as part of its reconciliation procedure.

The SSP Operator deploys and reconciles the common templates and the Template
Validator. If a user or script changes a common template, the changes are
reverted by the SSP Operator.

## Impact

Changes to common templates are overwritten.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get deployment -A | grep ssp-operator | awk '{print $1}')"
   ```

2. Check the `ssp-operator` logs for templates with reverted changes:

   ```bash
   $ oc -n $NAMESPACE logs --tail=-1 -l control-plane=ssp-operator | grep 'common template' -C 3
   ```

## Mitigation

Try to identify and resolve the cause of the changes.

Ensure that changes are made only to copies of templates, and not to the
templates themselves.

<!-- No downstream link. Modules cannot contain xrefs.-->