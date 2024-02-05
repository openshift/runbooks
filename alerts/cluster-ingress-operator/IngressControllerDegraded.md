# IngressControllerDegraded

## Meaning

This alert fires when the IngressController status is degraded.

## Impact

The routers won't be running in the cluster. This will cause an outage while
accessing the applications.

## Diagnosis

Ingress Controller may be degraded due to one or more reasons.

- Check the ingress operator logs using the following command:
```sh
oc logs <ingress operator pod> -n openshift-ingress-operator
```
- Check the router logs using the following commands:
```sh
oc logs <router pod> -n openshift-ingress
```
- Check the yaml file of the ingress controller and operator to see the reason
  for failure:
```sh
oc get ingresscontroller <ingresscontroller name> -n openshift-ingress-operator -o yaml
```

```sh
oc get deployment -n openshift-ingress-operator -o yaml
```

```sh
oc get events
```

## Mitigation

Try to fix the issue based on what you see in the status of yaml and errors
in the logs from the above mentioned
commands.