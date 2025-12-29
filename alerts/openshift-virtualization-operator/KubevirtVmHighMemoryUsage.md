# KubevirtVmHighMemoryUsage [Deprecated]

This alert is deprecated. You can safely ignore or silence it.


## Meaning

This alert fires when a container hosting a virtual machine (VM) has less than
20 MB free memory.

## Impact

The virtual machine running inside the container is terminated by the runtime
if the container's memory limit is exceeded.

## Diagnosis

1. Obtain the `virt-launcher` pod details:

   ```bash
   $ oc get pod <virt-launcher> -o yaml
   ```

2. Identify `compute` container processes with high memory usage in the
`virt-launcher` pod:

   ```bash
   $ oc exec -it <virt-launcher> -c compute -- top
   ```

## Mitigation

- Increase the memory limit in the `VirtualMachine` specification as in the
following example:

   ```yaml
   spec:
     running: false
     template:
       metadata:
         labels:
           kubevirt.io/vm: vm-name
       spec:
         domain:
           resources:
             limits:
               memory: 200Mi
             requests:
               memory: 128Mi
   ```
