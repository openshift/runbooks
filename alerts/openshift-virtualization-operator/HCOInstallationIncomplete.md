# HCOInstallationIncomplete

## Meaning
This alert fires when the HyperConverged Cluster Operator (HCO) runs for more
than an hour without a `HyperConverged` custom resource (CR).

This alert has the following causes:

- During the installation process, you installed the HCO but you did not create
the `HyperConverged` CR.
- During the uninstall process, you removed the `HyperConverged` CR before
uninstalling the HCO and the HCO is still running.

## Mitigation

The mitigation depends on whether you are installing or uninstalling
the HCO:

- Complete the installation by creating a `HyperConverged` CR with its
default values:

  * In version v4.23.0
    or above, use the `v1` API version:
    ```bash
    $ cat <<EOF | oc apply -f -
    apiVersion: hco.kubevirt.io/v1
    kind: HyperConverged
    metadata:
      name: kubevirt-hyperconverged
      namespace: openshift-cnv
    spec: {}
    EOF
    ```

  * In versions earlier than v4.23.0,
    use the `v1beta1` API version:
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

- Uninstall the HCO. If the uninstallation process continues to run, you must
resolve that issue in order to cancel the alert.

If you cannot resolve the issue, log in to the
[Red Hat Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.