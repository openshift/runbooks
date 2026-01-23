# Naming schemas

This document identifies the schema used by the `cluster-logging-operator` to
generate the configuration for the log collector from a `ClusterLogForwarder` instance.

## Outputs

### Non-Lokistack Outputs

The naming schema is `output_<OUTPUT_NAME>`,
where the `<OUTPUT_NAME>` corresponds to the CLF's named output, with all
punctuations replaced by underscores
(e.g., `output_my_splunk` will correspond to a CLF output named `my-splunk`).

### Lokistack Outputs

The naming schema is `output_<OUTPUT_NAME>_<INPUT_TYPE>`,
where the `<OUTPUT_NAME>` corresponds to the CLF's named Lokistack output,
with all punctuations replaced by underscores and `<INPUT_TYPE>` is the tenant.
(e.g `output_my_lokistack_application` will correspond to a CLF output named
`my-lokistack` receiving `application` logs.)
