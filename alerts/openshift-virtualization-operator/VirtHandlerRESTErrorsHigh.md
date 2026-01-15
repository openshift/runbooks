# VirtHandlerRESTErrorsHigh [Deprecated]

This alert has been deprecated; it does not indicate a genuine issue. If
triggered, it may be safely ignored and silenced.

## Meaning

More than 5% of REST calls failed in `virt-handler` in the last 60 minutes. This
alert usually indicates that the `virt-handler` pods have partially lost
connection to the API server.

This error is frequently caused by one of the following problems:

- The API server is overloaded, which causes timeouts. To verify if this is the
case, check the metrics of the API server, and view its response times and
overall calls.

- The `virt-handler` pod cannot reach the API server. This is commonly caused by
DNS issues on the node and networking connectivity issues.

## Impact

Node-related actions, such as starting and migrating workloads, are delayed on
the node that `virt-handler` is running on. Running workloads are not affected,
but reporting their current status might be delayed.

## Diagnosis

1. Set the `NAMESPACE` environment variable:

   ```bash
   $ export NAMESPACE="$(oc get kubevirt -A -o custom-columns="":.metadata.namespace)"
   ```

2. List the available `virt-handler` pods to identify the failing `virt-handler`
pod:

   ```bash
   $ oc get pods -n $NAMESPACE -l=kubevirt.io=virt-handler
   ```

3. Check the failing `virt-handler` pod log for error messages when connecting
to the API server:

   ```bash
   $ oc logs -n $NAMESPACE <virt-handler>
   ```

   Example error message:

   ```json
   {"component":"virt-handler","level":"error","msg":"Can't patch node my-node","pos":"heartbeat.go:96","reason":"the server has received too many API requests and has asked us to try again later","timestamp":"2023-11-06T11:11:41.099883Z","uid":"132c50c2-8d82-4e49-8857-dc737adcd6cc"}
   ```

## Mitigation

If the `virt-handler` cannot connect to the API server, delete the pod to force
a restart:

```bash
$ oc delete -n $NAMESPACE <virt-handler>
```

If you cannot resolve the issue, log in to the
[Customer Portal](https://access.redhat.com) and open a support case,
attaching the artifacts gathered during the diagnosis procedure.