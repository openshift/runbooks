# IngressWithoutClassName

## Meaning

This alert fires when there is an Ingress with an unset IngressClassName
for longer than one day.

## Impact

It is possible that a user could have created an Ingress with
some nonempty value for spec.ingressClassName that did not match an
OpenShift IngressClass object, and nevertheless intended for OpenShift
to expose this Ingress.  Again, it is impossible to determine reliably
what a user's intent was in such a scenario, but as OpenShift exposed
such an Ingress before this enhancement, changing this behavior could
break existing applications.

So, we considered modifying the ingress operator
to list all Ingresses and Routes in the cluster and publish a metric
for Routes that were created for Ingresses that OpenShift no longer
manage.

## Diagnosis

Check for alert messages on the UI.
Inspect the ingress object.
Inspect the route object. Check the status of it.
Check the logs of `cluster-openshift-controller-manager-operator`

## Mitigation
Figure out why the route which was created by ingress which OpenShift
no longer manages.
Delete that ingress and route if it is no longer needed.
