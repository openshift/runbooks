# HCOGoldenImageWithNoArchitectureAnnotation

## Meaning

When running on a heterogeneous cluster, a cluster with nodes of different
architectures, the DataImportCronTemplates (DICTs; also known as golden
images), in the hyperconverged cluster operator (HCO) should be annotated with
the `ssp.kubevirt.io/dict.architectures` annotation, where the value is the
list of the architectures supported by the image, that is defined in each DICT.

For pre-defined DICTs, this annotation is already set, but for custom DICTs
(user defined DICTs), this annotation must be set by the user in the
HyperConverged custom resource (CR).

For each DICT, if the `ssp.kubevirt.io/dict.architectures` annotation is
missing, Then HCO will trigger the `HCOGoldenImageWithNoArchitectureAnnotation`
alert for this specific DICT.

> **Note:** This alert is only triggered, if the `enableMultiArchBootImageImport`
> feature gate is enabled in the HyperConverged CR.

## Impact

When this alert is triggered, it means that the golden image created for this
DICT, is with undefined architecture. There is a risk that when this image
will be used as a boot image for a virtual machine, and the virtual machine
will be scheduled on a node with a CPU architecture different than the image
architecture,the virtual machine will fail to start.

## Diagnosis

Read the HyperConverged CR:

```bash
   # Get the namespace of the HyperConverged CR
$ NAMESPACE="$(oc get hyperconverged -A --no-headers | awk '{print $1}')"

#Read the HyperConverged CR
$ oc get hyperconverged -n "${NAMESPACE}" -o yaml
```

Go over the output of the command. If there are DICT objects under the
`spec.dataImportCronTemplates` field in the HyperConverged CR, then for each
one of them, check if the `ssp.kubevirt.io/dict.architectures` annotation is
set. If the annotation is not set, then this alert is triggered.

Below is an example for a HyperConverged CR with a valid DICT with the
`ssp.kubevirt.io/dict.architectures` annotation set:
```yaml
apiVersion: hco.kubevirt.io/v1beta1
kind: HyperConverged
...
spec:
  ...
  dataImportCronTemplates:
    - metadata:
        annotations:
          ...
          ssp.kubevirt.io/dict.architectures: amd64
        name: the-name-of-the-dict
      spec:
        ...
```

The `status.nodeInfo.workloadsArchitectures` shows the list of architectures
that are supported by the cluster.

User-defined DICTs are defined in the HyperConverged CR, in the
`spec.dataImportCronTemplates` field.

## Mitigation
First, check what architectures are supported by the image. You can use the
following command:

```bash
$ podman manifest inspect your-registry/your-image:latest
```

See here for
the [podman manifest inspect
documentation](https://docs.podman.io/en/latest/markdown/podman-manifest-inspect.1.html).

If the image is multi architecture manifest (fat manifest), it will include the
`manifests` field, which is a list of architectures supported by the image. If
the image is not a multi architecture manifest, you will need to find out what
is its architecture.

Then, edit the HyperConverged CR, to add the missing `ssp.kubevirt.io/dict.architectures`
annotation.

The format of the annotation is a comma-separated list of architectures;
e.g., `amd64,arm64,s390x`.

If the image does not support any of the architectures supported by the
cluster, you will need to either rebuild the image for one or more of
the architectures supported by the cluster, or remove the DICT from the
HyperConverged CR.

Find some more information about building multi-architecture images, see the
[podman documentation](https://docs.podman.io/en/latest/markdown/podman-manifest-create.1.html).

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.