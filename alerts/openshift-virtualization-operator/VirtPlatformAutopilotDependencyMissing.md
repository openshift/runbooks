# VirtPlatformAutopilotDependencyMissing

## Meaning

This alert fires when virt-platform-autopilot detects that an optional CRD it
would use (a soft dependency) is missing from the cluster.

The alert is driven by the `kubevirt_autopilot_missing_dependency` gauge,
labelled by `group`, `version`, and `kind`:

- `1` = the CRD is absent
- `0` = the CRD is present

The alert triggers on `kubevirt_autopilot_missing_dependency == 1` held for
`5m`. When a soft-dependency CRD is missing, the operator skips the assets that
require it and continues to manage everything else ("inform, don't crash").

## Impact

This is a `warning` alert. The cluster remains functional but is
feature-incomplete: the platform features that rely on the missing CRD (for
example, load-aware scheduling when the descheduler CRD is absent) are not
configured until the CRD is installed.

## Diagnosis

1. Read the `group`, `version`, and `kind` of the missing CRD from the alert
   labels.

2. Confirm the CRD is absent. Substitute the `kind` from the alert labels
   below. CRD names use the plural form (which cannot be derived
   deterministically), so this searches by the kind suffix:

   ```bash
   $ oc get crd -o name | grep -i '<kind>'
   ```

3. Decide whether the CRD is expected on this cluster. If it should be present,
   check the health of the operator that provides it. Replace `<provider-operator>`
   with the actual operator name (e.g., `descheduler` for the Kube Descheduler
Operator):

   ```bash
   $ oc get csv -A | grep -i <provider-operator>
   ```

## Mitigation

- **The feature is wanted:** install the operator that provides the CRD. For
  example, to enable load-aware scheduling install the Kube Descheduler
  Operator. First discover the operator package in the catalog:

  ```bash
  $ oc get packagemanifests -n openshift-marketplace | \
      grep cluster-kube-descheduler-operator
  ```

  Then install it from OperatorHub or by creating an OLM `Subscription` for the
  package; follow the operator's own documentation for the exact install steps.

  virt-platform-autopilot detects the CRD automatically once it appears and
  resumes managing the dependent assets. The alert resolves within a few
  minutes of `kubevirt_autopilot_missing_dependency` returning to `0`.

- **The feature is intentionally not installed** (for example, a development
  cluster): the alert is informational and safe to silence in Alertmanager. It
  does not affect any other platform functionality.

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.