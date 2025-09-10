# HCOGoldenImageWithNoSupportedArchitecture

## Meaning

When running on a heterogeneous cluster, the DataImportCronTemplate objects
(DICTs; also known as golden images) in the hyperconverged cluster operator
(HCO) must be annotated with the `ssp.kubevirt.io/dict.architectures`
annotation. The value of this annotation is a list of architectures
supported by the image, which is defined in each DICT.

For pre-defined DICTs, this annotation is already set. For custom DICTs
(user-defined DICTs), this annotation must be set by the user in the
`HyperConverged` custom resource (CR).

For each DICT, if the annotation does not include any architecture that is
supported by the cluster (which means that there is no node in the cluster with
the architectures listed in the DICT annotation), HCO triggers the
`HCOGoldenImageWithNoSupportedArchitecture` alert for this specific DICT.

> **Note:** This only triggers if the `enableMultiArchBootImageImport`
> feature gate is enabled in the `HyperConverged` CR.

## Impact

This alert triggers when the DICT is not supported by any of the nodes in the
cluster. HCO will not populate the SSP CR with this DICT, and the golden image
will not be available for use in the cluster.

## Diagnosis

1. Read the HyperConverged CR:

```bash
  # Get the namespace of the HyperConverged CR
  $ NAMESPACE="$(oc get hyperconverged -A --no-headers | awk '{print $1}')"

  # Read the HyperConverged CR
  $ oc get hyperconverged -n "${NAMESPACE}" -o yaml
```

2. Examine the following fields in the `HyperConverged` CR status:

   - The `status.nodeInfo.workloadsArchitectures` shows a list of architectures
supported by the cluster.
   - The `status.dataImportCronTemplates` field shows the list of DICTs that are
managed by HCO.

3. Find the specific DICT object that is triggering this alert by its name, as
specified in the alert message.

4. Check the `ssp.kubevirt.io/dict.architectures` annotation of the DICT. Unlike
the annotation in the `spec` field, this annotation contains only architectures
that are supported by both the image and by the cluster.

   - If the annotation is empty, there is no architecture supported by the image
and by the cluster.

   - The DICT `status` field will include the `conditions` field, with the
`Deployed` condition set to `False`, and the `reason` field set to
`UnsupportedArchitectures`. For DICTs with supported architectures, the `status`
field will not contain the `conditions` field.

   - The `status.workloadsArchitectures` field of the DICT shows the list of
architectures supported by the image, which was set in the
`ssp.kubevirt.io/dict.architectures` annotation in the source DICT.

### Example

```yaml
apiVersion: hco.kubevirt.io/v1beta1
kind: HyperConverged
...
status:
  ...
  dataImportCronTemplates:
    - metadata:
        annotations:
          ssp.kubevirt.io/dict.architectures: ""
        name: my-image
      spec:
          ...
      status:
        conditions:
          - message: DataImportCronTemplate has no supported architectures for the current
              cluster
            reason: UnsupportedArchitectures
            status: "False"
            type: Deployed
        originalSupportedArchitectures: someUnsupportedArch,otherUnsupportedArch
```

## Mitigation

The steps to mitigate this issue vary based on whether you are using pre-defined
DICTs or user-defined DICTs.

### Pre-defined DataImportCronTemplates

Pre-defined DICTs are not defined in the `spec.dataImportCronTemplates`
field in the `HyperConverged` CR. Instead, they are defined internally in the
HCO application.

All pre-defined DICTs are annotated with the `ssp.kubevirt.io/dict.architectures`
annotation, and all of them support the `amd64`, `arm64`, and `s390x`
architectures. If the cluster does not support any of these architectures,
you cannot use pre-defined DICTs in the cluster.

To mitigate the `HCOGoldenImageWithNoSupportedArchitecture` issue on pre-defined
DICTs, you can do one of the following:

- Add supported nodes to the cluster

- Disable the pre-defined DICTs in the `HyperConverged` CR to turn this alert
off:

  1. Find the DICT(s) that you want to disable in the
  `status.dataImportCronTemplates` field of the `HyperConverged` CR, as
  described [in the Diagnosis section](#diagnosis).

  2. Add the DICT to the `spec.dataImportCronTemplates` field in the
  `HyperConverged` CR. Add the `dataimportcrontemplate.kubevirt.io/enable`
  annotation with the value `false` to the DICT. Only the DICT name and the
  annotation are required.

       For example, to disable the `centos-stream10-image-cron` DICT:
       ```yaml
       apiVersion: hco.kubevirt.io/v1beta1
       kind: HyperConverged
       metadata:
         name: kubevirt-hyperconverged
       spec:
         dataImportCronTemplates:
         - metadata:
             name: centos-stream10-image-cron
             annotations:
               dataimportcrontemplate.kubevirt.io/enable: 'false'
        ```

     If you have a self-built image that is supported by the nodes in the
cluster,
you can modify the pre-defined DICT to use your image. To do so, add the DICT to
the `spec.dataImportCronTemplates` field in the `HyperConverged` CR and modify
its `spec.source.registry` field.

  > Tip: you can find the pre-defined DICTs in the
  > `status.dataImportCronTemplates` field of the `HyperConverged` CR, as
  > described [in the Diagnosis section](#diagnosis). Afterwards, you can copy
  > the DICT from the field, and modify it in the `spec.dataImportCronTemplates`
  > field.

3. Set the `ssp.kubevirt.io/dict.architectures` annotation to include all the
architectures supported by your image.

   For example:
  ```yaml
  apiVersion: hco.kubevirt.io/v1beta1
  kind: HyperConverged
  metadata:
    name: kubevirt-hyperconverged
  spec:
    dataImportCronTemplates:
    - metadata:
      annotations:
        cdi.kubevirt.io/storage.bind.immediate.requested: "true"
        ssp.kubevirt.io/dict.architectures: arch1,arch2
      name: centos-stream10-image-cron
    spec:
      garbageCollect: Outdated
      managedDataSource: centos-stream10
      schedule: "0 */12 * * *"
      template:
        spec:
          source:
            registry:
              url: docker://your-registry/your-image:latest
          storage:
            resources:
              requests:
                storage: 10Gi
  ```

### User-defined DataImportCronTemplates

User-defined DICTs are defined in the `spec.dataImportCronTemplates` field of
the of the HyperConverged CR.

1. Check what architectures are supported by the image:

   ```bash
   $ podman manifest inspect your-registry/your-image:latest
   ```

   For details, see the [podman manifest inspect
documentation](https://docs.podman.io/en/latest/markdown/podman-manifest-inspect.1.html).

   If the image is a multi-architecture manifest ("fat manifest"), it includes
the
`manifests` field, which is a list of architectures supported by the image. If
the image is not a multi-architecture manifest, you need to find out what is its
architecture.

2. Check that the `ssp.kubevirt.io/dict.architectures` annotation is set with
the correct value. If not, edit the `HyperConverged` CR to fix the annotation.
The format of the annotation is a comma-separated list of architectures, for
example: `amd64,arm64,s390x`.

3. If the image does not support any of the architectures supported by the
cluster, rebuild the image for one or more of the architectures supported by the
cluster, or remove the DICT from the `HyperConverged` CR.

It is also possible to disable the DICT, by adding it
the `dataimportcrontemplate.kubevirt.io/enable` annotation, with the value
of `false`. For example:

  ```yaml
  apiVersion: hco.kubevirt.io/v1beta1
  kind: HyperConverged
  metadata:
    name: kubevirt-hyperconverged
  spec:
    dataImportCronTemplates:
    - metadata:
      annotations:
        dataimportcrontemplate.kubevirt.io/enable: "false"
        ssp.kubevirt.io/dict.architectures: unsupported-arch1,unsupported-arch2
      name: my-image
    spec:
      ...
  ```

For more information about building multi-architecture images, see the
[podman documentation](https://docs.podman.io/en/latest/markdown/podman-manifest-create.1.html).

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.