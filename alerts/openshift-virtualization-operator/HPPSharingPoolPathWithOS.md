# HPPSharingPoolPathWithOS
<!-- Edited by Jiří Herrmann, 10 Nov 2022 -->

## Meaning

This alert fires when the hostpath provisioner (HPP) shares a file
system with other critical components, such as `kubelet` or the operating
system (OS).

HPP dynamically provisions hostpath volumes to provide storage for
persistent volume claims (PVCs).

## Impact

A shared hostpath pool puts pressure on the node's disks. The node
might have degraded performance and stability.

## Diagnosis

1. Configure the `HPP_NAMESPACE` environment variable:

   ```bash
   $ export HPP_NAMESPACE="$(oc get deployment -A | \
     grep hostpath-provisioner-operator | awk '{print $1}')"
   ```

2. Obtain the status of the `hostpath-provisioner-csi` daemon set
pods:

   ```bash
   $ oc -n $HPP_NAMESPACE get pods | grep hostpath-provisioner-csi
   ```

3. Check the `hostpath-provisioner-csi` logs to identify the shared
pool and path:

   ```bash
   $ oc -n $HPP_NAMESPACE logs <csi_daemonset> -c hostpath-provisioner
   ```

   Example output:

   ```text
   I0208 15:21:03.769731       1 utils.go:221] pool (<legacy, csi-data-dir>/csi),
   shares path with OS which can lead to node disk pressure
   ```

## Mitigation

Using the data obtained in the Diagnosis section, try to prevent the
pool path from being shared with the OS. The specific steps vary based
on the node and other circumstances.

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.
