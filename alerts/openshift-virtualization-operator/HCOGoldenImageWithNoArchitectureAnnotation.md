# HCOGoldenImageWithNoArchitectureAnnotation

## Meaning

When running on a heterogeneous cluster, the DataImportCronTemplate objects
(DICTs; also known as golden images) in the hyperconverged cluster operator
(HCO) must be annotated with the `ssp.kubevirt.io/dict.architectures`
annotation. The value of this annotation is a list of architectures supported by
the image, which is defined in each DICT.

For pre-defined DICTs, this annotation is already set. For custom DICTs
(user-defined DICTs), this annotation must be set by the user in the
`HyperConverged` custom resource (CR).

For each DICT, if the `ssp.kubevirt.io/dict.architectures` annotation is
missing, HCO triggers the `HCOGoldenImageWithNoArchitectureAnnotation`
alert for this specific DICT.

> **Note:** This alert only triggers if the `enableMultiArchBootImageImport`
> feature gate is enabled in the `HyperConverged` CR.

## Impact

This alert triggers when the golden image created for this DICT does not have a
defined architecture. If this image is used as a boot image for a virtual
machine (VM), and the VM will is scheduled on a node with a CPU architecture
different than the image architecture, the VM fails to start.

## Diagnosis

1. Read the `HyperConverged` CR:

```bash
  # Get the namespace of the HyperConverged CR
  $ NAMESPACE="$(oc get hyperconverged -A --no-headers | awk '{print $1}')"

  # Read the HyperConverged CR
  $ oc get hyperconverged -n "${NAMESPACE}" -o yaml
```

2. If this command lists any DICT objects under the
`spec.dataImportCronTemplates` field in the `HyperConverged` CR, check whether
the `ssp.kubevirt.io/dict.architectures` annotation is set for each of them. If
the annotation is not set, then this alert is triggered.

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

The `status.nodeInfo.workloadsArchitectures` field in the `HyperConverged` CR
shows the list of architectures that are supported by the cluster.

User-defined DICTs are defined in the `HyperConverged` CR, in the
`spec.dataImportCronTemplates` field.

## Mitigation

1. Check what architectures are supported by the image:

   ```bash
   $ podman manifest inspect your-registry/your-image:latest
   ```

   For details, see the
  [podman manifest inspect
documentation](https://docs.podman.io/en/latest/markdown/podman-manifest-inspect.1.html).

If the image is a multi-architecture manifest ("fat manifest"), it includes the
`manifests` field, which is a list of architectures supported by the image. If
the image is not a multi-architecture manifest, you need to find out what
is its architecture.

2. Edit the `HyperConverged` CR to add the missing
`ssp.kubevirt.io/dict.architectures` annotation.

  The format of the annotation is a comma-separated list of architectures,
  for example: `amd64,arm64,s390x`.

3. If the image does not support any of the architectures supported by the
cluster, rebuild the image for one or more of the architectures supported
by the cluster, or remove the DICT from the `HyperConverged` CR.

  For more information about building multi-architecture images, see the
  [podman
documentation](https://docs.podman.io/en/latest/markdown/podman-manifest-create.1.html).

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.