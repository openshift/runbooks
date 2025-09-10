# HCOMultiArchGoldenImagesDisabled

## Meaning

DataImportCron objects (DIC; also known as golden images) are used to create
boot images for virtual machines (VMs). The images are preloaded in the cluster,
and then they are used to create VM boot disks with a specific operating system.
By default, the preloaded images use the architecture of the cluster node that
was used to create the image.

If the `enableMultiArchBootImageImport` feature gate is enabled in the
HyperConverged custom resource (CR), multiple preloaded images are created for
each DataImportCronTemplate (DICT), one for each architecture supported by the
cluster and by the original image.

This allows the VMs to be scheduled on nodes with the same architecture as the
preloaded image.

This alert is triggered when running on a heterogeneous cluster (a cluster with
nodes of different architectures) while the `enableMultiArchBootImageImport`
feature gate is disabled in the HyperConverged CR.

## Impact

When running on a heterogeneous cluster, the preloaded image may use a different
architecture than the architecture of the node that the VM is scheduled on.

In this case, the VM will fail to start, as the image architecture is not
compatible with the node architecture.

## Diagnosis

HCO checks the workload node architectures in the cluster. By default, HCO
considers the worker nodes as the workload nodes. If the
`spec.workloads.nodePlacement` field in the HyperConverged CR is populated,
HCO considers the nodes that match the node selector in this field as the
workload nodes.

HCO publishes the list of the workload node architectures in the
`status.nodeInfo.workloadsArchitectures` field in the HyperConverged CR.

Read the HyperConverged CR:

```bash
oc get hyperconverged -n openshift-cnv kubevirt-hyperconverged -o yaml
```

The result will look similar to this:

```yaml
apiVersion: hco.kubevirt.io/v1beta1
kind: HyperConverged
spec:
  ...
  workloads: # check if the spec.workloads.nodePlacement field is populated
    nodePlacement:
  ...
status:
  ...
  nodeInfo:
    workloadsArchitectures:
      - amd64
      - arm64
...
```

## Mitigation

To address this issue, you can either enable the multi-arch boot image feature,
or modify the workloads node placement in the HyperConverged CR to include only
nodes with a single architecture.

### Enable the multi-arch boot image feature

The multi-arch boot image feature is in the alpha stage, and it is not enabled
by default. Enabling this feature will result in the creation of multiple
preloaded images for each DataImportCronTemplate (DICT), one for each
architecture supported by the cluster, and by the original image. However, this
feature is not generally available, and it is not fully supported.

To enable the multi-arch boot image feature, set the
`enableMultiArchBootImageImport` feature gate in the HyperConverged CR to `true`

If the HyperConverged CR contains the `spec.dataImportCronTemplates` field,
and this field is not empty, then you may need to add the
`ssp.kubevirt.io/dict.architectures` annotation to each DICT object in this
field. See
the
[HCOGoldenImageWithNoArchitectureAnnotation](HCOGoldenImageWithNoArchitectureAnnotation.md)
runbook for more details.

Edit the HyperConverged CR:

```bash
oc edit hyperconverged -n openshift-cnv kubevirt-hyperconverged -o yaml
```

The editor will be opened with the HyperConverged CR in YAML format.

Edit the CR to set the `enableMultiArchBootImageImport` feature gate to `true`,
and to add the `ssp.kubevirt.io/dict.architectures` annotation to each DICT
object in the `spec.dataImportCronTemplates` field, if needed.

```yaml
apiVersion: hco.kubevirt.io/v1beta1
kind: HyperConverged
spec:
  dataImportCronTemplates:
    ...
  ...
  featureGates:
    ...
    enableMultiArchBootImageImport: true
    ...
```

Save the changes and exit the editor.

### Modify the Workloads Node Placement

If you do not want to enable the multi arch boot image feature, you can modify
the workloads node placement in the HyperConverged CR to include only nodes with
a single architecture.

Edit the HyperConverged CR:
<!--DS:
```bash
oc edit hyperconverged -n openshift-cnv kubevirt-hyperconverged -o yaml
```
-->

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.