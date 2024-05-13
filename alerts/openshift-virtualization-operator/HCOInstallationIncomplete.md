# HCOInstallationIncomplete

## Meaning

This alert fires when the HyperConverged Cluster Operator (HCO) runs for
more than an hour without a `HyperConverged` custom resource (CR).

This alert has the following causes:

- During the installation process, you installed the HCO but you did not
create the `HyperConverged` CR.
- During the uninstall process, you removed the `HyperConverged` CR before
uninstalling the HCO and the HCO is still running.

## Mitigation

The mitigation depends on whether you are installing or uninstalling
the HCO:

- Complete the installation by creating a `HyperConverged` CR with its
default values:

  ```bash
  $ cat <<EOF | oc apply -f -
  apiVersion: hco.kubevirt.io/v1beta1
  kind: HyperConverged
  metadata:
    name: kubevirt-hyperconverged
    namespace: openshift-cnv
  spec: {}
  EOF
  ```

- Uninstall the HCO. If the uninstall process continues to run, you must
resolve that issue in order to cancel the alert.
