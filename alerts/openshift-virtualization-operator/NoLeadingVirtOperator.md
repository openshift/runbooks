# NoLeadingVirtOperator

## Meaning

This alert fires when no `virt-operator` pod with a leader lease has been
detected for 10 minutes, although the `virt-operator` pods are in a `Ready`
state. The alert indicates that no leader pod is available.

The `virt-operator` is the first Operator to start in a cluster. Its primary
responsibilities include the following:

- Installing, live updating, and live upgrading a cluster

- Monitoring the lifecycle of top-level controllers, such as `virt-controller`,
`virt-handler`, `virt-launcher`, and managing their reconciliation

- Certain cluster-wide tasks, such as certificate rotation and infrastructure
management

The `virt-operator` deployment has a default replica of 2 pods, with one pod
holding a leader lease.

## Impact

This alert indicates a failure at the level of the cluster. As a result,
critical cluster-wide management functionalities, such as certification
rotation, upgrade, and reconciliation of controllers, might not be available.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o jsonpath='{.items[].metadata.namespace}')"
   ```

2. Obtain the status of the `virt-operator` pods:

   ```bash
   $ oc -n $NAMESPACE get pods -l kubevirt.io=virt-operator
   ```

3. Check the `virt-operator` pod logs to determine the leader status:

   ```bash
   $ oc -n $NAMESPACE logs -l kubevirt.io=virt-operator | grep leader
   ```

   Leader pod example:

   ```text
   {"component":"virt-operator","level":"info","msg":"Attempting to acquire leader status","pos":"application.go:400","timestamp":"2021-11-30T12:15:18.635387Z"}
   I1130 12:15:18.635452       1 leaderelection.go:243] attempting to acquire leader lease <namespace>/virt-operator...
   I1130 12:15:19.216582       1 leaderelection.go:253] successfully acquired lease <namespace>/virt-operator
   {"component":"virt-operator","level":"info","msg":"Started leading","pos":"application.go:385","timestamp":"2021-11-30T12:15:19.216836Z"}
   ```

   Non-leader pod example:

   ```text
   {"component":"virt-operator","level":"info","msg":"Attempting to acquire leader status","pos":"application.go:400","timestamp":"2021-11-30T12:15:20.533696Z"}
   I1130 12:15:20.533792       1 leaderelection.go:243] attempting to acquire leader lease <namespace>/virt-operator...
   ```

4. Obtain the details of the affected `virt-operator` pods:

   ```bash
   $ oc -n $NAMESPACE describe pod <virt-operator>
   ```

## Mitigation

Based on the information obtained during the diagnosis procedure, try to
identify the root cause and resolve the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.