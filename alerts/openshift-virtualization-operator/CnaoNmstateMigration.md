# CnaoNmstateMigration

## Meaning

This alert fires when a `kubernetes-nmstate` deployment is detected and the
OpenShift Container Platform NMState Operator is not installed. This alert only
affects OpenShift
Virtualization 4.10.

The Cluster Network Addons Operator (CNAO) does not support `kubernetes-nmstate`
deployments in OpenShift Virtualization 4.11 and later.

## Impact

You cannot upgrade your cluster to OpenShift Virtualization 4.11.

## Mitigation

Install the OpenShift Container Platform NMState Operator from the OperatorHub.
CNAO automatically
transfers the `kubernetes-nmstate` deployment to the Operator.

Afterwards, you can upgrade to OpenShift Virtualization 4.11.