# VirtControllerRESTErrorsHigh
<!-- Edited by Jiří Herrmann, 8 Nov 2022 -->

## Meaning

More than 5% of REST calls failed in `virt-controller` in the last 60 minutes.

This is most likely because `virt-controller` has partially lost connection
to the API server.

This error is frequently caused by one of the following problems:

- The API server is overloaded, which causes timeouts. To verify if this
is the case, check the metrics of the API server, and view its response
times and overall calls.

- The `virt-controller` pod cannot reach the API server. This is commonly
caused by DNS issues on the node and networking connectivity issues.

## Impact

Node-related actions, such as starting and migrating, and scheduling virtual
machines, are delayed. Running workloads are not affected, but reporting
their current status might be delayed.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A \
     -o custom-columns="":.metadata.namespace)"
   ```

2. List the available `virt-controller` pods:

   ```bash
   $ oc get pods -n $NAMESPACE -l=kubevirt.io=virt-controller
   ```

3. Check the `virt-controller` logs for error messages when connecting
to the API server:

   ```bash
   $ oc logs -n  $NAMESPACE <virt-controller>
   ```

## Mitigation

- If the `virt-controller` pod cannot connect to the API server, delete
the pod to force a restart:

  ```bash
  $ oc delete -n $NAMESPACE <virt-controller>
  ```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the Diagnosis procedure.
