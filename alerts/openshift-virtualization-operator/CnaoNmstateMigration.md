# CnaoNmstateMigration
<!-- Edited by Jiří Herrmann, 8 Nov 2022 -->

## Meaning

This alert fires when a `kubernetes-nmstate` deployment is detected and the
Kubernetes NMState Operator is not installed. This alert only affects OpenShift
Virtualization 4.10.

The Cluster Network Addons Operator (CNAO) does not support `kubernetes-nmstate`
deployments in OpenShift Virtualization 4.11 and later.

## Impact

You cannot upgrade to OpenShift Virtualization 4.11.

## Mitigation

Install the Kubernetes NMState Operator from the OperatorHub. CNAO automatically
transfers the `kubernetes-nmstate` deployment to the Operator.
