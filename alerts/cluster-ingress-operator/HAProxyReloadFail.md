# HAProxyReloadFail

## Meaning

This alert fires when HAProxy fails to reload its configuration, which will
result in the router not picking up recently created or modified routes.

## Impact

The router won't pick up recently created or modified routes. This may cause
an outage for critical applications.

## Diagnosis

Check the router logs:
```sh
oc logs <router pod> -n openshift-ingress
```

Check if any recently added configuration in the haproxy config via ingress
controller CR caused the issue.

## Mitigation

Try to fix the configuration of the haproxy via ingress controller CR on the
basis of the output of the logs.