# CnaoDown
<!-- Edited by Jiří Herrmann, 3 Nov 2022 -->

## Meaning

This alert fires when the `Cluster-network-addons-operator` (CNAO) is down.
The CNAO deploys additional networking components on top of the cluster.

## Impact

If the CNAO is not running, the cluster cannot reconcile changes to virtual
machine components. As a result, the changes might fail to take effect.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get deployment -A | grep \
     cluster-network-addons-operator | awk '{print $1}')"
   ```

2. Check the status of the `cluster-network-addons-operator` pod:

   ```bash
   $ oc -n $NAMESPACE get pods -l name=cluster-network-addons-operator
   ```

3. Check the `cluster-network-addons-operator` logs for error messages:

   ```bash
   $ oc -n $NAMESPACE logs -l name=cluster-network-addons-operator
   ```

4. Obtain the details of the `cluster-network-addons-operator` pods:

   ```bash
   $ oc -n $NAMESPACE describe pods -l name=cluster-network-addons-operator
   ```

## Mitigation

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the Diagnosis procedure.
