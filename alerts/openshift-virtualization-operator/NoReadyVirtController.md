# NoReadyVirtController
<!-- Edited by Jiří Herrmann, 9 Nov 2022 -->

## Meaning

This alert fires when no available `virt-controller` devices have been
detected for 5 minutes.

The `virt-controller` devices monitor the custom resource definitions of
virtual machine instances (VMIs) and manage the associated pods. The devices
create pods for VMIs and manage the lifecycle of the pods.

Therefore, `virt-controller` devices are critical for all cluster-wide
virtualization functionality.

## Impact
Any actions related to VM lifecycle management fail. This notably includes
launching a new VMI or shutting down an existing VMI.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A
     -o custom-columns="":.metadata.namespace)"
   ```

2. Verify the number of `virt-controller` devices:

   ```bash
   $ oc get deployment -n $NAMESPACE virt-controller \
     -o jsonpath='{.status.readyReplicas}'
   ```

3. Check the status of the `virt-controller` deployment:

   ```bash
   $ oc -n $NAMESPACE get deploy virt-controller -o yaml
   ```

4. Obtain the details of the `virt-controller` deployment to check for
status conditions such as crashing pods or failure to pull images:

   ```bash
   $ oc -n $NAMESPACE describe deploy virt-controller
   ```

5. Obtain the details of the `virt-controller` pods:

   ```bash
   $ get pods -n $NAMESPACE | grep virt-controller
   ```

6. Check the logs of the `virt-controller` pods for error messages:

   ```bash
   $ oc logs -n $NAMESPACE <virt-controller>
   ```

7. Check the nodes for problems, suchs as a `NotReady` state:

   ```bash
   $ oc get nodes
   ```

## Mitigation

Based on the information obtained during Diagnosis, try to find and resolve
the cause of the issue.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the Diagnosis procedure.
