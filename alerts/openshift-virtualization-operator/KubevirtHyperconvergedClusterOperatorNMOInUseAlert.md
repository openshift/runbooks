# KubevirtHyperconvergedClusterOperatorNMOInUseAlert
<!-- Edited by apinnick, Nov 2022-->

## Meaning

This alert fires when _integrated_ Node Maintenance Operator (NMO) custom
resources (CRs) are detected. This alert only affects OpenShift Virtualization
4.10.

The Node Maintenance Operator is not included with OpenShift Virtualization
4.11.0 or later. Instead, the Operator is installed from OperatorHub.

The presence of `NodeMaintenance` CRs belonging to the `nodemaintenance.kubevirt.io`
API group indicates that the node specified in `spec.nodeName` was put
into maintenance mode. The target node has been cordoned off and drained.


## Impact

You cannot upgrade to OpenShift Virtualization 4.11.

## Diagnosis

1. Check the `kubevirt-hyperconverged` resource to see whether it is upgradeable:

   ```bash
   $ oc get hco kubevirt-hyperconverged -o json | \
     jq -r '.status.conditions[] | select(.type == "Upgradeable")'
   ```

   Example output:

   ```json
   {
     "lastTransitionTime": "2022-05-26T09:23:21Z",
     "message": "NMO custom resources have been found",
     "reason": "UpgradeBlocked",
     "status": "False",
     "type": "Upgradeable"
   }
   ```

2. Check for a ClusterServiceVersion (CSV) warning event such as the following:

   ```text
   Warning  NotUpgradeable   2m12s (x5 over 2m50s) kubevirt-hyperconvergedNode
   Maintenance Operator custom resources nodemaintenances.nodemaintenance.
     kubevirt.io have been found.
   Please remove them to allow upgrade. You can use NMO standalone operator
     if keeping the node(s) under maintenance is still required.
   ```

3. Check for NMO CRs belonging to the `nodemaintenance.kubevirt.io` API
group:

   ```bash
   $ oc get nodemaintenances.nodemaintenance.kubevirt.io
   ```

   Example output:

   ```text
   NAME                   AGE
   nodemaintenance-test   5m33s
   ```

## Mitigation

Remove all NMO CRs belonging to the `nodemaintenance.nodemaintenance.
kubevirt.io/` API group. After the integrated NMO resources are removed,
the alert is cleared and you can upgrade.

If a node must remain in maintenance mode during upgrade, install the Node
Maintenance Operator from OperatorHub. Then, create an NMO CR belonging
to the `nodemaintenance.nodemaintenance.medik8s.io/v1beta1` API group and
version for the node.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the Diagnosis procedure.
