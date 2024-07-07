# HAProxyDown

## Meaning

This alert fires when metrics report that HAProxy is down.

## Impact

Access to routes will fail. It may cause a severe outage for critical applications.

## Diagnosis

- Check the router logs:
```sh
oc logs <router pod> -n openshift-ingress
```

- Check the events:
```sh
oc get events -n openshift-ingress
```

- Check the load on the system where the routers are hosted.

## Mitigation

Based on the diagnosis, try to figure out the issue.
If the issue is configuration related then try to fix the haproxy config.
If the issue is load related try to fix the issues at infrastructure level.
