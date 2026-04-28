# UnsupportedHCOModification

## Meaning

This alert fires when a JSON Patch annotation is used to change an operand of
the HyperConverged Cluster Operator (HCO).

HCO configures OpenShift Virtualization and its supporting operators in an
opinionated way and
overwrites its operands when there is an unexpected change to them. Users must
not modify the operands directly.

However, if a change is required and it is not supported by the HCO API, you can
force HCO to set a change in an operator by using JSON Patch annotations. These
changes are not reverted by HCO during its reconciliation process.

## Impact

Incorrect use of JSON Patch annotations might lead to unexpected results or an
unstable environment.

Upgrading a system with JSON Patch annotations is dangerous because the
structure of the component custom resources might change.

## Diagnosis

Check the `annotation_name` in the alert details to identify the JSON Patch
annotation:

  ```text
  Labels
    alertname=UnsupportedHCOModification
    annotation_name=kubevirt.kubevirt.io/jsonpatch
    severity=info
  ```

To locate and review JSON Patch annotations on the `HyperConverged` resource:

  ```bash
  # Discover the HCO namespace (commonly 'openshift-cnv' or 'kubevirt-hyperconverged')
  NS=$(oc get hyperconverged -A -o jsonpath='{.items[0].metadata.namespace}')

  # Show annotations on the HyperConverged CR
  oc get hyperconverged -n "${NS}" kubevirt-hyperconverged \
    -o jsonpath='{.metadata.annotations}'; echo
  ```

When reviewing the annotations:
- The jsonpatch annotation keys are fixed and live on the `HyperConverged` CR:
  - `kubevirt.kubevirt.io/jsonpatch`
  - `containerizeddataimporter.kubevirt.io/jsonpatch`
  - `networkaddonsconfigs.kubevirt.io/jsonpatch`
  - `ssp.kubevirt.io/jsonpatch`
- Annotations that don't include the word `jsonpatch` are normal
  OpenShift Virtualization annotations and are unrelated to this alert.

## Mitigation

It is best to use the HCO API to change an operand. However, if the change can
only be done with a JSON Patch annotation, proceed with caution.

Remove JSON Patch annotations before upgrade to avoid potential issues.

To remove a specific JSON Patch annotation from the `HyperConverged` CR:

  ```bash
  # Example: remove the KubeVirt jsonpatch annotation key
  oc annotate hyperconverged -n "${NS}" kubevirt-hyperconverged \
    'kubevirt.kubevirt.io/jsonpatch-'

  # Other examples (use if those annotations are present):
  # Remove CDI jsonpatch
  # oc annotate hyperconverged -n "${NS}" kubevirt-hyperconverged 'containerizeddataimporter.kubevirt.io/jsonpatch-'
  # Remove CNAO jsonpatch
  # oc annotate hyperconverged -n "${NS}" kubevirt-hyperconverged 'networkaddonsconfigs.kubevirt.io/jsonpatch-'
  # Remove SSP jsonpatch
  # oc annotate hyperconverged -n "${NS}" kubevirt-hyperconverged 'ssp.kubevirt.io/jsonpatch-'
  ```