# KubevirtHyperconvergedClusterOperatorInstallationNotCompletedAlert
<!-- Edited: apinnick@redhat.com, Aug 2022 -->

## Meaning
This alert fires when the HyperConverged Cluster Operator (HCO) runs for
more than an hour without a `HyperConverged` custom resource (CR).

This alert has the following causes:

- During the installation process, you installed the HCO but you did not
create the `HyperConverged` CR.
- During the uninstall process, you removed the `HyperConverged` CR before
uninstalling the HCO and the HCO is still running.

## Mitigation

- Complete the installation by creating a `HyperConverged` CR with its
default values:

  ```bash
  $ cat <<EOF | oc apply -f -
  apiVersion: operators.coreos.com/v1
  kind: OperatorGroup
  metadata:
    name: hco-operatorgroup
    namespace: kubevirt-hyperconverged
  spec: {}
  EOF
  ```

  _or_

- Uninstall the HCO. If the uninstall process continues to run, you must
resolve that issue in order to cancel the alert.
