# UnmanagedRoutes

## Meaning

This alert fires when there is a Route owned by an unmanaged Ingress.

## Impact

The ingress-to-route controller does not remove Routes that earlier versions
of OpenShift created for Ingresses that specify `spec.ingressClassName`.
Thus, these Routes will continue to be in effect.
OpenShift does not update such Routes and does not recreate
them if the user deletes them.

In case any Routes existed in this state the alert would
help the administrator know that the Routes needed to be deleted,
or the Ingress modified to specify an appropriate IngressClass so
that OpenShift would once again reconcile the Routes.

## Diagnosis

Check for alert messages on the UI.
Inspect the ingress object.
Inspect the route object.
Check the logs of `cluster-openshift-controller-manager-operator`

## Mitigation

This alert will help the administrator to specify an appropriate
IngressClass in the Ingress object so that OpenShift would once
again reconcile the Routes.
